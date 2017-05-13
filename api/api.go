/*
 *    Copyright (C) 2015-2017 Christian Muehlhaeuser
 *
 *    This program is free software: you can redistribute it and/or modify
 *    it under the terms of the GNU Affero General Public License as published
 *    by the Free Software Foundation, either version 3 of the License, or
 *    (at your option) any later version.
 *
 *    This program is distributed in the hope that it will be useful,
 *    but WITHOUT ANY WARRANTY; without even the implied warranty of
 *    MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 *    GNU Affero General Public License for more details.
 *
 *    You should have received a copy of the GNU Affero General Public License
 *    along with this program.  If not, see <http://www.gnu.org/licenses/>.
 *
 *    Authors:
 *      Christian Muehlhaeuser <muesli@gmail.com>
 */

// Package api is Beehive's RESTful api for introspection and configuration
package api

import (
	"bytes"
	"fmt"
	"net/http"
	"net/url"
	"path"
	"regexp"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/emicklei/go-restful"
	"github.com/muesli/smolder"

	bee "github.com/muesli/beehive/bees"

	"github.com/muesli/beehive/api/context"
	"github.com/muesli/beehive/api/resources/actions"
	"github.com/muesli/beehive/api/resources/bees"
	"github.com/muesli/beehive/api/resources/chains"
	"github.com/muesli/beehive/api/resources/hives"
	"github.com/muesli/beehive/app"
)

var (
	bind         string
	canonicalURL string
)

const (
	defaultBind = "localhost:8181"
	defaultURL  = "http://localhost:8181"
)

// CanonicalURL returns the canonical URL of the API
func CanonicalURL() *url.URL {
	u, _ := url.Parse(canonicalURL)
	return u
}

func escapeURL(u string) string {
	return strings.Replace(url.QueryEscape(u), "%2F", "/", -1)
}

// Try to read a local/embedded asset. Gracefully fail and read index.html if
// if not found.
func readAssetOrIndex(path string) (string, []byte, error) {
	b, err := Asset(path)
	if err != nil {
		path = "config/index.html"
		b, err = Asset(path)
		if err != nil {
			return path, nil, err
		}
	}

	return path, b, nil
}

func assetHandler(req *restful.Request, resp *restful.Response) {
	var rootdir string
	if strings.HasPrefix(req.Request.URL.Path, "/images/") {
		rootdir = "./assets/bees"
	} else {
		rootdir = "./config"
	}

	subpath := req.PathParameter("subpath")
	sourceFile, b, err := readAssetOrIndex(path.Join(rootdir, subpath))
	if err != nil {
		log.Errorln("Failed reading", sourceFile)
		http.Error(resp.ResponseWriter, "Failed reading file", http.StatusInternalServerError)
		return
	}

	log.Printf("serving %s ... (from %s)", sourceFile, req.PathParameter("subpath"))

	if sourceFile == "config/index.html" {
		// Since we patch the content of the files, we must drop the integrity SHA-sums
		// TODO: Would be nicer to recalculate them
		re := regexp.MustCompile("integrity=\"([^\"]*)\"")
		b = re.ReplaceAll(b, []byte{})
	}
	if defaultURL != canonicalURL {
		// We're serving files on a non-default canonical URL
		// Make sure the HTML we serve references API & assets with the correct URL
		b = bytes.Replace(b, []byte(defaultURL), []byte(canonicalURL), -1)
		b = bytes.Replace(b, []byte(escapeURL(defaultURL)), []byte(escapeURL(canonicalURL)), -1)
	}

	http.ServeContent(
		resp.ResponseWriter,
		req.Request,
		sourceFile,
		time.Now(),
		bytes.NewReader(b))
}

func oauth2Handler(req *restful.Request, resp *restful.Response) {
	errHTML := []byte("<html>Failed retrieving OAuth2 access-token. Please check your Beehive logs!</html>")

	params := strings.Split(req.PathParameter("subpath"), "/")
	log.Printf("OAuth2 callback received: %s", params)

	if len(params) != 3 {
		log.Errorln("OAuth2: Missing parameters:", params)
		resp.Write(errHTML)
		return
	}

	subpath := params[0]
	id := params[1]
	secret := params[2]
	log.Printf("OAuth2 app ID: %s", id)
	log.Printf("OAuth2 app secret: %s", secret)

	code := req.QueryParameter("code")
	log.Printf("OAuth2 code: %s", code)

	f := bee.GetFactory(subpath)
	if f == nil {
		log.Errorln("OAuth2: No such hive:", subpath)
		resp.Write(errHTML)
		return
	}
	token, err := (*f).OAuth2AccessToken(id, secret, code)
	if err != nil {
		log.Errorln("OAuth2: This hive does not support OAuth2:", subpath)
		resp.Write(errHTML)
		return
	}

	s := fmt.Sprintf("<html>You're now logged in with %s!<br/><br/>Access token:<br/>"+
		"<b>%s</b><br/><br/>"+
		"Copy & paste this token into Beehive's admin interface. You can safely close this tab then.</html>", subpath, token.AccessToken)
	resp.Write([]byte(s))
}

// Run sets up the restful API container and an HTTP server go-routine
func Run() {
	// to see what happens in the package, uncomment the following
	//restful.TraceLogger(log.New(os.Stdout, "[restful] ", log.LstdFlags|log.Lshortfile))

	// Setup web-service
	smolderConfig := smolder.APIConfig{
		BaseURL:    canonicalURL,
		PathPrefix: "v1/",
	}
	context := &context.APIContext{
		Config: smolderConfig,
	}

	wsContainer := smolder.NewSmolderContainer(smolderConfig, nil, nil)
	wsContainer.Router(restful.CurlyRouter{})
	ws := new(restful.WebService)
	ws.Route(ws.GET("/images/{subpath:*}").To(assetHandler))
	ws.Route(ws.GET("/oauth2/{subpath:*}").To(oauth2Handler))
	ws.Route(ws.GET("/{subpath:*}").To(assetHandler))
	ws.Route(ws.GET("/").To(assetHandler))
	wsContainer.Add(ws)

	func(resources ...smolder.APIResource) {
		for _, r := range resources {
			r.Register(wsContainer, smolderConfig, context)
		}
	}(
		&hives.HiveResource{},
		&bees.BeeResource{},
		&chains.ChainResource{},
		&actions.ActionResource{},
	)

	server := &http.Server{Addr: bind, Handler: wsContainer}
	go func() {
		log.Fatal(server.ListenAndServe())
	}()
}

func init() {
	app.AddFlags([]app.CliFlag{
		{
			V:     &bind,
			Name:  "bind",
			Value: defaultBind,
			Desc:  "Which address to bind Beehive's API & admin interface to",
		},
	})
	app.AddFlags([]app.CliFlag{
		{
			V:     &canonicalURL,
			Name:  "canonicalurl",
			Value: defaultURL,
			Desc:  "Canonical URL for the API & admin interface",
		},
	})
}

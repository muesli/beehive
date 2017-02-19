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
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path"
	"regexp"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/emicklei/go-restful"
	"github.com/muesli/smolder"

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

func escapeURL(u string) string {
	return strings.Replace(url.QueryEscape(u), "%2F", "/", -1)
}

func configFromPathParam(req *restful.Request, resp *restful.Response) {
	rootdir := "./config"

	subpath := req.PathParameter("subpath")
	if _, err := os.Stat(path.Join(rootdir, subpath)); os.IsNotExist(err) || len(subpath) == 0 {
		subpath = "index.html"
	}
	actual := path.Join(rootdir, subpath)
	log.Printf("serving %s ... (from %s)", actual, req.PathParameter("subpath"))

	b, err := ioutil.ReadFile(actual)
	if err != nil {
		log.Errorln("Failed reading", actual)
		http.Error(resp.ResponseWriter, "Failed reading file", http.StatusInternalServerError)
		return
	}

	if defaultURL != canonicalURL {
		// We're serving files on a non-default canonical URL
		// Make sure the HTML we serve references API & assets with the correct URL
		if actual == "config/index.html" {
			// Since we patch the content of the files, we must drop the integrity SHA-sums
			// TODO: Would be nicer to recalculate them
			re := regexp.MustCompile("integrity=\"([^\"]*)\"")
			b = re.ReplaceAll(b, []byte{})
		}
		b = bytes.Replace(b, []byte(defaultURL), []byte(canonicalURL), -1)
		b = bytes.Replace(b, []byte(escapeURL(defaultURL)), []byte(escapeURL(canonicalURL)), -1)
	}

	http.ServeContent(
		resp.ResponseWriter,
		req.Request,
		actual,
		time.Now(),
		bytes.NewReader(b))
}

func imageFromPathParam(req *restful.Request, resp *restful.Response) {
	rootdir := "./assets/bees"

	subpath := req.PathParameter("subpath")
	actual := path.Join(rootdir, subpath)
	log.Printf("serving %s ... (from %s)", actual, req.PathParameter("subpath"))
	http.ServeFile(
		resp.ResponseWriter,
		req.Request,
		actual)
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
	ws.Route(ws.GET("/images/{subpath:*}").To(imageFromPathParam))
	ws.Route(ws.GET("/{subpath:*}").To(configFromPathParam))
	ws.Route(ws.GET("/").To(configFromPathParam))
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

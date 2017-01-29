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
	"log"
	"net/http"
	"path"

	"github.com/emicklei/go-restful"

	"github.com/muesli/beehive/api/context"
	"github.com/muesli/beehive/api/resources/actions"
	"github.com/muesli/beehive/api/resources/bees"
	"github.com/muesli/beehive/api/resources/chains"
	"github.com/muesli/beehive/api/resources/hives"
	"github.com/muesli/smolder"
)

func configFromPathParam(req *restful.Request, resp *restful.Response) {
	rootdir := "./config"

	subpath := req.PathParameter("subpath")
	if len(subpath) == 0 {
		subpath = "index.html"
	}
	actual := path.Join(rootdir, subpath)
	log.Printf("serving %s ... (from %s)\n", actual, req.PathParameter("subpath"))
	http.ServeFile(
		resp.ResponseWriter,
		req.Request,
		actual)
}

func imageFromPathParam(req *restful.Request, resp *restful.Response) {
	rootdir := "./assets/bees"

	subpath := req.PathParameter("subpath")
	actual := path.Join(rootdir, subpath)
	log.Printf("serving %s ... (from %s)\n", actual, req.PathParameter("subpath"))
	http.ServeFile(
		resp.ResponseWriter,
		req.Request,
		actual)
}

// Run sets up the restful API container and an HTTP server go-routine
func Run() {
	// to see what happens in the package, uncomment the following
	//restful.TraceLogger(log.New(os.Stdout, "[restful] ", log.LstdFlags|log.Lshortfile))

	context := &context.APIContext{}

	// Setup web-service
	smolderConfig := smolder.APIConfig{
		BaseURL:    "http://localhost:8181",
		PathPrefix: "v1/",
	}

	wsContainer := smolder.NewSmolderContainer(smolderConfig, nil, nil)
	wsContainer.Router(restful.CurlyRouter{})
	ws := new(restful.WebService)
	ws.Route(ws.GET("/config/").To(configFromPathParam))
	ws.Route(ws.GET("/config/{subpath:*}").To(configFromPathParam))
	ws.Route(ws.GET("/images/{subpath:*}").To(imageFromPathParam))
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

	server := &http.Server{Addr: ":8181", Handler: wsContainer}
	go func() {
		log.Fatal(server.ListenAndServe())
	}()
}

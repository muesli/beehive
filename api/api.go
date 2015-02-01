/*
 *    Copyright (C) 2015 Christian Muehlhaeuser
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

// beehive's RESTful api for introspection and configuration
package api

import (
	"log"
	"net/http"
	_ "strconv"

	"github.com/emicklei/go-restful"
	_ "github.com/emicklei/go-restful/swagger"
)

func Run() {
	// to see what happens in the package, uncomment the following
	//restful.TraceLogger(log.New(os.Stdout, "[restful] ", log.LstdFlags|log.Lshortfile))

	wsContainer := restful.NewContainer()
	b := BeeResource{}
	b.Register(wsContainer)
	h := HiveResource{}
	h.Register(wsContainer)

	log.Println("Starting JSON API on localhost:8080")
	server := &http.Server{Addr: ":8080", Handler: wsContainer}

	go func() {
		log.Fatal(server.ListenAndServe())
	}()
}

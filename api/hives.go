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
	_ "log"
	"net/http"
	_ "strconv"

	"github.com/muesli/beehive/bees"

	"github.com/emicklei/go-restful"
	_ "github.com/emicklei/go-restful/swagger"
)

// POST http://localhost:8080/v1/hives
//
// GET http://localhost:8080/v1/hives/1
//
// PUT http://localhost:8080/v1/hives/1
//
// DELETE http://localhost:8080/v1/hives/1
//

type Hive struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Image       string `json:"image"`
	Options     []bees.BeeOptionDescriptor
	Events      []bees.EventDescriptor
	Actions     []bees.ActionDescriptor
}

type Hives struct {
	Hives []Hive
}

type HiveResource struct {
}

func (r HiveResource) Register(container *restful.Container) {
	ws := new(restful.WebService)
	ws.Path("/v1/hives").
		Doc("Manage Hives").
		Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON)

	ws.Route(ws.GET("").To(r.findHives).
		// docs
		Doc("get all hives").
		Operation("findHives").
		Writes(Hive{}))

	container.Add(ws)
}

// GET http://localhost:8080/v1/hives
//
func (r HiveResource) findHives(request *restful.Request, response *restful.Response) {
	hives := bees.GetBeeFactories()
	if len(hives) == 0 {
		response.AddHeader("Content-Type", "text/plain")
		response.WriteErrorString(http.StatusNotFound, "404: No hives could be found.")
		return
	}

	res := Hives{}
	for _, hive := range hives {
		h := Hive{
			Id:          (*hive).Name(),
			Name:        (*hive).Name(),
			Description: (*hive).Description(),
			Image:       "/images/" + (*hive).Image(),
			Options:     (*hive).Options(),
			Events:      (*hive).Events(),
			Actions:     (*hive).Actions(),
		}

		res.Hives = append(res.Hives, h)
	}
	response.WriteEntity(res)
}

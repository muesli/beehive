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

// POST http://localhost:8080/v1/bees
//
// GET http://localhost:8080/v1/bees/1
//
// PUT http://localhost:8080/v1/bees/1
//
// DELETE http://localhost:8080/v1/bees/1
//

type Bee struct {
	Id, Name string
}

type BeeResource struct {
}

func (r BeeResource) Register(container *restful.Container) {
	ws := new(restful.WebService)
	ws.Path("/v1/bees").
		Doc("Manage Bees").
		Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON)

	ws.Route(ws.GET("/{bee-id}").To(r.findBee).
		// docs
		Doc("get a bee").
		Operation("findBee").
		Param(ws.PathParameter("bee-id", "identifier of the bee").DataType("string")).
		Writes(Bee{}))

	ws.Route(ws.PUT("/{bee-id}").To(r.updateBee).
		// docs
		Doc("update a bee").
		Operation("updateBee").
		Param(ws.PathParameter("bee-id", "identifier of the bee").DataType("string")).
		ReturnsError(409, "duplicate bee-id", nil).
		Reads(Bee{}))

	ws.Route(ws.POST("").To(r.createBee).
		// docs
		Doc("create a bee").
		Operation("createBee").
		Reads(Bee{}))

	ws.Route(ws.DELETE("/{bee-id}").To(r.removeBee).
		// docs
		Doc("delete a bee").
		Operation("removeBee").
		Param(ws.PathParameter("bee-id", "identifier of the bee").DataType("string")))

	container.Add(ws)
}

// GET http://localhost:8080/v1/bees/1
//
func (r BeeResource) findBee(request *restful.Request, response *restful.Response) {
	id := request.PathParameter("bee-id")

	b := bees.GetBee(id)
	if b == nil {
		response.AddHeader("Content-Type", "text/plain")
		response.WriteErrorString(http.StatusNotFound, "404: Bee could not be found.")
		return
	}
	response.WriteEntity(b)
}

// POST http://localhost:8080/v1/bees
//
func (r *BeeResource) createBee(request *restful.Request, response *restful.Response) {
	b := new(Bee)
	err := request.ReadEntity(b)
	if err != nil {
		response.AddHeader("Content-Type", "text/plain")
		response.WriteErrorString(http.StatusInternalServerError, err.Error())
		return
	}
	/*	b.Id = strconv.Itoa(len(r.bees) + 1) // simple id generation
		r.bees[b.Id] = *b
		response.WriteHeader(http.StatusCreated)
		response.WriteEntity(b)*/
}

// PUT http://localhost:8080/v1/bees/1
//
func (r *BeeResource) updateBee(request *restful.Request, response *restful.Response) {
	b := new(Bee)
	err := request.ReadEntity(&b)
	if err != nil {
		response.AddHeader("Content-Type", "text/plain")
		response.WriteErrorString(http.StatusInternalServerError, err.Error())
		return
	}
	/*	r.bees[b.Id] = *b
		response.WriteEntity(b)*/
}

// DELETE http://localhost:8080/v1/bees/1
//
func (r *BeeResource) removeBee(request *restful.Request, response *restful.Response) {
	//	id := request.PathParameter("bee-id")
	//	delete(r.bees, id)
}

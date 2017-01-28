/*
 *    Copyright (C) 2017 Christian Muehlhaeuser
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

package bees

import (
	"net/http"

	"github.com/emicklei/go-restful"
	"github.com/muesli/beehive/bees"
	"github.com/muesli/smolder"
)

// BeePostStruct holds all values of an incoming POST request
type BeePostStruct struct {
	Bee struct {
		Description string          `json:"description"`
		Active      bool            `json:"active"`
		Options     bees.BeeOptions `json:"options"`
	} `json:"bee"`
}

// BeePutStruct holds all values of an incoming PUT request
type BeePutStruct struct {
	BeePostStruct
}

// PutAuthRequired returns true because all requests need authentication
func (r *BeeResource) PutAuthRequired() bool {
	return false
}

// PutDoc returns the description of this API endpoint
func (r *BeeResource) PutDoc() string {
	return "update an existing bee"
}

// PutParams returns the parameters supported by this API endpoint
func (r *BeeResource) PutParams() []*restful.Parameter {
	return nil
}

// Put processes an incoming PUT (update) request
func (r *BeeResource) Put(context smolder.APIContext, request *restful.Request, response *restful.Response) {
	resp := BeeResponse{}
	resp.Init(context)

	pps := BeePutStruct{}
	err := request.ReadEntity(&pps)
	if err != nil {
		smolder.ErrorResponseHandler(request, response, smolder.NewErrorResponse(
			http.StatusBadRequest,
			false,
			"Can't parse PUT data",
			"BeeResource PUT"))
		return
	}

	id := request.PathParameter("bee-id")
	bee := bees.GetBee(id)
	if bee == nil {
		r.NotFound(request, response)
		return
	}

	(*bee).SetDescription(pps.Bee.Description)
	(*bee).SetOptions(pps.Bee.Options)

	if pps.Bee.Active {
		bees.RestartBee(bee)
	} else {
		(*bee).Stop()
	}

	resp.AddBee(bee)
	resp.Send(response)
}

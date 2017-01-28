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
	"github.com/muesli/beehive/bees"

	"github.com/emicklei/go-restful"
	"github.com/muesli/smolder"
)

// GetAuthRequired returns true because all requests need authentication
func (r *BeeResource) GetAuthRequired() bool {
	return false
}

// GetByIDsAuthRequired returns true because all requests need authentication
func (r *BeeResource) GetByIDsAuthRequired() bool {
	return false
}

// GetDoc returns the description of this API endpoint
func (r *BeeResource) GetDoc() string {
	return "retrieve bees"
}

// GetParams returns the parameters supported by this API endpoint
func (r *BeeResource) GetParams() []*restful.Parameter {
	params := []*restful.Parameter{}
	// params = append(params, restful.QueryParameter("user_id", "id of a user").DataType("int64"))

	return params
}

// GetByIDs sends out all items matching a set of IDs
func (r *BeeResource) GetByIDs(ctx smolder.APIContext, request *restful.Request, response *restful.Response, ids []string) {
	resp := BeeResponse{}
	resp.Init(ctx)

	for _, id := range ids {
		bee := bees.GetBee(id)
		if bee == nil {
			r.NotFound(request, response)
			return
		}

		resp.AddBee(bee)
	}

	resp.Send(response)
}

// Get sends out items matching the query parameters
func (r *BeeResource) Get(ctx smolder.APIContext, request *restful.Request, response *restful.Response, params map[string][]string) {
	//	ctxapi := ctx.(*context.APIContext)
	bees := bees.GetBees()
	if len(bees) == 0 { // err != nil {
		r.NotFound(request, response)
		return
	}

	resp := BeeResponse{}
	resp.Init(ctx)

	for _, bee := range bees {
		resp.AddBee(bee)
	}

	resp.Send(response)
}

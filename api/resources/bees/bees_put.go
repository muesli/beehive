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
	"github.com/emicklei/go-restful"
	"github.com/muesli/beehive/bees"
	"github.com/muesli/beehive/cfg"
	"github.com/muesli/smolder"
)

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
func (r *BeeResource) Put(context smolder.APIContext, data interface{}, request *restful.Request, response *restful.Response) {
	resp := BeeResponse{}
	resp.Init(context)

	pps := data.(*BeePostStruct)
	id := request.PathParameter("bee-id")
	bee := bees.GetBee(id)
	if bee == nil {
		r.NotFound(request, response)
		return
	}

	(*bee).SetDescription(pps.Bee.Description)
	(*bee).ReloadOptions(pps.Bee.Options)

	if pps.Bee.Active {
		bees.RestartBee(bee)
	} else {
		(*bee).Stop()
	}

	cfg.SaveCurrentConfig()
	resp.AddBee(bee)
	resp.Send(response)
}

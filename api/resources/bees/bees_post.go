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
	"github.com/muesli/smolder"
)

// BeePostStruct holds all values of an incoming POST request
type BeePostStruct struct {
	Bee struct {
		Name        string          `json:"name"`
		Namespace   string          `json:"namespace"`
		Description string          `json:"description"`
		Active      bool            `json:"active"`
		Options     bees.BeeOptions `json:"options"`
	} `json:"bee"`
}

// PostAuthRequired returns true because all requests need authentication
func (r *BeeResource) PostAuthRequired() bool {
	return false
}

// PostDoc returns the description of this API endpoint
func (r *BeeResource) PostDoc() string {
	return "create a new bee"
}

// PostParams returns the parameters supported by this API endpoint
func (r *BeeResource) PostParams() []*restful.Parameter {
	return nil
}

// Post processes an incoming POST (create) request
func (r *BeeResource) Post(context smolder.APIContext, data interface{}, request *restful.Request, response *restful.Response) {
	resp := BeeResponse{}
	resp.Init(context)

	pps := data.(BeePostStruct)
	c, err := bees.NewBeeConfig(pps.Bee.Name, pps.Bee.Namespace, pps.Bee.Description, pps.Bee.Options)
	if err != nil {
		smolder.ErrorResponseHandler(request, response, smolder.NewErrorResponse(
			422, // Go 1.7+: http.StatusUnprocessableEntity,
			false,
			err,
			"BeeResource POST"))
		return
	}

	bee := bees.StartBee(c)
	resp.AddBee(bee)

	resp.Send(response)
}

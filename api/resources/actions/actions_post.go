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

package actions

import (
	"net/http"

	"github.com/emicklei/go-restful"
	"github.com/muesli/beehive/bees"
	"github.com/muesli/smolder"
)

// ActionPostStruct holds all values of an incoming POST request
type ActionPostStruct struct {
	Action struct {
		Bee     string                `json:"bee"`
		Name    string                `json:"name"`
		Options bees.PlaceholderSlice `json:"options"`
	} `json:"action"`
}

// PostAuthRequired returns true because all requests need authentication
func (r *ActionResource) PostAuthRequired() bool {
	return false
}

// PostDoc returns the description of this API endpoint
func (r *ActionResource) PostDoc() string {
	return "create a new action"
}

// PostParams returns the parameters supported by this API endpoint
func (r *ActionResource) PostParams() []*restful.Parameter {
	return nil
}

// Post processes an incoming POST (create) request
func (r *ActionResource) Post(context smolder.APIContext, request *restful.Request, response *restful.Response) {
	resp := ActionResponse{}
	resp.Init(context)

	pps := ActionPostStruct{}
	err := request.ReadEntity(&pps)
	if err != nil {
		smolder.ErrorResponseHandler(request, response, smolder.NewErrorResponse(
			http.StatusBadRequest,
			false,
			"Can't parse POST data",
			"ActionResource POST"))
		return
	}

	action := bees.Action{
		ID:      bees.UUID(),
		Bee:     pps.Action.Bee,
		Name:    pps.Action.Name,
		Options: pps.Action.Options,
	}
	actions := append(bees.GetActions(), action)
	bees.SetActions(actions)

	resp.AddAction(&action)
	resp.Send(response)
}

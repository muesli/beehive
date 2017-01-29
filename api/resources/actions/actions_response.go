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

package actions

import (
	"github.com/muesli/beehive/bees"

	"github.com/muesli/smolder"
)

// ActionResponse is the common response to 'action' requests
type ActionResponse struct {
	smolder.Response

	Actions []actionInfoResponse `json:"actions,omitempty"`
	actions []*bees.Action
}

type actionInfoResponse struct {
	ID      string            `json:"id"`
	Bee     string            `json:"bee"`
	Name    string            `json:"name"`
	Options bees.Placeholders `json:"options"`
}

// Init a new response
func (r *ActionResponse) Init(context smolder.APIContext) {
	r.Parent = r
	r.Context = context

	r.Actions = []actionInfoResponse{}
}

// AddAction adds a action to the response
func (r *ActionResponse) AddAction(action *bees.Action) {
	r.actions = append(r.actions, action)
	r.Actions = append(r.Actions, prepareActionResponse(r.Context, action))
}

// EmptyResponse returns an empty API response for this endpoint if there's no data to respond with
func (r *ActionResponse) EmptyResponse() interface{} {
	if len(r.actions) == 0 {
		var out struct {
			Actions interface{} `json:"actions"`
		}
		out.Actions = []actionInfoResponse{}
		return out
	}
	return nil
}

func prepareActionResponse(context smolder.APIContext, action *bees.Action) actionInfoResponse {
	//	ctx := context.(*context.APIContext)
	resp := actionInfoResponse{
		ID:      (*action).ID,
		Bee:     (*action).Bee,
		Name:    (*action).Name,
		Options: (*action).Options,
	}

	return resp
}

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

package chains

import (
	"github.com/emicklei/go-restful"
	"github.com/muesli/beehive/bees"
	"github.com/muesli/smolder"
)

// DeleteAuthRequired returns true because all requests need authentication
func (r *ChainResource) DeleteAuthRequired() bool {
	return false
}

// DeleteDoc returns the description of this API endpoint
func (r *ChainResource) DeleteDoc() string {
	return "delete a chain"
}

// DeleteParams returns the parameters supported by this API endpoint
func (r *ChainResource) DeleteParams() []*restful.Parameter {
	return nil
}

// Delete processes an incoming DELETE request
func (r *ChainResource) Delete(context smolder.APIContext, request *restful.Request, response *restful.Response) {
	resp := ChainResponse{}
	resp.Init(context)

	id := request.PathParameter("chain-id")

	found := false
	chains := []bees.Chain{}
	actions := []bees.Action{}
	actionsList := []string{}

	for _, v := range bees.GetChains() {
		if v.Name == id {
			found = true
			actionsList = v.Actions
		} else {
			chains = append(chains, v)
		}
	}

	if found {
		for _, k := range bees.GetActions() {
			for _, t := range actionsList {
				// Delete the Action belongs to the chain
				if t == k.ID {
					continue
				} else {
					actions = append(actions, k)
				}
			}
		}
		bees.SetChains(chains)
		bees.SetActions(actions)
		resp.Send(response)
	} else {
		r.NotFound(request, response)
	}
}

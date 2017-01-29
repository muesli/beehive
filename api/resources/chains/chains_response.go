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

package chains

import (
	"sort"

	restful "github.com/emicklei/go-restful"
	"github.com/muesli/beehive/bees"

	"github.com/muesli/smolder"
)

// ChainResponse is the common response to 'chain' requests
type ChainResponse struct {
	smolder.Response

	Chains []chainInfoResponse `json:"chains,omitempty"`
	chains map[string]*bees.Chain
}

type chainInfoResponse struct {
	ID          string      `json:"id"`
	Name        string      `json:"name"`
	Description string      `json:"description"`
	Event       *bees.Event `json:"event"`
	Filters     []string    `json:"filters,omitempty"`
	Actions     []string    `json:"actions"`
}

// Init a new response
func (r *ChainResponse) Init(context smolder.APIContext) {
	r.Parent = r
	r.Context = context

	r.chains = make(map[string]*bees.Chain)
}

// AddChain adds a chain to the response
func (r *ChainResponse) AddChain(chain bees.Chain) {
	r.chains[chain.Name] = &chain
}

// Send responds to a request with http.StatusOK
func (r *ChainResponse) Send(response *restful.Response) {
	var keys []string
	for k := range r.chains {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		r.Chains = append(r.Chains, prepareChainResponse(r.Context, r.chains[k]))
	}

	r.Response.Send(response)
}

// EmptyResponse returns an empty API response for this endpoint if there's no data to respond with
func (r *ChainResponse) EmptyResponse() interface{} {
	if len(r.chains) == 0 {
		var out struct {
			Chains interface{} `json:"chains"`
		}
		out.Chains = []chainInfoResponse{}
		return out
	}
	return nil
}

func prepareChainResponse(context smolder.APIContext, chain *bees.Chain) chainInfoResponse {
	//	ctx := context.(*context.APIContext)
	resp := chainInfoResponse{
		ID:          (*chain).Name,
		Name:        (*chain).Name,
		Description: (*chain).Description,
		Event:       (*chain).Event,
		Actions:     (*chain).Actions,
		Filters:     (*chain).Filters,
	}

	return resp
}

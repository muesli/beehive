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
	"sort"
	"time"

	restful "github.com/emicklei/go-restful"
	"github.com/muesli/smolder"

	"github.com/muesli/beehive/api/resources/hives"
	"github.com/muesli/beehive/bees"
)

// BeeResponse is the common response to 'bee' requests
type BeeResponse struct {
	smolder.Response

	Bees []beeInfoResponse `json:"bees,omitempty"`
	bees map[string]*bees.BeeInterface

	Hives []hives.HiveInfoResponse `json:"hives,omitempty"`
	hives map[string]*bees.BeeFactoryInterface
}

type beeInfoResponse struct {
	ID          string           `json:"id"`
	Name        string           `json:"name"`
	Namespace   string           `json:"namespace"`
	Description string           `json:"description"`
	LastAction  time.Time        `json:"lastaction"`
	LastEvent   time.Time        `json:"lastevent"`
	Active      bool             `json:"active"`
	Options     []bees.BeeOption `json:"options"`
}

// Init a new response
func (r *BeeResponse) Init(context smolder.APIContext) {
	r.Parent = r
	r.Context = context

	r.bees = make(map[string]*bees.BeeInterface)
	r.hives = make(map[string]*bees.BeeFactoryInterface)
}

// AddBee adds a bee to the response
func (r *BeeResponse) AddBee(bee *bees.BeeInterface) {
	r.bees[(*bee).Name()] = bee

	hive := bees.GetFactory((*bee).Namespace())
	if hive == nil {
		panic("Hive for Bee not found")
	}

	r.hives[(*hive).Name()] = hive
	r.Hives = append(r.Hives, hives.PrepareHiveResponse(r.Context, hive))
}

// Send responds to a request with http.StatusOK
func (r *BeeResponse) Send(response *restful.Response) {
	var keys []string
	for k := range r.bees {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		r.Bees = append(r.Bees, prepareBeeResponse(r.Context, r.bees[k]))
	}

	r.Response.Send(response)
}

// EmptyResponse returns an empty API response for this endpoint if there's no data to respond with
func (r *BeeResponse) EmptyResponse() interface{} {
	if len(r.bees) == 0 {
		var out struct {
			Bees interface{} `json:"bees"`
		}
		out.Bees = []beeInfoResponse{}
		return out
	}
	return nil
}

func prepareBeeResponse(context smolder.APIContext, bee *bees.BeeInterface) beeInfoResponse {
	//	ctx := context.(*context.APIContext)
	resp := beeInfoResponse{
		ID:          (*bee).Name(),
		Name:        (*bee).Name(),
		Namespace:   (*bee).Namespace(),
		Description: (*bee).Description(),
		LastAction:  (*bee).LastAction(),
		LastEvent:   (*bee).LastEvent(),
		Active:      (*bee).IsRunning(),
		Options:     (*bee).Options(),
	}

	return resp
}

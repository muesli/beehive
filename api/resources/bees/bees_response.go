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
	"time"

	"github.com/muesli/beehive/bees"

	"github.com/muesli/smolder"
)

// BeeResponse is the common response to 'bee' requests
type BeeResponse struct {
	smolder.Response

	Bees []beeInfoResponse `json:"bees,omitempty"`
	bees []*bees.BeeInterface
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

	r.Bees = []beeInfoResponse{}
}

// AddBee adds a bee to the response
func (r *BeeResponse) AddBee(bee *bees.BeeInterface) {
	r.bees = append(r.bees, bee)
	r.Bees = append(r.Bees, prepareBeeResponse(r.Context, bee))
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

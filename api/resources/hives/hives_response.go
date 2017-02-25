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

package hives

import (
	"net/url"
	"path"
	"sort"

	restful "github.com/emicklei/go-restful"
	"github.com/muesli/smolder"

	"github.com/muesli/beehive/api/context"
	"github.com/muesli/beehive/bees"
)

// HiveResponse is the common response to 'hive' requests
type HiveResponse struct {
	smolder.Response

	Hives []hiveInfoResponse `json:"hives,omitempty"`
	hives map[string]*bees.BeeFactoryInterface
}

type hiveInfoResponse struct {
	ID          string                     `json:"id"`
	Name        string                     `json:"name"`
	Description string                     `json:"description"`
	Image       string                     `json:"image"`
	LogoColor   string                     `json:"logocolor"`
	Options     []bees.BeeOptionDescriptor `json:"options"`
	Events      []bees.EventDescriptor     `json:"events"`
	Actions     []bees.ActionDescriptor    `json:"actions"`
}

// Init a new response
func (r *HiveResponse) Init(context smolder.APIContext) {
	r.Parent = r
	r.Context = context

	r.hives = make(map[string]*bees.BeeFactoryInterface)
}

// AddHive adds a hive to the response
func (r *HiveResponse) AddHive(hive *bees.BeeFactoryInterface) {
	r.hives[(*hive).Name()] = hive
}

// Send responds to a request with http.StatusOK
func (r *HiveResponse) Send(response *restful.Response) {
	var keys []string
	for k := range r.hives {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		r.Hives = append(r.Hives, prepareHiveResponse(r.Context, r.hives[k]))
	}

	r.Response.Send(response)
}

// EmptyResponse returns an empty API response for this endpoint if there's no data to respond with
func (r *HiveResponse) EmptyResponse() interface{} {
	if len(r.hives) == 0 {
		var out struct {
			Hives interface{} `json:"hives"`
		}
		out.Hives = []hiveInfoResponse{}
		return out
	}
	return nil
}

func prepareHiveResponse(ctx smolder.APIContext, hive *bees.BeeFactoryInterface) hiveInfoResponse {
	u, _ := url.Parse(ctx.(*context.APIContext).Config.BaseURL)
	u.Path = path.Join(u.Path, "images", (*hive).Image())

	resp := hiveInfoResponse{
		ID:          (*hive).ID(),
		Name:        (*hive).Name(),
		Description: (*hive).Description(),
		Image:       u.String(),
		LogoColor:   (*hive).LogoColor(),
		Options:     (*hive).Options(),
		Events:      (*hive).Events(),
		Actions:     (*hive).Actions(),
	}

	return resp
}

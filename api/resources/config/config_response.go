/*
 *    Copyright (C) 2019 Sergio Rubio
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
 *      Sergio Rubio <sergio@rubio.im>
 */

package config

import (
	restful "github.com/emicklei/go-restful"
	"github.com/muesli/beehive/cfg"
	"github.com/muesli/smolder"
)

// HiveResponse is the common response to 'hive' requests
type ConfigResponse struct {
	smolder.Response

	Config cfg.Config `json:"config,omitempty"`
}

// Init a new response
func (r *ConfigResponse) Init(context smolder.APIContext) {
	r.Parent = r
	r.Context = context
}

// Send responds to a request with http.StatusOK
func (r *ConfigResponse) Send(response *restful.Response) {
	r.Response.Send(response)
}

// EmptyResponse returns an empty API response for this endpoint if there's no data to respond with
func (r *ConfigResponse) EmptyResponse() interface{} {
	return nil
}

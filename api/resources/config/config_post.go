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
	"github.com/emicklei/go-restful"
	"github.com/muesli/beehive/cfg"
	"github.com/muesli/smolder"
)

// ConfigPostStruct holds all values of an incoming POST request
type ConfigPostStruct struct {
	Config struct {
		Name string `json:"name"`
	} `json:"config"`
}

// PostAuthRequired returns true because all requests need authentication
func (r *ConfigResource) PostAuthRequired() bool {
	return false
}

// PostDoc returns the description of this API endpoint
func (r *ConfigResource) PostDoc() string {
	return "create a new config"
}

// PostParams returns the parameters supported by this API endpoint
func (r *ConfigResource) PostParams() []*restful.Parameter {
	return nil
}

// Post processes an incoming POST (create) request
func (r *ConfigResource) Post(context smolder.APIContext, data interface{}, request *restful.Request, response *restful.Response) {
	cfg.SaveCurrentConfig()
	resp := ConfigResponse{}
	resp.Init(context)
	resp.Send(response)
}

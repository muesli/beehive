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
 *      Sergio Rubio <rubiojr@frameos.org>
 */

package config

import (
	"github.com/emicklei/go-restful"
	"github.com/muesli/beehive/bees"
	"github.com/muesli/beehive/cfg"
	"github.com/muesli/smolder"
)

// GetAuthRequired returns true because all requests need authentication
func (r *ConfigResource) GetAuthRequired() bool {
	return false
}

// GetByIDsAuthRequired returns true because all requests need authentication
func (r *ConfigResource) GetByIDsAuthRequired() bool {
	return false
}

// GetDoc returns the description of this API endpoint
func (r *ConfigResource) GetDoc() string {
	return "retrieve config"
}

// GetParams returns the parameters supported by this API endpoint
func (r *ConfigResource) GetParams() []*restful.Parameter {
	params := []*restful.Parameter{}

	return params
}

// Get sends out items matching the query parameters
func (r *ConfigResource) Get(ctx smolder.APIContext, request *restful.Request, response *restful.Response, params map[string][]string) {
	resp := ConfigResponse{}
	resp.Init(ctx)

	config := cfg.GetConfig()
	config.Bees = bees.BeeConfigs()
	config.Chains = bees.GetChains()
	config.Actions = bees.GetActions()
	resp.Config = config
	resp.Send(response)
}

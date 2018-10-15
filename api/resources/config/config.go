/*
 *    Copyright (C) 2018 Sergio Rubio
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
	"github.com/muesli/smolder"
)

// ChainResource is the resource responsible for /chains
type ConfigResource struct {
	smolder.Resource
}

var (
	_ smolder.PostSupported = &ConfigResource{}
	_ smolder.GetSupported  = &ConfigResource{}
)

// Register this resource with the container to setup all the routes
func (r *ConfigResource) Register(container *restful.Container, config smolder.APIConfig, context smolder.APIContextFactory) {
	r.Name = "ConfigResource"
	r.TypeName = "config"
	r.Endpoint = "config"
	r.Doc = "Manage config"

	r.Config = config
	r.Context = context

	r.Init(container, r)
}

// Reads returns the model that will be read by POST, PUT & PATCH operations
func (r *ConfigResource) Reads() interface{} {
	return &ConfigPostStruct{}
}

// Returns returns the model that will be returned
func (r *ConfigResource) Returns() interface{} {
	return ConfigResponse{}
}

// Validate checks an incoming request for data errors
func (r *ConfigResource) Validate(context smolder.APIContext, data interface{}, request *restful.Request) error {
	return nil
}

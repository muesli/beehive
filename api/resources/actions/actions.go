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

package actions

import (
	"github.com/emicklei/go-restful"
	"github.com/muesli/smolder"
)

// ActionResource is the resource responsible for /actions
type ActionResource struct {
	smolder.Resource
}

var (
	_ smolder.GetIDSupported = &ActionResource{}
	_ smolder.GetSupported   = &ActionResource{}
	_ smolder.PostSupported  = &ActionResource{}
)

// Register this resource with the container to setup all the routes
func (r *ActionResource) Register(container *restful.Container, config smolder.APIConfig, context smolder.APIContextFactory) {
	r.Name = "ActionResource"
	r.TypeName = "action"
	r.Endpoint = "actions"
	r.Doc = "Manage actions"

	r.Config = config
	r.Context = context

	r.Init(container, r)
}

// Reads returns the model that will be read by POST, PUT & PATCH operations
func (r *ActionResource) Reads() interface{} {
	return ActionPostStruct{}
}

// Returns returns the model that will be returned
func (r *ActionResource) Returns() interface{} {
	return ActionResponse{}
}

// Validate checks an incoming request for data errors
func (r *ActionResource) Validate(context smolder.APIContext, data interface{}, request *restful.Request) error {
	//	ps := data.(*ActionPostStruct)
	// FIXME
	return nil
}

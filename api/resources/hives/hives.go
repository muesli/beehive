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

package hives

import (
	"github.com/emicklei/go-restful"
	"github.com/muesli/smolder"
)

// HiveResource is the resource responsible for /hives
type HiveResource struct {
	smolder.Resource
}

var (
	_ smolder.GetIDSupported = &HiveResource{}
	_ smolder.GetSupported   = &HiveResource{}
)

// Register this resource with the container to setup all the routes
func (r *HiveResource) Register(container *restful.Container, config smolder.APIConfig, context smolder.APIContextFactory) {
	r.Name = "HiveResource"
	r.TypeName = "hive"
	r.Endpoint = "hives"
	r.Doc = "Manage hives"

	r.Config = config
	r.Context = context

	r.Init(container, r)
}

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

package context

import (
	"strings"

	restful "github.com/emicklei/go-restful"
	"github.com/muesli/smolder"
)

// APIContext is polly's central context
type APIContext struct {
	Config smolder.APIConfig
}

// NewAPIContext returns a new polly context
func (context *APIContext) NewAPIContext() smolder.APIContext {
	ctx := &APIContext{
		Config: context.Config,
	}
	return ctx
}

// Authentication parses the request for an access-/authtoken and returns the matching user
func (context *APIContext) Authentication(request *restful.Request) (interface{}, error) {
	//FIXME: implement this properly

	t := request.QueryParameter("accesstoken")
	if len(t) == 0 {
		t = request.HeaderParameter("authorization")
		if strings.Index(t, " ") > 0 {
			t = strings.TrimSpace(strings.Split(t, " ")[1])
		}
	}

	return nil, nil // context.GetUserByAccessToken(t)
}

// LogSummary logs out the current context stats
func (context *APIContext) LogSummary() {
}

/*
 *    Copyright (C) 2014 Daniel 'grindhold' Brendle
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
 *      Daniel 'grindhold' Brendle <grindhold@skarphed.org>
 */

/*
   Please note that, in order to run this bee on a nagios-server, you
   have to provide the nagios status-script found at

   https://github.com/lizell/php-nagios-json/blob/master/statusJson.php

   just drop this script in the htdocs-folder of your nagios-installation
   and change the variable $statusFile to where the status.dat-file of your
   installation resides
*/

package nagiosbee

import (
	"github.com/muesli/beehive/modules"
)

type NagiosBee struct {
	name        string
	namespace   string
	description string

	url      string
	user     string
	password string
	services []Service

	eventChan chan modules.Event
}

type Service struct {
	Name      string
	Status    int
	LastEvent int
}

func (mod *NagiosBee) Name() string {
	return mod.name
}

func (mod *NagiosBee) Action(action modules.Action) []modules.Placeholder {
	return []modules.Placeholder{}
}

func (mod *NagiosBee) Run(cin chan modules.Event) {
	return
}

func (mod *NagiosBee) Namespace() string {
	return mod.namespace
}

func (mod *NagiosBee) Description() string {
	return mod.description
}

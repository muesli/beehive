/*
 *    Copyright (C) 2014 Christian Muehlhaeuser
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

// beehive's Anel PowerCtrl module.
package anelpowerctrlbee

import (
	"github.com/muesli/beehive/modules"
)

type AnelPowerCtrlBee struct {
	name        string
	namespace   string
	description string

	addr        string
	user        string
	password    string

	eventChan chan modules.Event
}

func (mod *AnelPowerCtrlBee) Name() string {
	return mod.name
}

func (mod *AnelPowerCtrlBee) Namespace() string {
	return mod.namespace
}

func (mod *AnelPowerCtrlBee) Description() string {
	return mod.description
}

func (mod *AnelPowerCtrlBee) Run(cin chan modules.Event) {
	mod.eventChan = cin
}

func (mod *AnelPowerCtrlBee) Action(action modules.Action) []modules.Placeholder {
	outs := []modules.Placeholder{}
	return outs
}

/*
 *    Copyright (C) 2016 Gonzalo Izquierdo
 *                  2017 Christian Muehlhaeuser
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
 *      Gonzalo Izquierdo <lalotone@gmail.com>
 *      Christian Muehlhaeuser <muesli@gmail.com>
 */

// Package transmissionbee is a Bee that can send torrents to Transmission.
package transmissionbee

import (
	"github.com/kr/pretty"
	"github.com/odwrtw/transmission"

	"github.com/muesli/beehive/bees"
)

// TransmissionBee is a Bee that can send torrents to Transmission.
type TransmissionBee struct {
	bees.Bee

	client *transmission.Client
}

// Action triggers the action passed to it.
func (mod *TransmissionBee) Action(action bees.Action) []bees.Placeholder {
	outs := []bees.Placeholder{}
	switch action.Name {
	case "add_torrent":
		torrentMsg := ""
		action.Options.Bind("torrent", &torrentMsg)

		_, err := mod.client.Add(torrentMsg)
		if err != nil {
			mod.LogErrorf("Error adding torrent/magnet: %s", err)
		}
	}
	return outs
}

// ReloadOptions parses the config options and initializes the Bee.
func (mod *TransmissionBee) ReloadOptions(options bees.BeeOptions) {
	mod.SetOptions(options)

	conf := transmission.Config{}
	options.Bind("url", &conf.Address)
	options.Bind("username", &conf.User)
	options.Bind("password", &conf.Password)

	t, err := transmission.New(conf)
	if err != nil {
		pretty.Println(err)
	}
	mod.client = t
}

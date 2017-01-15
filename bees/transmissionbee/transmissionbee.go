/*
 *    Copyright (C) 2016 Gonzalo Izquierdo
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
 */

package transmissionbee

import (
	"strings"

	"github.com/muesli/beehive/bees"
	"github.com/odwrtw/transmission"
)

type TransmissionBee struct {
	bees.Bee

	client *transmission.Client
}

func (mod *TransmissionBee) Action(action bees.Action) []bees.Placeholder {
	outs := []bees.Placeholder{}
	switch action.Name {
	case "add-torrent":
		torrentMsg := ""
		commandPrefix := ""
		action.Options.Bind("torrent", &torrentMsg)
		action.Options.Bind("commandPrefix", &commandPrefix)

		torrentMsg = strings.TrimSpace(strings.Replace(torrentMsg, commandPrefix, "", 1))
		_, err := mod.client.Add(torrentMsg)
		if err != nil {
			panic("Transmission: error adding torrent/magnet")
		}
	}
	return outs
}

func (mod *TransmissionBee) Run(eventChan chan bees.Event) {
}

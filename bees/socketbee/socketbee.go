/*
 *    Copyright (C) 2019 Christian Muehlhaeuser
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

// Package socketbee is a Bee that lets you transmit data via UDP sockets.
package socketbee

import (
	"log"
	"net"
	"strconv"

	"github.com/muesli/beehive/bees"
)

// SocketBee is a Bee that lets you transmit data via UDP sockets.
type SocketBee struct {
	bees.Bee

	eventChan chan bees.Event
}

// Run executes the Bee's event loop.
func (mod *SocketBee) Run(cin chan bees.Event) {
	mod.eventChan = cin

	select {
	case <-mod.SigChan:
		return
	}
}

// Action triggers the action passed to it.
func (mod *SocketBee) Action(action bees.Action) []bees.Placeholder {
	outs := []bees.Placeholder{}

	var data string
	var addr string
	var port int

	action.Options.Bind("address", &addr)
	action.Options.Bind("port", &port)
	action.Options.Bind("data", &data)

	switch action.Name {
	case "send":
		// log.Println("Sending", data, "to", addr, port)

		sa, err := net.ResolveUDPAddr("udp", addr+":"+strconv.Itoa(port))
		if err != nil {
			log.Panicln(err)
		}

		conn, err := net.DialUDP("udp", nil, sa)
		if err != nil {
			log.Panicln(err)
		}

		defer conn.Close()
		_, err = conn.Write([]byte(data))
		if err != nil {
			log.Panicln(err)
		}

	default:
		panic("Unknown action triggered in " + mod.Name() + ": " + action.Name)
	}

	return outs
}

// ReloadOptions parses the config options and initializes the Bee.
func (mod *SocketBee) ReloadOptions(options bees.BeeOptions) {
	mod.SetOptions(options)
}

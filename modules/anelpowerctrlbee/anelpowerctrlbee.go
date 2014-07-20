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
	"log"
	"net"
	"strconv"
	"time"
)

type AnelPowerCtrlBee struct {
	modules.Module

	addr     string
	user     string
	password string
}

func (mod *AnelPowerCtrlBee) Run(cin chan modules.Event) {
}

func (mod *AnelPowerCtrlBee) anelSwitch(socket int, state bool) bool {
	conn, err := net.ListenUDP("udp", &net.UDPAddr{Port: 0})
	if err != nil {
		log.Fatal(err)
	}
	conn.SetDeadline(time.Now().Add(3 * time.Second))

	addr, err := net.ResolveUDPAddr("udp", mod.addr+":75")
	if err != nil {
		log.Fatal(err)
	}

	stateToken := "off"
	if state {
		stateToken = "on"
	}
	b := "Sw_" + stateToken + strconv.Itoa(socket) + mod.user + mod.password

	_, err = conn.WriteToUDP([]byte(b), addr)
	if err != nil {
		log.Fatal(err)
	}

	return true
}

func (mod *AnelPowerCtrlBee) Action(action modules.Action) []modules.Placeholder {
	outs := []modules.Placeholder{}

	switch action.Name {
	case "switch":
		socket := 0
		state := false

		for _, opt := range action.Options {
			if opt.Name == "socket" {
				socket = int(opt.Value.(float64))
			}
			if opt.Name == "state" {
				state = opt.Value.(bool)
			}
		}

		mod.anelSwitch(socket, state)

	default:
		panic("Unknown action triggered in " +mod.Name()+": "+action.Name)
	}

	return outs
}

/*
 *    Copyright (C) 2014 Daniel 'grindhold' Brendle
 *                  2014-2017 Christian Muehlhaeuser
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
 *      Christian Muehlhaeuser <muesli@gmail.com>
 */

package notificationbee

import (
	"log"
	"strings"

	"github.com/guelfey/go.dbus"
	"github.com/muesli/beehive/bees"
)

const (
	URGENCY_LOW      = uint32(iota)
	URGENCY_NORMAL   = uint32(iota)
	URGENCY_CRITICAL = uint32(iota)
)

var (
	urgency_map map[string]uint32 = map[string]uint32{
		"":         URGENCY_NORMAL,
		"normal":   URGENCY_NORMAL,
		"low":      URGENCY_LOW,
		"critical": URGENCY_CRITICAL,
	}
)

type NotificationBee struct {
	bees.Bee
	conn     *dbus.Conn
	notifier *dbus.Object
}

func (mod *NotificationBee) Run(cin chan bees.Event) {
	conn, err := dbus.SessionBus()
	mod.conn = conn
	if err != nil {
		panic(err)
	}
	mod.notifier = mod.conn.Object("org.freedesktop.Notifications", "/org/freedesktop/Notifications")
}

func (mod *NotificationBee) Action(action bees.Action) []bees.Placeholder {
	outs := []bees.Placeholder{}

	switch action.Name {
	case "notify":
		text := ""
		u := ""
		urgency := URGENCY_NORMAL

		action.Options.Bind("text", &text)
		action.Options.Bind("urgency", &u)
		text = strings.TrimSpace(text)
		urgency, _ = urgency_map[u]

		if len(text) > 0 {
			call := mod.notifier.Call("org.freedesktop.Notifications.Notify", 0, "", uint32(0),
				"", "Beehive", text, []string{},
				map[string]dbus.Variant{"urgency": dbus.MakeVariant(urgency)}, int32(5000))

			if call.Err != nil {
				log.Println("(" + string(urgency) + ") Failed to print message: " + text)
			}
		}

	default:
		panic("Unknown action triggered in " + mod.Name() + ": " + action.Name)
	}
	return outs
}

func (mod *NotificationBee) ReloadOptions(options bees.BeeOptions) {
	//FIXME: implement this
	mod.SetOptions(options)
}

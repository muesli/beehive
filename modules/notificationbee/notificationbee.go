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

package notificationbee

import (
	"github.com/guelfey/go.dbus"
	"github.com/muesli/beehive/modules"
	"log"
)

type NotificationBee struct {
	modules.Module
	conn     *dbus.Conn
	notifier *dbus.Object
}

func (mod *NotificationBee) Run(cin chan modules.Event) {
	conn, err := dbus.SessionBus()
	mod.conn = conn
	if err != nil {
		panic(err)
	}
	mod.notifier = mod.conn.Object("org.freedesktop.Notifications", "/org/freedesktop/Notifications")
}

func (mod *NotificationBee) Action(action modules.Action) []modules.Placeholder {
	outs := []modules.Placeholder{}
	switch action.Name {
	case "notify":
		text := ""
		urgency := ""
		for _, opt := range action.Options {
			if opt.Name == "text" {
				text = opt.Value.(string)
			}
			if opt.Name == "urgency" {
				urgency = opt.Value.(string)
			}
			call := mod.notifier.Call("org.freedesktop.Notifications.Notify", 1, "", uint32(0),
				"", "Beehive", text, []string{},
				map[string]dbus.Variant{}, int32(5000))
			if call.Err != nil {
				log.Println("(" + urgency + ") Failed to print message: " + text)
			}
		}
	default:
		return outs
	}
	return outs
}

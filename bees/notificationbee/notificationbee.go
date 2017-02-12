// +build dragonfly freebsd linux netbsd openbsd solaris

/*
 *    Copyright (C) 2014      Daniel 'grindhold' Brendle
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

// Package notificationbee is a Bee that can trigger freedesktop.org
// DBus notifications.
package notificationbee

import (
	"strings"

	"github.com/guelfey/go.dbus"

	"github.com/muesli/beehive/bees"
)

// Urgency level iota
const (
	UrgencyLow      = uint32(iota)
	UrgencyNormal   = uint32(iota)
	UrgencyCritical = uint32(iota)
)

var (
	urgencyMap map[string]uint32 = map[string]uint32{
		"":         UrgencyNormal,
		"normal":   UrgencyNormal,
		"low":      UrgencyLow,
		"critical": UrgencyCritical,
	}
)

// NotificationBee is a Bee that can trigger freedesktop.org DBus
// notifications.
type NotificationBee struct {
	bees.Bee
	conn     *dbus.Conn
	notifier *dbus.Object
}

// Run executes the Bee's event loop.
func (mod *NotificationBee) Run(cin chan bees.Event) {
	conn, err := dbus.SessionBus()
	mod.conn = conn
	if err != nil {
		panic(err)
	}
	mod.notifier = mod.conn.Object("org.freedesktop.Notifications", "/org/freedesktop/Notifications")
}

// Action triggers the action passed to it.
func (mod *NotificationBee) Action(action bees.Action) []bees.Placeholder {
	outs := []bees.Placeholder{}

	switch action.Name {
	case "notify":
		text := ""
		u := ""
		urgency := UrgencyNormal

		action.Options.Bind("text", &text)
		action.Options.Bind("urgency", &u)
		text = strings.TrimSpace(text)
		urgency, _ = urgencyMap[u]

		if len(text) > 0 {
			call := mod.notifier.Call("org.freedesktop.Notifications.Notify", 0, "", uint32(0),
				"", "Beehive", text, []string{},
				map[string]dbus.Variant{"urgency": dbus.MakeVariant(urgency)}, int32(5000))

			if call.Err != nil {
				mod.Logln("(" + string(urgency) + ") Failed to print message: " + text)
			}
		}

	default:
		panic("Unknown action triggered in " + mod.Name() + ": " + action.Name)
	}
	return outs
}

// ReloadOptions parses the config options and initializes the Bee.
func (mod *NotificationBee) ReloadOptions(options bees.BeeOptions) {
	mod.SetOptions(options)
}

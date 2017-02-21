// +build darwin

/*
 *    Copyright (C) 2014-2017 Christian Muehlhaeuser
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

// Package notificationbee is a Bee that can trigger desktop notifications.
package notificationbee

import "github.com/deckarep/gosx-notifier"

// Run executes the Bee's event loop.
func (mod *NotificationBee) execAction(text string, urgency uint32) {
	note := gosxnotifier.NewNotification(text)
	note.Subtitle = "this is subtitle"
	note.Push()
}

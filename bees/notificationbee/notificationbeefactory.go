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

package notificationbee

import (
	"github.com/muesli/beehive/bees"
)

// NotificationBeeFactory is a factory for NotificationBees.
type NotificationBeeFactory struct {
	bees.BeeFactory
}

// New returns a new Bee instance configured with the supplied options.
func (factory *NotificationBeeFactory) New(name, description string, options bees.BeeOptions) bees.BeeInterface {
	bee := NotificationBee{
		Bee: bees.NewBee(name, factory.ID(), description, options),
	}
	bee.ReloadOptions(options)

	return &bee
}

// ID returns the ID of this Bee.
func (factory *NotificationBeeFactory) ID() string {
	return "notificationbee"
}

// Name returns the name of this Bee.
func (factory *NotificationBeeFactory) Name() string {
	return "Desktop Notifications"
}

// Description returns the description of this Bee.
func (factory *NotificationBeeFactory) Description() string {
	return "Shows desktop notifications"
}

// Image returns the filename of an image for this Bee.
func (factory *NotificationBeeFactory) Image() string {
	return factory.ID() + ".png"
}

// LogoColor returns the preferred logo background color (used by the admin interface).
func (factory *NotificationBeeFactory) LogoColor() string {
	return "#00aeef"
}

// Actions describes the available actions provided by this Bee.
func (factory *NotificationBeeFactory) Actions() []bees.ActionDescriptor {
	actions := []bees.ActionDescriptor{
		{
			Namespace:   factory.Name(),
			Name:        "notify",
			Description: "Shows the given text as notification message",
			Options: []bees.PlaceholderDescriptor{
				{
					Name:        "text",
					Description: "The content of the notification",
					Type:        "string",
					Mandatory:   true,
				},
				{
					Name:        "urgency",
					Description: "The urgencylevel to display the notification with ('low', 'normal' or 'critical')",
					Type:        "string",
				},
			},
		},
	}
	return actions
}

func init() {
	f := NotificationBeeFactory{}
	bees.RegisterFactory(&f)
}

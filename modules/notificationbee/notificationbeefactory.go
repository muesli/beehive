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
	"github.com/muesli/beehive/modules"
)

type NotificationBeeFactory struct {
	modules.ModuleFactory
}

func (factory *NotificationBeeFactory) New(name, description string, options modules.BeeOptions) modules.ModuleInterface {
	bee := NotificationBee{}
	bee.Module = modules.Module{name, factory.Name(), description}
	return &bee
}

func (factory *NotificationBeeFactory) Name() string {
	return "notificationbee"
}

func (factory *NotificationBeeFactory) Description() string {
	return "A bee that shows desktop-notifications"
}

func (factory *NotificationBeeFactory) Actions() []modules.ActionDescriptor {
	actions := []modules.ActionDescriptor{
		modules.ActionDescriptor{
			Namespace:   factory.Name(),
			Name:        "notify",
			Description: "Shows the given text as notification message",
			Options: []modules.PlaceholderDescriptor{
				modules.PlaceholderDescriptor{
					Name:        "text",
					Description: "The content of the notification",
					Type:        "string",
				},
				modules.PlaceholderDescriptor{
					Name:        "urgency",
					Description: "The urgencylevel to display the notification with",
					Type:        "string",
				},
			},
		},
	}
	return actions
}

func init() {
	f := NotificationBeeFactory{}
	modules.RegisterFactory(&f)
}

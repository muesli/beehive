/*
 *    Copyright (C) 2017 Timm Schäuble
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
 *		Timm Schäuble <tymmm1+gh@gmail.com>
 */

package simplepushbee

import (
	"github.com/muesli/beehive/bees"
)

// SimplepushBeeFactory is a factory for SimplepushBees.
type SimplepushBeeFactory struct {
	bees.BeeFactory
}

// New returns a new Bee instance configured with the supplied options.
func (factory *SimplepushBeeFactory) New(name, description string, options bees.BeeOptions) bees.BeeInterface {
	bee := SimplepushBee{
		Bee: bees.NewBee(name, factory.ID(), description, options),
	}
	bee.ReloadOptions(options)

	return &bee
}

// ID returns the ID of this Bee.
func (factory *SimplepushBeeFactory) ID() string {
	return "simplepushbee"
}

// Name returns the name of this Bee.
func (factory *SimplepushBeeFactory) Name() string {
	return "Simplepush"
}

// Description returns the description of this Bee.
func (factory *SimplepushBeeFactory) Description() string {
	return "Lets you send (encrypted) push notifications to Android"
}

// Image returns the filename of an image for this Bee.
func (factory *SimplepushBeeFactory) Image() string {
	return factory.ID() + ".png"
}

// LogoColor returns the preferred logo background color (used by the admin interface).
func (factory *SimplepushBeeFactory) LogoColor() string {
	return "#2c3e50"
}

// Options returns the options available to configure this Bee.
func (factory *SimplepushBeeFactory) Options() []bees.BeeOptionDescriptor {
	opts := []bees.BeeOptionDescriptor{
		{
			Name:        "key",
			Description: "Simplepush key which you get after installing the app",
			Type:        "string",
			Mandatory:   true,
		},
		{
			Name:        "password",
			Description: "Password for end-to-end encryption (optional)",
			Type:        "string",
			Mandatory:   false,
		},
		{
			Name:        "salt",
			Description: "Salt for end-to-end encryption (optional)",
			Type:        "url",
			Mandatory:   false,
		},
	}
	return opts
}

// Events describes the available events provided by this Bee.
func (factory *SimplepushBeeFactory) Events() []bees.EventDescriptor {
	events := []bees.EventDescriptor{}
	return events
}

// Actions describes the available actions provided by this Bee.
func (factory *SimplepushBeeFactory) Actions() []bees.ActionDescriptor {
	actions := []bees.ActionDescriptor{
		{
			Namespace:   factory.Name(),
			Name:        "send",
			Description: "Sends a push notification",
			Options: []bees.PlaceholderDescriptor{
				{
					Name:        "title",
					Description: "Title of push notification (optional)",
					Type:        "string",
					Mandatory:   false,
				},
				{
					Name:        "message",
					Description: "Content of push notification",
					Type:        "string",
					Mandatory:   true,
				},
				{
					Name:        "event",
					Description: "Event id for customizing vibration and ringtone (optional)",
					Type:        "string",
					Mandatory:   false,
				},
			},
		},
	}
	return actions
}

func init() {
	f := SimplepushBeeFactory{}
	bees.RegisterFactory(&f)
}

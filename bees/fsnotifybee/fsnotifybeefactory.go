/*
 *    Copyright (C) 2017 Sergio Rubio
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
 *      Sergio Rubio <sergio@rubio.im>
 */

package fsnotifybee

import (
	"github.com/muesli/beehive/bees"
)

// FSNotifyBeeFactory is a factory for FSNotifyBees.
type FSNotifyBeeFactory struct {
	bees.BeeFactory
}

// New returns a new Bee instance configured with the supplied options.
func (factory *FSNotifyBeeFactory) New(name, description string, options bees.BeeOptions) bees.BeeInterface {
	bee := FSNotifyBee{
		Bee: bees.NewBee(name, factory.ID(), description, options),
	}
	bee.ReloadOptions(options)

	return &bee
}

// ID returns the ID of this Bee.
func (factory *FSNotifyBeeFactory) ID() string {
	return "fsnotifybee"
}

// Name returns the name of this Bee.
func (factory *FSNotifyBeeFactory) Name() string {
	return "FSNotify"
}

// Description returns the description of this Bee.
func (factory *FSNotifyBeeFactory) Description() string {
	return "Monitor filesystem paths"
}

// Image returns the filename of an image for this Bee.
func (factory *FSNotifyBeeFactory) Image() string {
	return factory.ID() + ".png"
}

// LogoColor returns the preferred logo background color (used by the admin interface).
func (factory *FSNotifyBeeFactory) LogoColor() string {
	return "#4b4b4b"
}

// Options returns the options available to configure this Bee.
func (factory *FSNotifyBeeFactory) Options() []bees.BeeOptionDescriptor {
	opts := []bees.BeeOptionDescriptor{
		{
			Name:        "path",
			Description: "Filesystem path to monitor",
			Type:        "string",
			Mandatory:   true,
		},
	}
	return opts
}

// Events describes the available events provided by this Bee.
func (factory *FSNotifyBeeFactory) Events() []bees.EventDescriptor {
	events := []bees.EventDescriptor{
		{
			Namespace:   factory.Name(),
			Name:        "fsevent",
			Description: "Filesystem event",
			Options: []bees.PlaceholderDescriptor{
				{
					Name:        "type", // CREATE, CHMOD, RENAME, REMOVE
					Description: "The event type received",
					Type:        "string",
				},
				{
					Name:        "path",
					Description: "Canonical path to the file or directory",
					Type:        "string",
				},
			},
		},
	}
	return events
}

func init() {
	f := FSNotifyBeeFactory{}
	bees.RegisterFactory(&f)
}

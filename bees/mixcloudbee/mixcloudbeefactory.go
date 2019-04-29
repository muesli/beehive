/*
 *    Copyright (C) 2018 Stefan Derkits
 *                  2018 Christian Muehlhaeuser
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
 *      Stefan Derkits <stefan@derkits.at>
 *      Christian Muehlhaeuser <muesli@gmail.com>
 */

package mixcloudbee

import "github.com/muesli/beehive/bees"

// mixcloudBeeFactory is a factory for mixcloudBees.
type MixcloudBeeFactory struct {
	bees.BeeFactory
}

// New returns a new Bee instance configured with the supplied options.
func (factory *MixcloudBeeFactory) New(name, description string, options bees.BeeOptions) bees.BeeInterface {
	bee := MixcloudBee{
		Bee: bees.NewBee(name, factory.ID(), description, options),
	}
	bee.ReloadOptions(options)

	return &bee
}

// ID returns the ID of this Bee.
func (factory *MixcloudBeeFactory) ID() string {
	return "mixcloudbee"
}

// Name returns the name of this Bee.
func (factory *MixcloudBeeFactory) Name() string {
	return "Mixcloud"
}

// Description returns the description of this Bee.
func (factory *MixcloudBeeFactory) Description() string {
	return "Interact with Mixcloud"
}

// Image returns the filename of an image for this Bee.
func (factory *MixcloudBeeFactory) Image() string {
	return factory.ID() + ".png"
}

// LogoColor returns the preferred logo background color (used by the admin interface).
func (factory *MixcloudBeeFactory) LogoColor() string {
	return "#52aad8"
}

// Options returns the options available to configure this Bee.
func (factory *MixcloudBeeFactory) Options() []bees.BeeOptionDescriptor {
	opts := []bees.BeeOptionDescriptor{
		{
			Name:        "feed",
			Description: "Feed to follow",
			Type:        "string",
			Mandatory:   true,
		},
	}
	return opts
}

// Events describes the available events provided by this Bee.
func (factory *MixcloudBeeFactory) Events() []bees.EventDescriptor {
	events := []bees.EventDescriptor{
		{
			Namespace:   factory.Name(),
			Name:        "new_cloudcast",
			Description: "A new cloudcast is available",
			Options: []bees.PlaceholderDescriptor{
				{
					Name:        "name",
					Description: "Name of the Cloudcast",
					Type:        "string",
				},
				{
					Name:        "url",
					Description: "Cloudcast URL",
					Type:        "string",
				},
				{
					Name:        "slug",
					Description: "Cloudcast Slug",
					Type:        "string",
				},
			},
		},
	}
	return events
}

func init() {
	f := MixcloudBeeFactory{}
	bees.RegisterFactory(&f)
}

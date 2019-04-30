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

package youtubebee

import (
	"github.com/muesli/beehive/bees"
)

// WebBeeFactory is a factory for WebBees.
type YoutubeBeeFactory struct {
	bees.BeeFactory
}

// New returns a new Bee instance configured with the supplied options.
func (factory *YoutubeBeeFactory) New(name, description string, options bees.BeeOptions) bees.BeeInterface {
	bee := YoutubeBee{
		Bee: bees.NewBee(name, factory.ID(), description, options),
	}
	bee.ReloadOptions(options)

	return &bee
}

// ID returns the ID of this Bee.
func (factory *YoutubeBeeFactory) ID() string {
	return "youtubebee"
}

// Name returns the name of this Bee.
func (factory *YoutubeBeeFactory) Name() string {
	return "Youtube Bee"
}

// Description returns the description of this Bee.
func (factory *YoutubeBeeFactory) Description() string {
	return "HTTP Server that listens to a Youtube channel"
}

// Image returns the filename of an image for this Bee.
func (factory *YoutubeBeeFactory) Image() string {
	return factory.ID() + ".png"
}

// LogoColor returns the preferred logo background color (used by the admin interface).
func (factory *YoutubeBeeFactory) LogoColor() string {
	return "#ff0000"
}

// Options returns the options available to configure this Bee.
func (factory *YoutubeBeeFactory) Options() []bees.BeeOptionDescriptor {
	opts := []bees.BeeOptionDescriptor{
		{
			Name:        "address",
			Description: "Which addr to listen on, eg: 0.0.0.0:12345",
			Type:        "address",
			Mandatory:   true,
		},
		{
			Name:        "channel",
			Description: "What is the link of the channel you want to receive push notifications for?",
			Type:        "url",
			Mandatory:   true,
		},
	}
	return opts
}

// Events describes the available events provided by this Bee.
func (factory *YoutubeBeeFactory) Events() []bees.EventDescriptor {
	events := []bees.EventDescriptor{
		{
			Namespace:   factory.Name(),
			Name:        "push",
			Description: "A push notification was sent by the Youtube channel",
			Options: []bees.PlaceholderDescriptor{
				{
					Name:        "channelUrl",
					Description: "The url of the channel push notification was sent from",
					Type:        "url",
				},
				{
					Name:        "vidUrl",
					Description: "The url of the video relevant to the push notification",
					Type:        "url",
				},
			},
		},
	}
	return events
}

func init() {
	f := YoutubeBeeFactory{}
	bees.RegisterFactory(&f)
}

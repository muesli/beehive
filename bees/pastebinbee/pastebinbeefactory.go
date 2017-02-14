/*
 *    Copyright (C) 2017 Sebastian Ławniczak
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
 *      Sebastian Ławniczak <seb@seblw.me>
 */

package pastebinbee

import (
	"github.com/muesli/beehive/bees"
)

// PastebinBeeFactory is a factory for PastebinBees.
type PastebinBeeFactory struct {
	bees.BeeFactory
}

// New returns a new Bee instance configured with the supplied options.
func (factory *PastebinBeeFactory) New(name, description string, options bees.BeeOptions) bees.BeeInterface {
	bee := PastebinBee{
		Bee: bees.NewBee(name, factory.ID(), description, options),
	}
	bee.ReloadOptions(options)

	return &bee
}

// ID returns the ID of this Bee.
func (factory *PastebinBeeFactory) ID() string {
	return "pastebinbee"
}

// Name returns the name of this Bee.
func (factory *PastebinBeeFactory) Name() string {
	return "Pastebin"
}

// Description returns the description of this Bee.
func (factory *PastebinBeeFactory) Description() string {
	return "Posts a snippet on Pastebin"
}

// Image returns the filename of an image for this Bee.
func (factory *PastebinBeeFactory) Image() string {
	return factory.ID() + ".png"
}

// LogoColor returns the preferred logo background color (used by the admin interface).
func (factory *PastebinBeeFactory) LogoColor() string {
	return "#00796B"
}

// Options returns the options available to configure this Bee.
func (factory *PastebinBeeFactory) Options() []bees.BeeOptionDescriptor {
	opts := []bees.BeeOptionDescriptor{
		{
			Name:        "api_dev_key",
			Description: "Developer key for Pastebin API",
			Type:        "string",
		},
	}
	return opts
}

// Events describes the available events provided by this Bee.
func (factory *PastebinBeeFactory) Events() []bees.EventDescriptor {
	events := []bees.EventDescriptor{}
	return events
}

// Actions describes the available actions provided by this Bee.
func (factory *PastebinBeeFactory) Actions() []bees.ActionDescriptor {
	actions := []bees.ActionDescriptor{
		{
			Namespace:   factory.Name(),
			Name:        "post",
			Description: "Posts a snippet on Pastebin",
			Options: []bees.PlaceholderDescriptor{
				{
					Name:        "title",
					Description: "Title",
					Type:        "string",
					Mandatory:   false,
				},
				{
					Name:        "content",
					Description: "Content",
					Type:        "string",
					Mandatory:   true,
				},
				{
					Name:        "expire",
					Description: "Expire date (available: N, 10M, 1H, 1D, 1W, 2W, 1M)",
					Type:        "string",
					Mandatory:   true,
				},
				{
					Name:        "exposure",
					Description: "Exposure (0 - public, 1 - unlisted, 2 - private)",
					Type:        "string",
					Mandatory:   true,
				},
			},
		},
	}
	return actions
}

func init() {
	f := PastebinBeeFactory{}
	bees.RegisterFactory(&f)
}

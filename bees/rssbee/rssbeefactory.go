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

package rssbee

import (
	"github.com/muesli/beehive/bees"
)

// RSSBeeFactory is a factory for RSSBees.
type RSSBeeFactory struct {
	bees.BeeFactory
}

// New returns a new Bee instance configured with the supplied options.
func (factory *RSSBeeFactory) New(name, description string, options bees.BeeOptions) bees.BeeInterface {
	bee := RSSBee{
		Bee: bees.NewBee(name, factory.ID(), description, options),
	}
	bee.ReloadOptions(options)

	return &bee
}

// ID returns the ID of this Bee.
func (factory *RSSBeeFactory) ID() string {
	return "rssbee"
}

// Name returns the name of this Bee.
func (factory *RSSBeeFactory) Name() string {
	return "RSS Feeds"
}

// Description returns the description of this Bee.
func (factory *RSSBeeFactory) Description() string {
	return "Reacts to RSS-feed updates"
}

// Image returns the filename of an image for this Bee.
func (factory *RSSBeeFactory) Image() string {
	return factory.ID() + ".png"
}

// LogoColor returns the preferred logo background color (used by the admin interface).
func (factory *RSSBeeFactory) LogoColor() string {
	return "#ec7505"
}

// Options returns the options available to configure this Bee.
func (factory *RSSBeeFactory) Options() []bees.BeeOptionDescriptor {
	opts := []bees.BeeOptionDescriptor{
		{
			Name:        "url",
			Description: "URL of the RSS-feed you want to monitor",
			Type:        "url",
			Mandatory:   true,
		},
		{
			Name:        "skip_first",
			Description: "Whether to skip already existing entries",
			Type:        "bool",
			Mandatory:   false,
			Default:     false,
		},
		{
			Name:        "skip_first_allow_newest",
			Description: "Whether to skip already existing entries, but allow the newest item to get through once",
			Type:        "bool",
			Mandatory:   false,
			Default:     false,
		},
	}
	return opts
}

// Events describes the available events provided by this Bee.
func (factory *RSSBeeFactory) Events() []bees.EventDescriptor {
	events := []bees.EventDescriptor{
		{
			Namespace:   factory.Name(),
			Name:        "new_item",
			Description: "A new item has been received through the Feed",
			Options: []bees.PlaceholderDescriptor{
				{
					Name:        "title",
					Description: "Title of the Item",
					Type:        "string",
				},
				{
					Name:        "links",
					Description: "Links referenced by the Item",
					Type:        "[]string",
				},
				{
					Name:        "description",
					Description: "Description of the Item",
					Type:        "string",
				},
				{
					Name:        "author",
					Description: "The person who wrote the Item",
					Type:        "string",
				},
				{
					Name:        "categories",
					Description: "Categories that the Item belongs to",
					Type:        "[]string",
				},
				{
					Name:        "comments",
					Description: "Comments of the Item",
					Type:        "string",
				},
				{
					Name:        "enclosures",
					Description: "Enclosures related to Item",
					Type:        "[]string",
				},
				{
					Name:        "guid",
					Description: "Global unique ID attached to the Item",
					Type:        "string",
				},
				{
					Name:        "pubdate",
					Description: "Date the Item was published on",
					Type:        "string",
				},
				{
					Name:        "source",
					Description: "Source of the Item",
					Type:        "string",
				},
			},
		},
	}
	return events
}

func init() {
	f := RSSBeeFactory{}
	bees.RegisterFactory(&f)
}

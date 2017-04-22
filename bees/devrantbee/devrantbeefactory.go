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
 *      Nicolas Martin <penguwingithub@gmail.com>
 *      Christian Muehlhaeuser <muesli@gmail.com>
 */

package devrantbee

import (
	"github.com/muesli/beehive/bees"
)

// DevrantBeeFactory is a factory for DevrantBees.
type DevrantBeeFactory struct {
	bees.BeeFactory
}

// New returns a new Bee instance configured with the supplied options.
func (factory *DevrantBeeFactory) New(name, description string, options bees.BeeOptions) bees.BeeInterface {
	bee := DevrantBee{
		Bee: bees.NewBee(name, factory.ID(), description, options),
	}
	bee.ReloadOptions(options)

	return &bee
}

// ID returns the ID of this Bee.
func (factory *DevrantBeeFactory) ID() string {
	return "devrantbee"
}

// Name returns the name of this Bee.
func (factory *DevrantBeeFactory) Name() string {
	return "Devrant"
}

// Description returns the description of this Bee.
func (factory *DevrantBeeFactory) Description() string {
	return "Retrieves rants from Devrant"
}

// Image returns the filename of an image for this Bee.
func (factory *DevrantBeeFactory) Image() string {
	return factory.ID() + ".png"
}

// LogoColor returns the preferred logo background color (used by the admin interface).
func (factory *DevrantBeeFactory) LogoColor() string {
	return "#54556E"
}

// Events describes the available events provided by this Bee.
func (factory *DevrantBeeFactory) Events() []bees.EventDescriptor {
	events := []bees.EventDescriptor{
		{
			Namespace:   factory.Name(),
			Name:        "rant",
			Description: "is triggered after rants were fetched",
			Options: []bees.PlaceholderDescriptor{
				{
					Name:        "Id",
					Description: "ID of a rant",
					Type:        "int",
				},
				{
					Name:        "Text",
					Description: "Text of a rant",
					Type:        "string",
				},
				{
					Name:        "Upvotes",
					Description: "Sum of upvotes",
					Type:        "int",
				},
				{
					Name:        "Downvotes",
					Description: "Sum of downvotes",
					Type:        "int",
				},
				{
					Name:        "Score",
					Description: "Current score of a rant",
					Type:        "int",
				},
				{
					Name:        "CreatedTime",
					Description: "Creation time of a rant",
					Type:        "int",
				},
				{
					Name:        "NumComments",
					Description: "Number of comments",
					Type:        "int",
				},
				{
					Name:        "Tags",
					Description: "Tagged tags",
					Type:        "[]string",
				},
				{
					Name:        "UserId",
					Description: "ID of the user who posted the rant",
					Type:        "int",
				},
				{
					Name:        "UserUsername",
					Description: "Username of the ranter",
					Type:        "string",
				},
				{
					Name:        "UserScore",
					Description: "Total user score",
					Type:        "int",
				},
			},
		},
	}
	return events
}

// Actions describes the available actions provided by this Bee.
func (factory *DevrantBeeFactory) Actions() []bees.ActionDescriptor {
	actions := []bees.ActionDescriptor{
		{
			Namespace:   factory.Name(),
			Name:        "surprise",
			Description: "Fetches a surprise rant",
			Options:     []bees.PlaceholderDescriptor{},
		},
		{
			Namespace:   factory.Name(),
			Name:        "weekly",
			Description: "Fetches the top weekly rants",
			Options:     []bees.PlaceholderDescriptor{},
		},
		{
			Namespace:   factory.Name(),
			Name:        "rant",
			Description: "Fetches new rants",
			Options: []bees.PlaceholderDescriptor{
				{
					Name:        "limit",
					Description: "How many rants to retrieve",
					Type:        "int",
					Mandatory:   true,
				},
			},
		},
	}
	return actions
}

func init() {
	f := DevrantBeeFactory{}
	bees.RegisterFactory(&f)
}

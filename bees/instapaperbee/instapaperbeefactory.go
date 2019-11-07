/*
 *    Copyright (C) 2019 Adam Petrovic
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
 *      Adam Petrovic <adam@petrovic.com.au>
 */

package instapaperbee

import "github.com/muesli/beehive/bees"

type InstapaperBeeFactory struct {
	bees.BeeFactory
}

func (factory *InstapaperBeeFactory) New(name, description string, options bees.BeeOptions) bees.BeeInterface {
	bee := InstapaperBee{
		Bee: bees.NewBee(name, factory.ID(), description, options),
	}
	bee.ReloadOptions(options)
	return &bee
}

func (factory *InstapaperBeeFactory) ID() string {
	return "instapaperbee"
}

func (factory *InstapaperBeeFactory) Name() string {
	return "Instapaper"
}

func (factory *InstapaperBeeFactory) Description() string {
	return "Add to Instapaper"
}

func (factory *InstapaperBeeFactory) Image() string {
	return factory.ID() + ".png"
}

func (factory *InstapaperBeeFactory) LogoColor() string {
	return "#808080"
}

func (factory *InstapaperBeeFactory) Options() []bees.BeeOptionDescriptor {
	opts := []bees.BeeOptionDescriptor{
		{
			Name:        "username",
			Description: "Instapaper Username / Email",
			Type:        "string",
			Mandatory:   true,
		},
		{
			Name:        "password",
			Description: "Instapaper Password",
			Type:        "string",
			Mandatory:   true,
		},
	}
	return opts
}

func (factory *InstapaperBeeFactory) Events() []bees.EventDescriptor {
	return []bees.EventDescriptor{}
}

func (factory *InstapaperBeeFactory) Actions() []bees.ActionDescriptor {
	actions := []bees.ActionDescriptor{
		{
			Namespace:   factory.Name(),
			Name:        "save",
			Description: "Saves a URL to Instapaper",
			Options: []bees.PlaceholderDescriptor{
				{
					Name:        "title",
					Description: "Article title",
					Type:        "string",
					Mandatory:   false,
				},
				{
					Name:        "url",
					Description: "Article URL",
					Type:        "string",
					Mandatory:   true,
				},
			},
		},
	}
	return actions
}

func init() {
	f := InstapaperBeeFactory{}
	bees.RegisterFactory(&f)
}

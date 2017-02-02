/*
 *    Copyright (C) 2014      Stefan 'glaxx' Luecke
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
 *		Stefan Luecke <glaxx@glaxx.net>
 *      Christian Muehlhaeuser <muesli@gmail.com>
 */

package cronbee

import (
	"github.com/muesli/beehive/bees"
)

// CronBeeFactory is a factory for CronBees.
type CronBeeFactory struct {
	bees.BeeFactory
}

// New returns a new Bee instance configured with the supplied options.
func (factory *CronBeeFactory) New(name, description string, options bees.BeeOptions) bees.BeeInterface {
	bee := CronBee{
		Bee: bees.NewBee(name, factory.Name(), description, options),
	}
	bee.ReloadOptions(options)

	return &bee
}

// Name returns the name of this Bee.
func (factory *CronBeeFactory) Name() string {
	return "cronbee"
}

// Description returns the description of this Bee.
func (factory *CronBeeFactory) Description() string {
	return "A bee that triggers events in given intervals"
}

// Image returns the filename of an image for this Bee.
func (factory *CronBeeFactory) Image() string {
	return factory.Name() + ".png"
}

// LogoColor returns the preferred logo background color (used by the admin interface).
func (factory *CronBeeFactory) LogoColor() string {
	return "#aaaaaa"
}

// Options returns the options available to configure this Bee.
func (factory *CronBeeFactory) Options() []bees.BeeOptionDescriptor {
	opts := []bees.BeeOptionDescriptor{
		{
			Name:        "second",
			Description: "00-59 for a specific time; * for ignore",
			Type:        "string",
		},
		{
			Name:        "minute",
			Description: "00-59 for a specific time; * for ignore",
			Type:        "string",
		},
		{
			Name:        "hour",
			Description: "00-23 for a specific time; * for ignore",
			Type:        "string",
		},
		{
			Name:        "day_of_week",
			Description: "0-6 0 = Sunday 6 = Saturday; * for ignore",
			Type:        "string",
		},
		{
			Name:        "day_of_month",
			Description: "01-31 for a specific time; * for ignore)",
			Type:        "string",
		},
		{
			Name:        "month",
			Description: "01 - 12 for a specific time; * for ignore)",
			Type:        "string",
		},
	}
	return opts
}

// Events describes the available events provided by this Bee.
func (factory *CronBeeFactory) Events() []bees.EventDescriptor {
	events := []bees.EventDescriptor{
		{
			Namespace:   factory.Name(),
			Name:        "time",
			Description: "The time has come ...",
			Options: []bees.PlaceholderDescriptor{
				{
					Name:        "timestamp",
					Description: "Timestamp of the next event",
					Type:        "string",
				},
			},
		},
	}
	return events
}

func init() {
	f := CronBeeFactory{}
	bees.RegisterFactory(&f)
}

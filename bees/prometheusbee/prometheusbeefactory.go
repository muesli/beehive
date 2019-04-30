/*
 *    Copyright (C) 2019 CalmBit
 *                  2014-2019 Christian Muehlhaeuser
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
 *      CalmBit <calmbit@posteo.net>
 *      Christian Muehlhaeuser <muesli@gmail.com>
 */

package prometheusbee

import (
	"github.com/muesli/beehive/bees"
)

// PrometheusBeeFactory is a factory for PrometheusBees.
type PrometheusBeeFactory struct {
	bees.BeeFactory
}

// New returns a new Bee instance configured with the supplied options.
func (factory *PrometheusBeeFactory) New(name, description string, options bees.BeeOptions) bees.BeeInterface {
	bee := PrometheusBee{
		Bee: bees.NewBee(name, factory.ID(), description, options),
	}
	bee.ReloadOptions(options)

	return &bee
}

// ID returns the ID of this Bee.
func (factory *PrometheusBeeFactory) ID() string {
	return "prometheusbee"
}

// Name returns the name of this Bee.
func (factory *PrometheusBeeFactory) Name() string {
	return "Prometheus"
}

// Description returns the description of this Bee.
func (factory *PrometheusBeeFactory) Description() string {
	return "Allows for the export of data to Prometheus"
}

// Image returns the filename of an image for this Bee.
func (factory *PrometheusBeeFactory) Image() string {
	return factory.ID() + ".png"
}

// LogoColor returns the preferred logo background color (used by the admin interface).
func (factory *PrometheusBeeFactory) LogoColor() string {
	return "#e59030"
}

// Options returns the options available to configure this Bee.
func (factory *PrometheusBeeFactory) Options() []bees.BeeOptionDescriptor {
	opts := []bees.BeeOptionDescriptor{
		{
			Name:        "address",
			Description: "Which addr to listen on, eg: 0.0.0.0:2112",
			Type:        "address",
			Mandatory:   true,
		},
		{
			Name:        "counter_vec_name",
			Description: "The name of the counter vector",
			Type:        "string",
			Mandatory:   true,
			Default:     "counter_vec",
		},
		{
			Name:        "gauge_vec_name",
			Description: "The name of the gauge vector",
			Type:        "string",
			Mandatory:   true,
			Default:     "gauge_vec",
		},
		{
			Name:        "histogram_vec_name",
			Description: "The name of the histogram vector",
			Type:        "string",
			Mandatory:   true,
			Default:     "histogram_vec",
		},
		{
			Name:        "summary_vec_name",
			Description: "The name of the summary vector",
			Type:        "string",
			Mandatory:   true,
			Default:     "summary_vec",
		},
	}
	return opts
}

// Actions describes the available actions provided by this Bee.
func (factory *PrometheusBeeFactory) Actions() []bees.ActionDescriptor {
	actions := []bees.ActionDescriptor{
		{
			Namespace:   factory.Name(),
			Name:        "counter_inc",
			Description: "Increments the value of a counter by 1",
			Options: []bees.PlaceholderDescriptor{
				{
					Name:        "label",
					Description: "Label of the counter to increment",
					Type:        "string",
					Mandatory:   true,
				},
			},
		},
		{
			Namespace:   factory.Name(),
			Name:        "counter_add",
			Description: "Adds an arbitrary, positive value to a counter",
			Options: []bees.PlaceholderDescriptor{
				{
					Name:        "label",
					Description: "Label of the counter to add to",
					Type:        "string",
					Mandatory:   true,
				},
				{
					Name:        "value",
					Description: "Value to add to the counter",
					Type:        "float64",
					Mandatory:   true,
				},
			},
		},
		{
			Namespace:   factory.Name(),
			Name:        "gauge_set",
			Description: "Sets a gauge's value to an arbitrary value",
			Options: []bees.PlaceholderDescriptor{
				{
					Name:        "label",
					Description: "Label of the gauge to set",
					Type:        "string",
					Mandatory:   true,
				},
				{
					Name:        "value",
					Description: "New value for the gauge",
					Type:        "float64",
					Mandatory:   true,
				},
			},
		},
		{
			Namespace:   factory.Name(),
			Name:        "gauge_inc",
			Description: "Increments the value of a gauge by 1",
			Options: []bees.PlaceholderDescriptor{
				{
					Name:        "label",
					Description: "Label of the gauge to increment",
					Type:        "string",
					Mandatory:   true,
				},
			},
		},
		{
			Namespace:   factory.Name(),
			Name:        "gauge_dec",
			Description: "Decrements the value of a gauge by 1",
			Options: []bees.PlaceholderDescriptor{
				{
					Name:        "label",
					Description: "Label of the gauge to decrement",
					Type:        "string",
					Mandatory:   true,
				},
			},
		},
		{
			Namespace:   factory.Name(),
			Name:        "gauge_add",
			Description: "Adds an arbitrary value to a gauge",
			Options: []bees.PlaceholderDescriptor{
				{
					Name:        "label",
					Description: "Label of the gauge to add to",
					Type:        "string",
					Mandatory:   true,
				},
				{
					Name:        "value",
					Description: "Value to add to the counter",
					Type:        "float64",
					Mandatory:   true,
				},
			},
		},
		{
			Namespace:   factory.Name(),
			Name:        "gauge_sub",
			Description: "Subtracts an arbitrary value from a gauge",
			Options: []bees.PlaceholderDescriptor{
				{
					Name:        "label",
					Description: "Label of the gauge to subtract from",
					Type:        "string",
					Mandatory:   true,
				},
				{
					Name:        "value",
					Description: "Value to subtract from the gauge",
					Type:        "float64",
					Mandatory:   true,
				},
			},
		},
		{
			Namespace:   factory.Name(),
			Name:        "gauge_set_to_current_time",
			Description: "Sets a gauge's vaue to the current time as a unix timestamp",
			Options: []bees.PlaceholderDescriptor{
				{
					Name:        "label",
					Description: "Label of the gauge to set",
					Type:        "string",
					Mandatory:   true,
				},
			},
		},
		{
			Namespace:   factory.Name(),
			Name:        "histogram_observe",
			Description: "Records an observation in a histogram",
			Options: []bees.PlaceholderDescriptor{
				{
					Name:        "label",
					Description: "Label of the histogram to add an observation to",
					Type:        "string",
					Mandatory:   true,
				},
				{
					Name:        "value",
					Description: "Value to observe",
					Type:        "float64",
					Mandatory:   true,
				},
			},
		},
		{
			Namespace:   factory.Name(),
			Name:        "summary_observe",
			Description: "Records an observation in a summary",
			Options: []bees.PlaceholderDescriptor{
				{
					Name:        "label",
					Description: "Label of the summary to add an observation to",
					Type:        "string",
					Mandatory:   true,
				},
				{
					Name:        "value",
					Description: "Value to observe",
					Type:        "float64",
					Mandatory:   true,
				},
			},
		},
	}
	return actions
}

func init() {
	f := PrometheusBeeFactory{}
	bees.RegisterFactory(&f)
}

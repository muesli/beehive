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
 *      Nicolas Martin <penguwingit@gmail.com>
 */

package openweathermapbee

import (
	"github.com/muesli/beehive/bees"
)

// OpenweathermapBeeFactory is a factory for openweathermapbees
type OpenweathermapBeeFactory struct {
	bees.BeeFactory
}

// New returns a new Bee instance configured with the supplied options.
func (factory *OpenweathermapBeeFactory) New(name, description string, options bees.BeeOptions) bees.BeeInterface {
	bee := OpenweathermapBee{
		Bee: bees.NewBee(name, factory.ID(), description, options),
	}
	bee.ReloadOptions(options)

	return &bee
}

// ID returns the ID of this Bee.
func (factory *OpenweathermapBeeFactory) ID() string {
	return "openweathermapbee"
}

// Name returns the name of this Bee.
func (factory *OpenweathermapBeeFactory) Name() string {
	return "OpenWeatherMap"
}

// Description returns the description of this Bee.
func (factory *OpenweathermapBeeFactory) Description() string {
	return "Retrieves weather information from openweathermap.org"
}

// Image returns the filename of an image for this Bee.
func (factory *OpenweathermapBeeFactory) Image() string {
	return factory.ID() + ".png"
}

// LogoColor returns the preferred logo background color (used by the admin interface).
func (factory *OpenweathermapBeeFactory) LogoColor() string {
	return "#7D1640"
}

// Options returns the options available to configure this Bee.
func (factory *OpenweathermapBeeFactory) Options() []bees.BeeOptionDescriptor {
	opts := []bees.BeeOptionDescriptor{
		{
			Name:        "unit",
			Description: "Your preferred unit",
			Type:        "string",
			Default:     "c", // celcius -> The right one
			Mandatory:   true,
		},
		{
			Name:        "language",
			Description: "Your preferred language",
			Type:        "string",
			Default:     "en", // english
			Mandatory:   true,
		},
		{
			Name:        "key",
			Description: "Your OpenWeatherMap api key",
			Type:        "string",
			Mandatory:   true,
		},
	}
	return opts
}

// Events describes the available events provided by this Bee.
func (factory *OpenweathermapBeeFactory) Events() []bees.EventDescriptor {
	events := []bees.EventDescriptor{
		{
			Namespace:   factory.Name(),
			Name:        "current_weather",
			Description: "Current weather holds measurement informations",
			Options: []bees.PlaceholderDescriptor{
				{
					Name:        "geo_pos_longitude",
					Type:        "float64",
					Description: "Geo position longitude",
				},
				{
					Name:        "geo_pos_latitude",
					Type:        "float64",
					Description: "Geo position latitude",
				},
				{
					Name:        "sys_type",
					Type:        "int",
					Description: "sys type",
				},
				{
					Name:        "sys_id",
					Type:        "int",
					Description: "sys ID",
				},
				{
					Name:        "sys_message",
					Type:        "float64",
					Description: "sys message",
				},
				{
					Name:        "sys_country",
					Type:        "string",
					Description: "sys country code",
				},
				{
					Name:        "sys_sunrise",
					Type:        "int",
					Description: "sys sunrise",
				},
				{
					Name:        "sys_sunset",
					Type:        "int",
					Description: "sys sunset",
				},
				{
					Name:        "base",
					Type:        "string",
					Description: "current weather base",
				},
				{
					Name:        "main_temp",
					Type:        "float64",
					Description: "current main temperature",
				},
				{
					Name:        "main_temp_min",
					Type:        "float64",
					Description: "current main minimum temperature",
				},
				{
					Name:        "main_temp_max",
					Type:        "float64",
					Description: "current main maximum temperature",
				},
				{
					Name:        "main_pressure",
					Type:        "float64",
					Description: "current air pressure",
				},
				{
					Name:        "main_sealevel",
					Type:        "float64",
					Description: "current sealevel",
				},
				{
					Name:        "main_grndlevel",
					Type:        "float64",
					Description: "current groundlevel",
				},
				{
					Name:        "main_humidity",
					Type:        "float64",
					Description: "main humidity",
				},
				{
					Name:        "wind_speed",
					Type:        "float64",
					Description: "wind speed",
				},
				{
					Name:        "wind_deg",
					Type:        "float64",
					Description: "wind degree",
				},
				{
					Name:        "clouds_all",
					Type:        "int",
					Description: "clouds",
				},
				{
					Name:        "rain",
					Type:        "map[string]float64",
					Description: "rain",
				},
				{
					Name:        "snow",
					Type:        "map[string]float64",
					Description: "snow",
				},
				{
					Name:        "dt",
					Type:        "int",
					Description: "dt",
				},
				{
					Name:        "id",
					Type:        "int",
					Description: "id",
				},
				{
					Name:        "name",
					Type:        "string",
					Description: "name",
				},
				{
					Name:        "cod",
					Type:        "int",
					Description: "cod",
				},
				{
					Name:        "unit",
					Type:        "string",
					Description: "unit meassurement system",
				},
				{
					Name:        "lang",
					Type:        "string",
					Description: "language",
				},
				{
					Name:        "key",
					Type:        "string",
					Description: "key",
				},
			},
		},
		{
			Namespace:   factory.Name(),
			Name:        "main_weather",
			Description: "returns main measurement weather informations",
			Options: []bees.PlaceholderDescriptor{
				{
					Name:        "id",
					Type:        "int",
					Description: "weather id",
				},
				{
					Name:        "main",
					Type:        "string",
					Description: "main weather information",
				},
				{
					Name:        "description",
					Type:        "string",
					Description: "current weather main description",
				},
				{
					Name:        "icon",
					Type:        "string",
					Description: "weather icon",
				},
			},
		},
	}
	return events
}

// Actions describes the available actions provided by this Bee.
func (factory *OpenweathermapBeeFactory) Actions() []bees.ActionDescriptor {
	actions := []bees.ActionDescriptor{
		{
			Namespace:   factory.Name(),
			Name:        "get_current_weather",
			Description: "fetch current weather",
			Options: []bees.PlaceholderDescriptor{
				{
					Name:        "location",
					Description: "desired location",
					Type:        "string",
					Mandatory:   true,
				},
			},
		},
	}
	return actions
}

func init() {
	f := OpenweathermapBeeFactory{}
	bees.RegisterFactory(&f)
}

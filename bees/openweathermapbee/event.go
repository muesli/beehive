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

	owm "github.com/briandowns/openweathermap"
)

// TriggerCurrentWeatherEvent triggers all current weather events
func (mod *OpenweathermapBee) TriggerCurrentWeatherEvent() {
	ev := bees.Event{
		Bee:  mod.Name(),
		Name: "current_weather",
		Options: []bees.Placeholder{
			{
				Name:  "geo_pos_longitude",
				Type:  "float64",
				Value: mod.current.GeoPos.Longitude,
			},
			{
				Name:  "geo_pos_latitude",
				Type:  "float64",
				Value: mod.current.GeoPos.Latitude,
			},
			{
				Name:  "sys_type",
				Type:  "int",
				Value: mod.current.Sys.Type,
			},
			{
				Name:  "sys_id",
				Type:  "int",
				Value: mod.current.Sys.ID,
			},
			{
				Name:  "sys_message",
				Type:  "float64",
				Value: mod.current.Sys.Message,
			},
			{
				Name:  "sys_country",
				Type:  "string",
				Value: mod.current.Sys.Country,
			},
			{
				Name:  "sys_sunrise",
				Type:  "int",
				Value: mod.current.Sys.Sunrise,
			},
			{
				Name:  "sys_sunset",
				Type:  "int",
				Value: mod.current.Sys.Sunset,
			},
			{
				Name:  "base",
				Type:  "string",
				Value: mod.current.Base,
			},
			{
				Name:  "main_temp",
				Type:  "float64",
				Value: mod.current.Main.Temp,
			},
			{
				Name:  "main_temp_min",
				Type:  "float64",
				Value: mod.current.Main.TempMin,
			},
			{
				Name:  "main_temp_max",
				Type:  "float64",
				Value: mod.current.Main.TempMax,
			},
			{
				Name:  "main_pressure",
				Type:  "float64",
				Value: mod.current.Main.Pressure,
			},
			{
				Name:  "main_sealevel",
				Type:  "float64",
				Value: mod.current.Main.SeaLevel,
			},
			{
				Name:  "main_grndlevel",
				Type:  "float64",
				Value: mod.current.Main.GrndLevel,
			},
			{
				Name:  "main_humidity",
				Type:  "float64",
				Value: mod.current.Main.Humidity,
			},
			{
				Name:  "wind_speed",
				Type:  "float64",
				Value: mod.current.Wind.Speed,
			},
			{
				Name:  "wind_deg",
				Type:  "float64",
				Value: mod.current.Wind.Deg,
			},
			{
				Name:  "clouds_all",
				Type:  "int",
				Value: mod.current.Clouds.All,
			},
			{
				Name:  "rain",
				Type:  "map[string]float64",
				Value: mod.current.Rain,
			},
			{
				Name:  "snow",
				Type:  "map[string]float64",
				Value: mod.current.Snow,
			},
			{
				Name:  "dt",
				Type:  "int",
				Value: mod.current.Dt,
			},
			{
				Name:  "id",
				Type:  "int",
				Value: mod.current.ID,
			},
			{
				Name:  "name",
				Type:  "string",
				Value: mod.current.Name,
			},
			{
				Name:  "cod",
				Type:  "int",
				Value: mod.current.Cod,
			},
			{
				Name:  "unit",
				Type:  "string",
				Value: mod.current.Unit,
			},
			{
				Name:  "lang",
				Type:  "string",
				Value: mod.current.Lang,
			},
			{
				Name:  "key",
				Type:  "string",
				Value: mod.current.Key,
			},
		},
	}
	mod.evchan <- ev

	for _, v := range mod.current.Weather {
		mod.TriggerWeatherInformationEvent(&v)
	}
}

// WeatherInformationEvent triggers a weather event
func (mod *OpenweathermapBee) TriggerWeatherInformationEvent(v *owm.Weather) {
	weather := bees.Event{
		Bee:  mod.Name(),
		Name: "main_weather",
		Options: []bees.Placeholder{
			{
				Name:  "id",
				Type:  "int",
				Value: v.ID,
			},
			{
				Name:  "main",
				Type:  "string",
				Value: v.Main,
			},
			{
				Name:  "description",
				Type:  "string",
				Value: v.Description,
			},
			{
				Name:  "icon",
				Type:  "string",
				Value: v.Icon,
			},
		},
	}
	mod.evchan <- weather
}

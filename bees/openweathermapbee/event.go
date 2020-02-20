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
 *      Nicolas Martin <penguwin@penguwin.eu>
 */

package openweathermapbee

import (
	owm "github.com/briandowns/openweathermap"
)

// TriggerCurrentWeatherEvent triggers all current weather events
func (mod *OpenweathermapBee) TriggerCurrentWeatherEvent() {
	properties := map[string]interface{}{
		"geo_pos_longitude": mod.current.GeoPos.Longitude,
		"geo_pos_latitude":  mod.current.GeoPos.Latitude,
		"sys_type":          mod.current.Sys.Type,
		"sys_id":            mod.current.Sys.ID,
		"sys_message":       mod.current.Sys.Message,
		"sys_country":       mod.current.Sys.Country,
		"sys_sunrise":       mod.current.Sys.Sunrise,
		"sys_sunset":        mod.current.Sys.Sunset,
		"base":              mod.current.Base,
		"main_temp":         mod.current.Main.Temp,
		"main_temp_min":     mod.current.Main.TempMin,
		"main_temp_max":     mod.current.Main.TempMax,
		"main_pressure":     mod.current.Main.Pressure,
		"main_sealevel":     mod.current.Main.SeaLevel,
		"main_grndlevel":    mod.current.Main.GrndLevel,
		"main_humidity":     mod.current.Main.Humidity,
		"wind_speed":        mod.current.Wind.Speed,
		"wind_deg":          mod.current.Wind.Deg,
		"clouds_all":        mod.current.Clouds.All,
		"rain":              mod.current.Rain,
		"snow":              mod.current.Snow,
		"dt":                mod.current.Dt,
		"id":                mod.current.ID,
		"name":              mod.current.Name,
		"cod":               mod.current.Cod,
		"unit":              mod.current.Unit,
		"lang":              mod.current.Lang,
		"key":               mod.current.Key,
	}

	mod.evchan <- mod.CreateEvent("current_weather", properties)

	for _, v := range mod.current.Weather {
		mod.TriggerWeatherInformationEvent(&v)
	}
}

// WeatherInformationEvent triggers a weather event
func (mod *OpenweathermapBee) TriggerWeatherInformationEvent(v *owm.Weather) {
	properties := map[string]interface{}{
		"id":          v.ID,
		"main":        v.Main,
		"description": v.Description,
		"icon":        v.Icon,
	}
	mod.evchan <- mod.CreateEvent("main_weather", properties)
}

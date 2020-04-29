/*
 *    Copyright (C) 2020 Sergio Rubio
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
 *      Sergio Rubio <sergio@rubio.im>
 */

// Package sunbee is a bee that sends sunrise/sunset event based on the
// selected location.
package sunbee

import (
	"math"
	"strconv"
	"time"

	"github.com/kelvins/sunrisesunset"
	"github.com/muesli/beehive/bees"
	"github.com/muesli/gominatim"
)

// SunBee is an example for a Bee skeleton, designed to help you get started
// with writing your own Bees.
type SunBee struct {
	bees.Bee
	query    string
	offset   int64
	lat, lon float64
}

// Run executes the Bee's event loop.
func (mod *SunBee) Run(eventChan chan bees.Event) {
	gominatim.SetServer("https://nominatim.openstreetmap.org/")
	qry := gominatim.SearchQuery{
		Q: mod.query,
	}
	resp, err := qry.Get()
	if err != nil {
		mod.LogFatal("Error geocoding %s. err: %v", mod.query, err)
	}

	mod.lat, err = strconv.ParseFloat(resp[0].Lat, 64)
	if err != nil {
		mod.LogFatal("failed parsing latitude from response. err: %v", err)
	}
	mod.lon, err = strconv.ParseFloat(resp[0].Lon, 64)
	if err != nil {
		mod.LogFatal("failed parsing longitude from response. err: %v", err)
	}

	for {
		select {
		case <-mod.SigChan:
			return
		case <-time.After(time.Duration(1 * time.Minute)):
			mod.check(eventChan)
		}
	}
}

// Action triggers the action passed to it.
func (mod *SunBee) Action(action bees.Action) []bees.Placeholder {
	return []bees.Placeholder{}
}

// ReloadOptions parses the config options and initializes the Bee.
func (mod *SunBee) ReloadOptions(options bees.BeeOptions) {
	mod.SetOptions(options)

	mod.ContextSet("sunset", false)
	mod.ContextSet("sunrise", false)
	options.Bind("query", &mod.query)
	options.Bind("offset", mod.offset)
}

func (mod *SunBee) sunset(secondsTo int64, eventChan chan bees.Event) {
	if mod.ContextValue("sunset").(bool) {
		return
	}

	ev := bees.Event{
		Bee:     mod.Name(),
		Name:    "sunset",
		Options: []bees.Placeholder{},
	}

	eventChan <- ev
	mod.ContextSet("sunset", true)
	mod.ContextSet("sunrise", false)
}

func (mod *SunBee) sunrise(secondsTo int64, eventChan chan bees.Event) {
	if mod.ContextValue("sunrise").(bool) {
		return
	}

	ev := bees.Event{
		Bee:     mod.Name(),
		Name:    "sunrise",
		Options: []bees.Placeholder{},
	}

	eventChan <- ev
	mod.ContextSet("sunset", false)
	mod.ContextSet("sunrise", true)
}

func (mod *SunBee) check(eventChan chan bees.Event) {
	now := time.Now().UTC()
	p := sunrisesunset.Parameters{
		Latitude:  mod.lat,
		Longitude: mod.lon,
		UtcOffset: 0,
		Date:      now,
	}

	// Calculate the sunrise and sunset times
	sunrise, sunset, err := p.GetSunriseSunset()
	if err != nil {
		mod.LogFatal("error retrieving sunrise/sunset time. err: %v", err)
	}

	tsunset := time.Date(now.Year(), now.Month(), now.Day(), sunset.Hour(), sunset.Minute(), 0, 0, time.UTC)
	tsunrise := time.Date(now.Year(), now.Month(), now.Day(), sunrise.Hour(), sunrise.Minute(), 0, 0, time.UTC)

	timeToSunset := tsunset.Unix() - now.Unix()
	// sunset already passed
	if timeToSunset < 0 {
		timeToSunset = timeToSunset + 24*60*60
	}
	timeToSunrise := tsunrise.Unix() - now.Unix()
	// sunrise already passed
	if timeToSunrise < 0 {
		timeToSunrise = timeToSunrise + 24*60*60
	}

	f := mod.sunrise
	var evt string
	var tdiff int64
	if timeToSunset < timeToSunrise {
		tdiff = timeToSunset
		f = mod.sunset
		evt = "sunset"
	} else {
		tdiff = timeToSunrise
		evt = "sunrise"
	}
	timeTo := float64(tdiff) / 3600.0
	mod.LogDebugf("Time remaining to %s event in %s (%f, %f): %.2f hours\n", evt, mod.query, mod.lat, mod.lon, timeTo)

	// if sunrise/sunset less than mod.offset seconds away, callback
	if math.Abs(float64(tdiff)) <= float64(mod.offset) {
		f(tdiff, eventChan)
	}
}

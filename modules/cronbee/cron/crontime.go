/*
 *        Copyright (C) 2014 Stefan 'glaxx' Luecke
 *
 *        This program is free software: you can redistribute it and/or modify
 *        it under the terms of the GNU Affero General Public License as published
 *        by the Free Software Foundation, either version 3 of the License, or
 *        (at your option) any later version.
 *
 *        This program is distributed in the hope that it will be useful,
 *        but WITHOUT ANY WARRANTY; without even the implied warranty of
 *        MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 *        GNU Affero General Public License for more details.
 *
 *        You should have received a copy of the GNU Affero General Public License
 *        along with this program.      If not, see <http://www.gnu.org/licenses/>.
 *
 *        Authors:	Stefan Luecke <glaxx@glaxx.net>
 */

 /*
 * TODO List:
 * - Test leap year behavior 
 * - Test
 */

package cron

import (
//	"fmt"
	"time"
//	"sort"
	"log"
)

type crontime struct {
	second []int
	minute []int
	hour   []int
	dow    []int //Day of Week
	dom    []int //Day of Month
	month  []int
	CalculatedTime time.Time
}

// This functions returns a time.Duration until the next Event based
// upon the parsed input.
func (c *crontime) NextEvent() time.Duration{
	c.CalculatedTime = time.Now() // Ignore all Events in the Past & initial 'result'
	c.CalculatedTime = setNanoecond(c.CalculatedTime, 10000)
	c.nextValidMonth()
	c.nextValidDay()
	c.nextValidHour()
	c.nextValidMinute()
	c.nextValidSecond()
	log.Println("Cronbee has found the next time stamp: ", c.CalculatedTime)
	return c.CalculatedTime.Sub(time.Now())
}

// Calculates the next valid Month based upon the previous results.
func (c *crontime) nextValidMonth() {
	for _, mon := range c.month {
		if time.Now().Year() == c.CalculatedTime.Year() {
			if !hasPassed(mon, int(c.CalculatedTime.Month())) {
				c.CalculatedTime = setMonth(c.CalculatedTime, mon)
				return
			}
		} else {
			c.CalculatedTime = setMonth(c.CalculatedTime, mon)
			return
		}
	}
	// If no result was found try it again in the following year
	c.CalculatedTime = c.CalculatedTime.AddDate(1, 0, 0)
	c.nextValidMonth()
}

// Calculates the next valid Day based upon the previous results.
func (c *crontime) nextValidDay() {
	for _, dom := range c.dom {
		if c.CalculatedTime.Month() == time.Now().Month() {
			if !hasPassed(dom, c.CalculatedTime.Day()) {
				for _, dow := range c.dow {
					if monthHasDow(dow, dom, int(c.CalculatedTime.Month()), c.CalculatedTime.Year()){
						c.CalculatedTime = setDay(c.CalculatedTime, dom)
						return
					}
				}
			}
		} else {
			for _, dow := range c.dow {
				if monthHasDow(dow, dom, int(c.CalculatedTime.Month()), c.CalculatedTime.Year()){
					c.CalculatedTime = setDay(c.CalculatedTime, dom)
					return
				}
			}
		}
	}
	// If no result was found try it again in the following month.
	c.CalculatedTime = c.CalculatedTime.AddDate(0, 1, 0)
	c.nextValidMonth()
	c.nextValidDay()
}

// Calculates the next valid Hour based upon the previous results.
func (c *crontime) nextValidHour() {
	for _, hour := range c.hour {
		if c.CalculatedTime.Day() == time.Now().Day() {
			if !hasPassed(hour, c.CalculatedTime.Hour()) {
				c.CalculatedTime = setHour(c.CalculatedTime, hour)
				return
			}
		} else {
			c.CalculatedTime = setHour(c.CalculatedTime, hour)
			return
		}
	}
	// If no result was found try it again in the following day.
	c.CalculatedTime = c.CalculatedTime.AddDate(0, 0, 1)
	c.nextValidDay()
	c.nextValidHour()
}

// Calculates the next valid Minute based upon the previous results.
func (c *crontime) nextValidMinute() {
	for _, min := range c.minute {
		if c.CalculatedTime.Hour() == time.Now().Hour() {
			if !hasPassed(min, c.CalculatedTime.Minute()) {
				c.CalculatedTime = setMinute(c.CalculatedTime, min)
				return
			}
		} else {
			c.CalculatedTime = setMinute(c.CalculatedTime, min)
			return
		}
	}
	c.CalculatedTime = c.CalculatedTime.Add(1 * time.Hour)
	c.nextValidHour()
	c.nextValidMinute()
}

// Calculates the next valid Second based upon the previous results.
func (c *crontime) nextValidSecond() {
	for _, sec := range c.second {
		if c.CalculatedTime.Minute() == time.Now().Minute() {
			// check if sec is in the past. <= prevents triggering the same event twice
			if !(sec <= c.CalculatedTime.Second()){
				c.CalculatedTime = setSecond(c.CalculatedTime, sec)
				return
			}
		} else {
			c.CalculatedTime = setSecond(c.CalculatedTime, sec)
			return
		}
	}
	c.CalculatedTime = c.CalculatedTime.Add(1 * time.Minute)
	c.nextValidMinute()
	c.nextValidSecond()
}

func hasPassed(value, tstamp int) bool{
	return value < tstamp
}

// Check if the combination of day(of month), month and year is the weekday dow.
func monthHasDow(dow, dom, month, year int) bool{
	Nday := dom % 7
	var Nmonth int
	switch month{
		case 1: Nmonth = 0
		case 2: Nmonth = 3
		case 3: Nmonth = 3
		case 4: Nmonth = 6
		case 5: Nmonth = 1
		case 6: Nmonth = 4
		case 7: Nmonth = 6
		case 8: Nmonth = 2
		case 9: Nmonth = 5
		case 10: Nmonth = 0
		case 11: Nmonth = 3
		case 12: Nmonth = 5
	}
	var Nyear int
	temp := year % 100
	if temp != 0{
		Nyear = (temp + (temp / 4)) % 7	
	} else {
		Nyear = 0
	}
	Ncent := (3 - ((year / 100) %4)) * 2
	var Nsj int
	if isLeapYear(year) {
		Nsj = -1
	} else {
		Nsj = 0
	}
	W := (Nday + Nmonth + Nyear + Ncent + Nsj) % 7
	return dow == W
}

func isLeapYear(year int) bool{
	return year % 4 == 0 && (year % 100 != 0 || year % 400 == 0)
}

func setMonth(tstamp time.Time, month int) time.Time {
	if month >= 12 { panic("ERROR") }
	return tstamp.AddDate(0, -absolute(int(tstamp.Month()), month), 0)
}

func setDay(tstamp time.Time, day int) time.Time {
	if day >= 31 { panic("ERROR") }
	return tstamp.AddDate(0, 0, -absolute(tstamp.Day(), day))
}

func setHour(tstamp time.Time, hour int) time.Time {
	if hour >= 24 { panic("ERROR") }
	return tstamp.Add(time.Duration(-absolute(tstamp.Hour(), hour)) * time.Hour)
}

func setMinute(tstamp time.Time, minute int) time.Time {
	if minute >= 60 { panic("ERROR") }
	return tstamp.Add(time.Duration(-absolute(tstamp.Minute(), minute)) * time.Minute)
}

func setSecond(tstamp time.Time, second int) time.Time {
	if second >= 60 { panic("ERROR") }
	return tstamp.Add(time.Duration(-absolute(tstamp.Second(), second)) * time.Second)
}

func setNanoecond(tstamp time.Time, nanosecond int) time.Time {
	return tstamp.Add(time.Duration(-absolute(tstamp.Nanosecond(), nanosecond)) * time.Nanosecond)
}

func absolute(a, b int) int {
	return a - b
}
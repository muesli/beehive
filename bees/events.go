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

// Package bees is Beehive's central module system.
package bees

import (
	"fmt"
	"runtime/debug"

	log "github.com/sirupsen/logrus"
)

// An Event describes an event including its parameters.
type Event struct {
	Bee     string
	Name    string
	Options Placeholders
}

var (
	eventsIn = make(chan Event)
)

// handleEvents handles incoming events and executes matching Chains.
func handleEvents() {
	for {
		event, ok := <-eventsIn
		if !ok {
			log.Println()
			log.Println("Stopped event handler!")
			break
		}

		bee := GetBee(event.Bee)
		(*bee).LogEvent()

		log.Debugln()
		log.Debugln("Event received:", event.Bee, "/", event.Name, "-", GetEventDescriptor(&event).Description)
		for _, v := range event.Options {
			vv := truncateString(fmt.Sprintln(v), 1000)
			log.Debugln("\tOptions:", vv)
		}

		go func() {
			defer func() {
				if e := recover(); e != nil {
					log.Printf("Fatal chain event: %s %s", e, debug.Stack())
				}
			}()

			execChains(&event)
		}()
	}
}

func truncateString(str string, num int) string {
	bnoden := str
	if len(str) > num {
		if num > 3 {
			num -= 3
		}
		bnoden = str[0:num] + "... (" + fmt.Sprint(len(str)-num) + " more characters)"
	}
	return bnoden
}

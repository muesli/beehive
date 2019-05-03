/*
 *    Copyright (C) 2017 Christian Muehlhaeuser
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
	"sort"
	"sync"
	"time"
)

// LogMessage stores a log message with its timestamp, type and originating Bee
type LogMessage struct {
	ID          string
	Bee         string
	Message     string
	MessageType uint
	Timestamp   time.Time
}

var (
	logs     = make(map[string][]LogMessage)
	logMutex sync.RWMutex
)

// MessageType defines the log level of the log entry we're dealing with
type MessageType uint

const (
	// LogInfo is for info-level log entries
	LogInfo MessageType = iota

	// LogError is for error-level log entries
	LogError MessageType = iota

	// LogFatal is for fatal-level log entries
	LogFatal MessageType = iota

	// LogDebug is for debug-level log entries
	LogDebug MessageType = iota
)

// LogSorter is used for sorting an array of LogMessages by their timestamp
type LogSorter []LogMessage

func (a LogSorter) Len() int           { return len(a) }
func (a LogSorter) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a LogSorter) Less(i, j int) bool { return !a[i].Timestamp.Before(a[j].Timestamp) }

// NewLogMessage returns a newly composed LogMessage
func NewLogMessage(bee string, message string, messageType MessageType) LogMessage {
	return LogMessage{
		ID:          UUID(),
		Bee:         bee,
		Message:     message,
		MessageType: uint(messageType),
		Timestamp:   time.Now(),
	}
}

// Log adds a new LogMessage to the log
func Log(bee string, message string, messageType MessageType) {
	logMutex.Lock()
	defer logMutex.Unlock()

	logs[bee] = append(logs[bee], NewLogMessage(bee, message, messageType))
}

// GetLogs returns all logs for a Bee.
func GetLogs(bee string) []LogMessage {
	r := []LogMessage{}

	logMutex.RLock()
	for b, ls := range logs {
		if len(bee) == 0 || bee == b {
			for _, l := range ls {
				r = append(r, l)
			}
		}
	}
	logMutex.RUnlock()

	sort.Sort(LogSorter(r))
	return r
}

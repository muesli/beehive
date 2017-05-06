/*
 *    Copyright (C) 2017      Sergio Rubio
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

package fsnotifybee

import (
	"github.com/fsnotify/fsnotify"
	"github.com/muesli/beehive/bees"
)

type FSNotifyBee struct {
	bees.Bee
}

func (mod *FSNotifyBee) Run(eventChan chan bees.Event) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		mod.LogFatal("error: could not create the fswatcher: %v", err)
	}
	defer watcher.Close()

	path := mod.Options().Value("path").(string)
	mod.Logf("monitoring %s\n", path)

	err = watcher.Add(path)
	if err != nil {
		mod.LogErrorf("error: watching %s failed: %v", path, err)
		return
	}
	for {
		select {
		case <-mod.SigChan:
			return
		case event := <-watcher.Events:
			if event.Op != fsnotify.Write && event.Op != 0 {
				sendEvent(mod.Name(), event.Op.String(), event.Name, eventChan)
			}
		case <-watcher.Errors:
		}
	}
}

func (mod *FSNotifyBee) ReloadOptions(options bees.BeeOptions) {
	mod.SetOptions(options)
}

func sendEvent(bee, etype, path string, eventChan chan bees.Event) {
	event := bees.Event{
		Bee:  bee,
		Name: "event",
		Options: []bees.Placeholder{
			{
				Name:  "type",
				Type:  "string",
				Value: etype,
			},
			{
				Name:  "path",
				Type:  "string",
				Value: path,
			},
		},
	}
	eventChan <- event
}

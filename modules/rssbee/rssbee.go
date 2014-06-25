/*
 *    Copyright (C) 2014 Michael Wendland
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
 *      Michael Wendland <michiwend@michiwend.com>
 */

// RSS module for beehive
package rssbee

import (
	"github.com/muesli/beehive/modules"
//    rss "github.com/jteeuwen/go-pkg-rss"
)

var (
    eventChan chan modules.Event
)

type RSSBee struct {
}

func (mod *RSSBee) Name() string {
	return "rssbee"
}

func (mod *RSSBee) Description() string {
	return "A bee that manages RSS-feeds"
}

func (mod *RSSBee) Events() []modules.EventDescriptor {
	events := []modules.EventDescriptor{
        modules.EventDescriptor{
            Namespace: mod.Name(),
            Name: "incoming",
            Description: "A new post has been received through the Feed",
            Options: []modules.PlaceholderDescriptor{
            },
        },
    }
	return events
}

func (mod *RSSBee) Actions() []modules.ActionDescriptor {
	actions := []modules.ActionDescriptor{}
	return actions
}

func (mod *RSSBee) Run(eventChan chan modules.Event) {
//	eventChan <- ev
}

func (mod *RSSBee) Action(action modules.Action) []modules.Placeholder {
	return []modules.Placeholder{}
}

func init() {
	bee := RSSBee{}

	modules.RegisterModule(&bee)
}

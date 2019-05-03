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

// Package travisbee is a bee for monitoring and reacting to the status of
// TravisCI builds.
package travisbee

import (
	"context"
	"time"

	"github.com/muesli/beehive/bees"
	"github.com/shuheiktgw/go-travis"
)

// TravisBee is a bee for monitoring and reacting to the status of
// TravisCI builds.
type TravisBee struct {
	bees.Bee

	eventChan chan bees.Event
	client    *travis.Client
	apiToken  string
	builds    map[uint]BuildTracker
}

// BuildTracker is a marker struct to denote "tracked builds" that TravisBee
// cares about.
type BuildTracker struct {
	id       uint
	state    string
	lastTime time.Time
}

// Run executes the Bee's event loop.
func (mod *TravisBee) Run(eventChan chan bees.Event) {

	mod.eventChan = eventChan
	mod.client = travis.NewClient(travis.ApiOrgUrl, mod.apiToken)
	mod.builds = make(map[uint]BuildTracker)

	since := time.Now()
	timeout := time.Duration(time.Second * 10)
	for {
		select {
		case <-mod.SigChan:
			return
		case <-time.After(timeout):
			mod.getBuilds(since)

		}
		since = time.Now()
		timeout = time.Duration(time.Minute)
	}
}

func (mod *TravisBee) getBuilds(since time.Time) {
	var opt = travis.BuildsOption{Limit: 10, Offset: 0, SortBy: "started_at:desc"}
	builds, _, err := mod.client.Builds.List(context.Background(), &opt)
	if err != nil {
		mod.LogErrorf("Error getting builds from travis-ci: %v", err)
		panic("Fatal error processing travis builds!")
	}

	for i, currentBuild := range builds {
		if b, ok := mod.builds[currentBuild.Id]; ok {
			b.lastTime = time.Now()
			if b.state != currentBuild.State {
				mod.handleStateChange(&b, currentBuild)
				b.state = currentBuild.State
			}
			mod.builds[currentBuild.Id] = b
		} else {
			mod.Logf("[%d] %d - %s / %s", i, currentBuild.Id, currentBuild.StartedAt, currentBuild.State)

			// If a build was just barely created (i.e. hasn't moved into any other
			// stage yet) then we can't even get a started_at time - instead, let's
			// just make it a new build, figuring it's obviously something we havent
			// seen before.
			if currentBuild.State == "created" {
				mod.handleNewBuild(currentBuild, since)
			} else {
				t, err := time.Parse(time.RFC3339, currentBuild.StartedAt)
				if err != nil {
					mod.LogErrorf("Error parsing time %s - %v", currentBuild.StartedAt, err)
				}
				if t.After(since) {
					mod.handleNewBuild(currentBuild, since)
				}
			}
		}
	}
}

func (mod *TravisBee) handleNewBuild(build *travis.Build, since time.Time) {
	mod.builds[build.Id] = BuildTracker{id: build.Id, state: build.State, lastTime: time.Now()}
	mod.Logf("Tracking build %d - state %s", build.Id, build.State)

	ev := bees.Event{
		Bee:  mod.Name(),
		Name: "build_started",
		Options: []bees.Placeholder{
			{
				Name:  "id",
				Type:  "uint",
				Value: build.Id,
			},
			{
				Name:  "state",
				Type:  "string",
				Value: build.State,
			},
			{
				Name:  "repo_slug",
				Type:  "string",
				Value: build.Repository.Slug,
			},
			{
				Name:  "duration",
				Type:  "uint",
				Value: build.Duration,
			},
		},
	}

	mod.eventChan <- ev
}

func (mod *TravisBee) handleStateChange(bt *BuildTracker, build *travis.Build) {
	mod.Logf("State changed! (was %s, is now %s)", bt.state, build.State)

	ev := bees.Event{
		Bee:  mod.Name(),
		Name: "build_status_change",
		Options: []bees.Placeholder{
			{
				Name:  "id",
				Type:  "uint",
				Value: build.Id,
			},
			{
				Name:  "state",
				Type:  "string",
				Value: build.State,
			},
			{
				Name:  "last_state",
				Type:  "string",
				Value: bt.state,
			},
			{
				Name:  "repo_slug",
				Type:  "string",
				Value: build.Repository.Slug,
			},
			{
				Name:  "duration",
				Type:  "uint",
				Value: build.Duration,
			},
		},
	}

	mod.eventChan <- ev

	if build.State == "canceled" || build.State == "passed" ||
		build.State == "failed" || build.State == "errored" {
		mod.handleBuildFinish(bt, build)
	}
}

func (mod *TravisBee) handleBuildFinish(bt *BuildTracker, build *travis.Build) {
	mod.Logf("Build %d has finished with state %s", build.Id, build.State)
	ev := bees.Event{
		Bee:  mod.Name(),
		Name: "build_finished",
		Options: []bees.Placeholder{
			{
				Name:  "id",
				Type:  "uint",
				Value: build.Id,
			},
			{
				Name:  "state",
				Type:  "string",
				Value: build.State,
			},
			{
				Name:  "repo_slug",
				Type:  "string",
				Value: build.Repository.Slug,
			},
			{
				Name:  "duration",
				Type:  "uint",
				Value: build.Duration,
			},
		},
	}
	mod.eventChan <- ev
	mod.Logf("Now untracking build %d", bt.id)
	delete(mod.builds, bt.id)
}

// ReloadOptions parses the config options and initializes the Bee.
func (mod *TravisBee) ReloadOptions(options bees.BeeOptions) {
	mod.SetOptions(options)

	options.Bind("api_key", &mod.apiToken)
}

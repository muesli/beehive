/*
 *    Copyright (C) 2014      Daniel 'grindhold' Brendle
 *                  2014-2017 Christian Muehlhaeuser
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
 *      Daniel 'grindhold' Brendle <grindhold@skarphed.org>
 *      Christian Muehlhaeuser <muesli@gmail.com>
 */

// Package jenkinsbee is a Bee that can interface with a Jenkins server.
package jenkinsbee

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"
	"github.com/muesli/beehive/bees"
)

// JenkinsBee is a Bee that can interface with a Jenkins server.
type JenkinsBee struct {
	bees.Bee

	url      string
	user     string
	password string

	Jobs      map[string]Job `json:"jobs"`
	eventChan chan bees.Event
}

type report struct {
	Jobs []Job `json:"jobs"`
}

// Job represents the JSON API response for a Jenkins job.
type Job struct {
	Name      string    `json:"name"`
	URL       string    `json:"url"`
	Color     string    `json:"color"`
	LastBuild LastBuild `json:"lastBuild"`
}

type LastBuild struct {
	Building bool `json:"building"`
}

func (mod *JenkinsBee) announceStatusChange(j Job) {
	event := bees.Event{
		Bee:  mod.Name(),
		Name: "status_change",
		Options: []bees.Placeholder{
			{
				Name:  "name",
				Type:  "string",
				Value: j.Name,
			},
			{
				Name:  "url",
				Type:  "string",
				Value: j.URL,
			},
			{
				Name:  "status",
				Type:  "string",
				Value: j.Color,
			},
		},
	}
	mod.eventChan <- event
}

func (mod *JenkinsBee) announceBuildStateChange(j Job) {
	event := bees.Event{
		Bee:  mod.Name(),
		Name: "build_status_change",
		Options: []bees.Placeholder{
			{
				Name:  "name",
				Type:  "string",
				Value: j.Name,
			},
			{
				Name:  "url",
				Type:  "string",
				Value: j.URL,
			},
			{
				Name:  "status",
				Type:  "string",
				Value: j.Color,
			},
			{
				Name:  "building",
				Type:  "bool",
				Value: j.LastBuild.Building,
			},
		},
	}
	mod.eventChan <- event
}

func (mod *JenkinsBee) poll() {
	request, err := http.NewRequest("GET", mod.url+"/api/json?tree=jobs[name,url,color,lastBuild[building]]", nil)
	if err != nil {
		mod.LogErrorf("Could not build request: %s", err)
		return
	}
	request.SetBasicAuth(mod.user, mod.password)

	client := http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		mod.LogErrorf("Could not call API on %s/api/json: %s", mod.url, err)
		return
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		mod.LogErrorf("Could not read data from API: %s", err)
		return
	}
	rep := new(report)
	err = json.Unmarshal(body, &rep)
	if err != nil {
		mod.LogErrorf("Failed to unmarshal JSON: %s", err)
		return
	}

	jobmap := make(map[string]Job)
	for job := range rep.Jobs {
		if oldState, ok := mod.Jobs[rep.Jobs[job].Name]; !ok {
			// There is no record of this job
			mod.announceStatusChange(rep.Jobs[job])
		} else {
			// There exists a record of this job
			if oldState.Color != rep.Jobs[job].Color {
				// The status is different from last time
				mod.announceStatusChange(rep.Jobs[job])
			}
			if oldState.LastBuild.Building != rep.Jobs[job].LastBuild.Building {
				mod.announceBuildStateChange(rep.Jobs[job])
			}
		}
		jobmap[rep.Jobs[job].Name] = rep.Jobs[job]
	}
	mod.Jobs = jobmap
}

// Run executes the Bee's event loop.
func (mod *JenkinsBee) Run(cin chan bees.Event) {
	mod.eventChan = cin
	for {
		select {
		case <-mod.SigChan:
			return

		case <-time.After(time.Duration(10 * time.Second)):
			mod.poll()
		}
	}
}

func (mod *JenkinsBee) triggerBuild(jobname string) {
	client := http.Client{}
	request, err := http.NewRequest("GET", mod.url+"/job/"+jobname+"/build", nil)
	if err != nil {
		mod.LogErrorf("Could not build request: %s", err)
		return
	}
	request.SetBasicAuth(mod.user, mod.password)
	if _, err := client.Do(request); err != nil {
		mod.LogErrorf("Could not trigger build: %s", err)
	}
}

// Action triggers the action passed to it.
func (mod *JenkinsBee) Action(action bees.Action) []bees.Placeholder {
	outs := []bees.Placeholder{}
	switch action.Name {
	case "trigger":
		jobname := ""
		action.Options.Bind("job", &jobname)

		mod.triggerBuild(jobname)

	default:
		panic("Unknown action triggered in " + mod.Name() + ": " + action.Name)
	}
	return outs
}

// ReloadOptions parses the config options and initializes the Bee.
func (mod *JenkinsBee) ReloadOptions(options bees.BeeOptions) {
	mod.SetOptions(options)

	options.Bind("url", &mod.url)
	options.Bind("user", &mod.user)
	options.Bind("password", &mod.password)
}

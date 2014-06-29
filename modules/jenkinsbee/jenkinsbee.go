/*
 *    Copyright (C) 2014 Daniel 'grindhold' Brendle
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
 */

package jenkinsbee

import (
	"encoding/json"
	"github.com/muesli/beehive/modules"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type JenkinsBee struct {
	modules.Module

	url       string
	Jobs      map[string]Job `json:"jobs"`
	eventChan chan modules.Event
}

type report struct {
	Jobs []Job `json:"jobs"`
}

type Job struct {
	Name  string `json:"name"`
	Url   string `json:"url"`
	Color string `json:"color"`
}

func (mod *JenkinsBee) announceStatusChange(j Job) {
	event := modules.Event{
		Bee:  mod.Name(),
		Name: "statuschange",
		Options: []modules.Placeholder{
			modules.Placeholder{
				Name:  "name",
				Type:  "string",
				Value: j.Name,
			},
			modules.Placeholder{
				Name:  "url",
				Type:  "string",
				Value: j.Url,
			},
			modules.Placeholder{
				Name:  "status",
				Type:  "string",
				Value: j.Color,
			},
		},
	}
	mod.eventChan <- event
}

func (mod *JenkinsBee) Run(cin chan modules.Event) {
	mod.eventChan = cin
	for {
		resp, err := http.Get(mod.url + "/api/json")
		if err != nil {
			log.Println("Could not call API on " + mod.url + "/api/json")
			time.Sleep(5 * time.Second)
			continue
		}
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Println("Could not read data of API-Call")
			time.Sleep(5 * time.Second)
			continue
		}
		rep := new(report)
		err = json.Unmarshal(body, &rep)
		if err != nil {
			log.Println("Failed to unmarshal JSON")
			time.Sleep(5 * time.Second)
			continue
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
			}
			jobmap[rep.Jobs[job].Name] = rep.Jobs[job]
		}
		mod.Jobs = jobmap
		time.Sleep(5 * time.Second)
	}
}

func (mod *JenkinsBee) triggerBuild(jobname string) {
	if _, err := http.Get(mod.url + "/job/" + jobname + "/build"); err != nil {
		log.Println("Could not trigger build")
	}
	return
}

func (mod *JenkinsBee) Action(action modules.Action) []modules.Placeholder {
	outs := []modules.Placeholder{}
	switch action.Name {
	case "trigger":
		jobname := ""
		for _, opt := range action.Options {
			if opt.Name == "job" {
				jobname = opt.Value.(string)
			}
		}
		mod.triggerBuild(jobname)
	default:
		return outs
	}
	return outs
}

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
	"github.com/muesli/beehive/bees"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type JenkinsBee struct {
	bees.Module

	url       string
	user     string
	password string

	Jobs      map[string]Job `json:"jobs"`
	eventChan chan bees.Event
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
	event := bees.Event{
		Bee:  mod.Name(),
		Name: "statuschange",
		Options: []bees.Placeholder{
			bees.Placeholder{
				Name:  "name",
				Type:  "string",
				Value: j.Name,
			},
			bees.Placeholder{
				Name:  "url",
				Type:  "string",
				Value: j.Url,
			},
			bees.Placeholder{
				Name:  "status",
				Type:  "string",
				Value: j.Color,
			},
		},
	}
	mod.eventChan <- event
}

func (mod *JenkinsBee) Run(cin chan bees.Event) {
	mod.eventChan = cin
	for {
		select {
			case <-mod.SigChan:
				return

			default:
		}
		time.Sleep(10 * time.Second)

		request, err := http.NewRequest("GET", mod.url + "/api/json", nil)
		if err != nil {
			log.Println("Could not build request")
			break
		}
		request.SetBasicAuth(mod.user, mod.password)

		client := http.Client{}
		resp, err := client.Do(request)
		if err != nil {
			log.Println("Could not call API on " + mod.url + "/api/json", err)
			continue
		}

		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Println("Could not read data of API-Call")
			continue
		}
		rep := new(report)
		err = json.Unmarshal(body, &rep)
		if err != nil {
			log.Println("Failed to unmarshal JSON")
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
	}
}

func (mod *JenkinsBee) triggerBuild(jobname string) {
	client := http.Client{}
	request, err := http.NewRequest("GET", mod.url + "/job/" + jobname + "/build", nil)
	if err != nil {
		log.Println("Could not build request")
		return
	}
	request.SetBasicAuth(mod.user, mod.password)
	if _, err := client.Do(request); err != nil {
		log.Println("Could not trigger build")
	}
}

func (mod *JenkinsBee) Action(action bees.Action) []bees.Placeholder {
	outs := []bees.Placeholder{}
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
		panic("Unknown action triggered in " +mod.Name()+": "+action.Name)
	}
	return outs
}

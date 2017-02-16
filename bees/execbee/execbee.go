/*
 *    Copyright (C) 2015      Dominik Schmidt
 *                  2015-2017 Christian Muehlhaeuser
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
 *      Dominik Schmidt <domme@tomahawk-player.org>
 *      Christian Muehlhaeuser <muesli@gmail.com>
 */

// Package execbee is a Bee that can launch external processes.
package execbee

import (
	"bufio"
	"os"
	"os/exec"
	"strings"

	"github.com/muesli/beehive/bees"
)

// ExecBee is a Bee that can launch external processes.
type ExecBee struct {
	bees.Bee

	eventChan chan bees.Event
}

// Action triggers the action passed to it.
func (mod *ExecBee) Action(action bees.Action) []bees.Placeholder {
	outs := []bees.Placeholder{}

	switch action.Name {
	case "execute":
		var command string
		action.Options.Bind("command", &command)
		mod.Logln("Executing locally: ", command)

		go func() {
			c := strings.Split(command, " ")
			cmd := exec.Command(c[0], c[1:]...)

			// read and print stdout
			outReader, err := cmd.StdoutPipe()
			if err != nil {
				mod.LogFatal("Error creating StdoutPipe for Cmd", err)
				return
			}
			outBuffer := []string{}
			outScanner := bufio.NewScanner(outReader)
			go func() {
				for outScanner.Scan() {
					foo := outScanner.Text()
					mod.Logln("execbee: std: | ", foo)
					outBuffer = append(outBuffer, foo)
				}
			}()

			// read and print stderr
			errReader, err := cmd.StderrPipe()
			if err != nil {
				mod.LogFatal(os.Stderr, "Error creating StderrPipe for Cmd", err)
				return
			}
			errBuffer := []string{}
			errScanner := bufio.NewScanner(errReader)
			go func() {
				for errScanner.Scan() {
					foo := errScanner.Text()
					mod.Logln("Err: | ", foo)
					errBuffer = append(errBuffer, foo)
				}
			}()

			err = cmd.Start()
			if err != nil {
				mod.LogFatal("Error starting Cmd", err)
			}

			err = cmd.Wait()
			if err != nil {
				mod.LogFatal("Error waiting for Cmd", err)
			}

			ev := bees.Event{
				Bee:  mod.Name(),
				Name: "result",
				Options: []bees.Placeholder{
					{
						Name:  "stdout",
						Type:  "string",
						Value: strings.Join(outBuffer, "\n"),
					},
					{
						Name:  "stderr",
						Type:  "string",
						Value: strings.Join(errBuffer, "\n"),
					},
				},
			}
			mod.eventChan <- ev
		}()

	default:
		panic("Unknown action triggered in " + mod.Name() + ": " + action.Name)
	}

	return outs
}

// Run executes the Bee's event loop.
func (mod *ExecBee) Run(eventChan chan bees.Event) {
	mod.eventChan = eventChan
}

// ReloadOptions parses the config options and initializes the Bee.
func (mod *ExecBee) ReloadOptions(options bees.BeeOptions) {
	mod.SetOptions(options)
}

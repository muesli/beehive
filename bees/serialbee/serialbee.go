// +build !darwin

// FIXME: This bee doesn't build on macOS
// XXX: https://gist.github.com/prologic/b5fa148410a26f917b4d944b72d847c0

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

// Package serialbee is a Bee that can send & receive data on a serial port.
package serialbee

import (
	"bytes"
	"encoding/binary"
	"io"

	"github.com/huin/goserial"

	"github.com/muesli/beehive/bees"
)

// SerialBee is a Bee that can send & receive data on a serial port.
type SerialBee struct {
	bees.Bee

	conn io.ReadWriteCloser

	device   string
	baudrate int
}

// Action triggers the action passed to it.
func (mod *SerialBee) Action(action bees.Action) []bees.Placeholder {
	outs := []bees.Placeholder{}
	text := ""

	switch action.Name {
	case "send":
		action.Options.Bind("text", &text)

		bufOut := new(bytes.Buffer)
		err := binary.Write(bufOut, binary.LittleEndian, []byte(text))
		if err != nil {
			panic(err)
		}

		_, err = mod.conn.Write(bufOut.Bytes())
		if err != nil {
			panic(err)
		}

	default:
		panic("Unknown action triggered in " + mod.Name() + ": " + action.Name)
	}

	return outs
}

func (mod *SerialBee) handleEvents(eventChan chan bees.Event) error {
	text := ""
	c := []byte{0}
	for {
		_, err := mod.conn.Read(c)
		if err != nil {
			return err
		}
		if c[0] == 10 || c[0] == 13 {
			break
		}

		text += string(c[0])
	}

	if len(text) > 0 {
		ev := bees.Event{
			Bee:  mod.Name(),
			Name: "message",
			Options: []bees.Placeholder{
				{
					Name:  "port",
					Type:  "string",
					Value: mod.device,
				},
				{
					Name:  "text",
					Type:  "string",
					Value: text,
				},
			},
		}
		eventChan <- ev
	}

	return nil
}

// Run executes the Bee's event loop.
func (mod *SerialBee) Run(eventChan chan bees.Event) {
	if mod.baudrate == 0 || mod.device == "" {
		return
	}

	var err error
	c := &goserial.Config{Name: mod.device, Baud: mod.baudrate}
	mod.conn, err = goserial.OpenPort(c)
	if err != nil {
		mod.LogFatal(err)
	}
	defer mod.conn.Close()

	go func() {
		for {
			if err := mod.handleEvents(eventChan); err != nil {
				return
			}
		}
	}()

	select {
	case <-mod.SigChan:
		return
	}
}

// ReloadOptions parses the config options and initializes the Bee.
func (mod *SerialBee) ReloadOptions(options bees.BeeOptions) {
	mod.SetOptions(options)

	options.Bind("device", &mod.device)
	options.Bind("baudrate", &mod.baudrate)
}

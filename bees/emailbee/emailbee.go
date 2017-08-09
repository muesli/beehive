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

// Package emailbee is a Bee that is able to send emails.
package emailbee

import (
	"net"

	"strconv"

	"github.com/go-gomail/gomail"
	"github.com/muesli/beehive/bees"
)

// EmailBee is a Bee that is able to send emails.
type EmailBee struct {
	bees.Bee

	username string
	password string
	server   string
}

// Action triggers the action passed to it.
func (mod *EmailBee) Action(action bees.Action) []bees.Placeholder {
	outs := []bees.Placeholder{}

	switch action.Name {
	case "send":
		var from, to, plainText, htmlText, subject string
		action.Options.Bind("sender", &from)
		action.Options.Bind("recipient", &to)
		action.Options.Bind("subject", &subject)
		action.Options.Bind("text", &plainText)
		action.Options.Bind("html", &htmlText)

		m := gomail.NewMessage()
		if len(from) > 0 {
			m.SetHeader("From", from)
		} else {
			m.SetHeader("From", mod.username)
		}
		m.SetHeader("To", to)
		m.SetHeader("Subject", subject)
		if plainText != "" {
			m.SetBody("text/plain", plainText)
		}
		if htmlText != "" {
			m.SetBody("text/html", htmlText)
		}

		host, portstr, err := net.SplitHostPort(mod.server)
		if err != nil {
			host = mod.server
			portstr = "587"
		}
		port, _ := strconv.Atoi(portstr)
		s, _ := gomail.NewDialer(host, port, mod.username, mod.password).Dial()

		// Send the email.
		if err := gomail.Send(s, m); err != nil {
			panic(err)
		}

	default:
		panic("Unknown action triggered in " + mod.Name() + ": " + action.Name)
	}

	return outs
}

// ReloadOptions parses the config options and initializes the Bee.
func (mod *EmailBee) ReloadOptions(options bees.BeeOptions) {
	mod.SetOptions(options)

	options.Bind("username", &mod.username)
	options.Bind("password", &mod.password)
	options.Bind("address", &mod.server)
}

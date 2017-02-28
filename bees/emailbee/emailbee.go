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
	"regexp"
	"strconv"
	"strings"

	"github.com/muesli/beehive/bees"
	gomail "gopkg.in/gomail.v2"
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
		var to, mailText, subject string
		var addOn bool
		action.Options.Bind("recipient", &to)
		action.Options.Bind("text", &mailText)
		action.Options.Bind("subject", &subject)
		action.Options.Bind("HTML", &addOn)

		var host string
		var port int
		if strings.Index(mod.server, ":") != -1 {
			host = strings.Split(mod.server, ":")[0]
			port, _ = strconv.Atoi(strings.Split(mod.server, ":")[1])
		} else {
			host = mod.server
			port = 587
		}

		re, _ := regexp.Compile("\\<[\\S\\s]+?\\>")
		plainText := re.ReplaceAllString(mailText, "")

		m := gomail.NewMessage()
		m.SetHeader("From", mod.username)
		m.SetHeader("To", to)
		m.SetHeader("Subject", subject)
		m.SetBody("text/plain", plainText)

		if addOn {
			m.AddAlternative("text/html", mailText)
		}

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

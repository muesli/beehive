/*
 *    Copyright (C) 2014 Christian Muehlhaeuser
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

// beehive's email module.
package emailbee

import (
	_ "log"
	"net/smtp"
	"strings"

	"github.com/muesli/beehive/bees"
)

type EmailBee struct {
	bees.Bee

	username string
	password string
	server   string
}

// Interface impl

func (mod *EmailBee) Action(action bees.Action) []bees.Placeholder {
	outs := []bees.Placeholder{}

	switch action.Name {
	case "send":
		to := ""
		text := ""
		subject := ""

		action.Options.Bind("recipient", &to)
		action.Options.Bind("text", &text)
		action.Options.Bind("subject", &subject)

		text = "Subject: " + subject + "\n\n" + text
		auth := smtp.PlainAuth("", mod.username, mod.password, mod.server[:strings.Index(mod.server, ":")])
		err := smtp.SendMail(mod.server, auth, mod.username, []string{to}, []byte(text))
		if err != nil {
			panic(err)
		}

	default:
		panic("Unknown action triggered in " + mod.Name() + ": " + action.Name)
	}

	return outs
}

func (mod *EmailBee) SetOptions(options bees.BeeOptions) {
	//FIXME: implement this
}

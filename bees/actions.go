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

// Package bees is Beehive's central module system.
package bees

import (
	"bytes"
	"text/template"

	log "github.com/sirupsen/logrus"

	"github.com/muesli/beehive/templatehelper"
)

// Action describes an action.
type Action struct {
	ID      string
	Bee     string
	Name    string
	Options Placeholders
}

var (
	actions []Action
)

// GetActions returns all configured actions.
func GetActions() []Action {
	return actions
}

// GetAction returns one action with a specific ID.
func GetAction(id string) *Action {
	for _, a := range actions {
		if a.ID == id {
			return &a
		}
	}

	return nil
}

// SetActions sets the currently configured actions.
func SetActions(as []Action) {
	actions = as
}

// execAction executes an action and map its ins & outs.
func execAction(action Action, opts map[string]interface{}) bool {
	a := Action{
		Bee:  action.Bee,
		Name: action.Name,
	}

	for _, opt := range action.Options {
		ph := Placeholder{
			Name: opt.Name,
		}

		switch opt.Value.(type) {
		case string:
			var value bytes.Buffer

			tmpl, err := template.New(action.Bee + "_" + action.Name + "_" + opt.Name).Funcs(templatehelper.FuncMap).Parse(opt.Value.(string))
			if err == nil {
				err = tmpl.Execute(&value, opts)
			}
			if err != nil {
				panic(err)
			}

			ph.Type = "string"
			ph.Value = value.String()

		default:
			ph.Type = opt.Type
			ph.Value = opt.Value
		}
		a.Options = append(a.Options, ph)
	}

	bee := GetBee(a.Bee)
	if (*bee).IsRunning() {
		(*bee).LogAction()

		log.Debugln("\tExecuting action:", a.Bee, "/", a.Name, "-", GetActionDescriptor(&a).Description)
		for _, v := range a.Options {
			log.Debugln("\t\tOptions:", v)
		}

		(*bee).Action(a)
	} else {
		log.Debugln("\tNot executing action on stopped bee:", a.Bee, "/", a.Name, "-", GetActionDescriptor(&a).Description)
		for _, v := range a.Options {
			log.Debugln("\t\tOptions:", v)
		}
	}

	return true
}

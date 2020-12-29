/*
 *    Copyright (C) 2020 Anthony Corrado
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
 *      Anthony Corrado <anthony@synetz.fr>
 */

// Package jirabee is a Bee that can interface with Jira
package jirabee

import (
	"github.com/muesli/beehive/bees"
)

// JiraBeeFactory is a factory for JiraBees
type JiraBeeFactory struct {
	bees.BeeFactory
}

// New returns a new Bee instance configured with the supplied options.
func (factory *JiraBeeFactory) New(name, description string, options bees.BeeOptions) bees.BeeInterface {
	bee := JiraBee{
		Bee: bees.NewBee(name, factory.ID(), description, options),
	}
	bee.ReloadOptions(options)

	return &bee
}

// ID returns the ID of this Bee.
func (factory *JiraBeeFactory) ID() string {
	return "jirabee"
}

// Name returns the name of this Bee.
func (factory *JiraBeeFactory) Name() string {
	return "Jira"
}

// Description returns the desciption of this Bee.
func (factory *JiraBeeFactory) Description() string {
	return "Reacts to events on Jira"
}

// Image returns the filename of an image for this Bee.
func (factory *JiraBeeFactory) Image() string {
	return factory.ID() + ".png"
}

// LogoColor returns the preferred logo background color (used by the admin interface).
func (factory *JiraBeeFactory) LogoColor() string {
	return "#6098d0"
}

// Options returns the options available to configure this Bee.
func (factory *JiraBeeFactory) Options() []bees.BeeOptionDescriptor {
	opts := []bees.BeeOptionDescriptor{
		{
			Name:        "url",
			Description: "URL of the JIRA instance (for example, https://myjira.atlassian.com)",
			Type:        "string",
			Mandatory:   true,
		},
		{
			Name:        "username",
			Description: "Username used to access the JIRA API",
			Type:        "string",
			Mandatory:   true,
		},
		{
			Name:        "password",
			Description: "Password or API Token (for the cloud version) used to access the JIRA API",
			Type:        "string",
			Mandatory:   true,
		},
	}
	return opts
}

// Events describes the available events provided by this Bee.
func (factory *JiraBeeFactory) Events() []bees.EventDescriptor {
	events := []bees.EventDescriptor{}
	return events
}

// Actions describes the available actions provided by this Bee.
func (factory *JiraBeeFactory) Actions() []bees.ActionDescriptor {
	actions := []bees.ActionDescriptor{
		{
			Namespace:   factory.Name(),
			Name:        "create_issue",
			Description: "Create an issue",
			Options: []bees.PlaceholderDescriptor{
				{
					Name:        "project",
					Description: "Project name where to create the issue",
					Type:        "string",
					Mandatory:   true,
				},
				{
					Name:        "reporter_email",
					Description: "Reporter email address",
					Type:        "string",
					Mandatory:   false,
				},
				{
					Name:        "assignee_email",
					Description: "Assignee email address",
					Type:        "string",
					Mandatory:   false,
				},
				{
					Name:        "issue_type",
					Description: "Type of the issue (Story, Bug, ...)",
					Type:        "string",
					Mandatory:   true,
				},
				{
					Name:        "issue_summary",
					Description: "Summary of the issue (title)",
					Type:        "string",
					Mandatory:   true,
				},
				{
					Name:        "issue_description",
					Description: "Description of the issue",
					Type:        "string",
					Mandatory:   true,
				},
			},
		},
	}

	return actions
}

func init() {
	f := JiraBeeFactory{}
	bees.RegisterFactory(&f)
}

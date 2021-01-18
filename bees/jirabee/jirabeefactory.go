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
	return "#0052cc"
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
		{
			Name:        "address",
			Description: "Which address to listen on, eg: 0.0.0.0:12345",
			Type:        "address",
			Mandatory:   true,
		},
	}
	return opts
}

// Events describes the available events provided by this Bee.
func (factory *JiraBeeFactory) Events() []bees.EventDescriptor {
	events := []bees.EventDescriptor{
		{
			Namespace:   factory.Name(),
			Name:        "issue_created",
			Description: "Event triggered when an issue is created",
			Options: []bees.PlaceholderDescriptor{
				{
					Name:        "key",
					Description: "Key of the issue (example, BEEH-123)",
					Type:        "string",
				},
				{
					Name:        "title",
					Description: "Title of the issue",
					Type:        "string",
				},
				{
					Name:        "description",
					Description: "Description of the issue",
					Type:        "string",
				},
			},
		},
		{
			Namespace:   factory.Name(),
			Name:        "issue_status_updated",
			Description: "Event triggered when an issue status is updated",
			Options: []bees.PlaceholderDescriptor{
				{
					Name:        "key",
					Description: "Key of the issue (example, BEEH-123)",
					Type:        "string",
				},
				{
					Name:        "fromStatus",
					Description: "Previous status of the issue",
					Type:        "string",
				},
				{
					Name:        "toStatus",
					Description: "New status of the issue",
					Type:        "string",
				},
			},
		},
	}
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
		{
			Namespace:   factory.Name(),
			Name:        "update_issue_status",
			Description: "Update an issue status",
			Options: []bees.PlaceholderDescriptor{
				{
					Name:        "issue_key",
					Description: "Key of the issue (for example, BEEH-123)",
					Type:        "string",
					Mandatory:   true,
				},
				{
					Name:        "issue_new_status",
					Description: "New status of the issue (for example, Done)",
					Type:        "string",
					Mandatory:   true,
				},
			},
		},
		{
			Namespace:   factory.Name(),
			Name:        "comment_issue",
			Description: "Comment an issue",
			Options: []bees.PlaceholderDescriptor{
				{
					Name:        "issue_key",
					Description: "Key of the issue (for example, BEEH-123)",
					Type:        "string",
					Mandatory:   true,
				},
				{
					Name:        "comment_body",
					Description: "Body of the comment",
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

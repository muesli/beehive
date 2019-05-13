/*
 *    Copyright (C) 2017 Christian Muehlhaeuser
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
 *      Nicolas Martin <penguwingithub@gmail.com>
 *      Christian Muehlhaeuser <muesli@gmail.com>
 */

// Package githubbee is a Bee that can interface with GitHub
package githubbee

import (
	"github.com/muesli/beehive/bees"
)

// GitHubBeeFactory is a factory for GitHubBees
type GitHubBeeFactory struct {
	bees.BeeFactory
}

// New returns a new Bee instance configured with the supplied options.
func (factory *GitHubBeeFactory) New(name, description string, options bees.BeeOptions) bees.BeeInterface {
	bee := GitHubBee{
		Bee: bees.NewBee(name, factory.ID(), description, options),
	}
	bee.ReloadOptions(options)

	return &bee
}

// ID returns the ID of this Bee.
func (factory *GitHubBeeFactory) ID() string {
	return "githubbee"
}

// Name returns the name of this Bee.
func (factory *GitHubBeeFactory) Name() string {
	return "GitHub"
}

// Description returns the desciption of this Bee.
func (factory *GitHubBeeFactory) Description() string {
	return "Reacts to events on GitHub"
}

// Image returns the filename of an image for this Bee.
func (factory *GitHubBeeFactory) Image() string {
	return factory.ID() + ".png"
}

// LogoColor returns the preferred logo background color (used by the admin interface).
func (factory *GitHubBeeFactory) LogoColor() string {
	return "#6098d0"
}

// Options returns the options available to configure this Bee.
func (factory *GitHubBeeFactory) Options() []bees.BeeOptionDescriptor {
	opts := []bees.BeeOptionDescriptor{
		{
			Name:        "accesstoken",
			Description: "Your GitHub access token",
			Type:        "string",
			Mandatory:   true,
		},
		{
			Name:        "owner",
			Description: "Owner of the repository to watch",
			Type:        "string",
			Mandatory:   true,
		},
		{
			Name:        "repository",
			Description: "Name of the repository to watch",
			Type:        "string",
			Mandatory:   true,
		},
	}
	return opts
}

// Events describes the available events provided by this Bee.
func (factory *GitHubBeeFactory) Events() []bees.EventDescriptor {
	events := []bees.EventDescriptor{
		{
			Namespace:   factory.Name(),
			Name:        "push",
			Description: "Commits were pushed to a repository",
			Options: []bees.PlaceholderDescriptor{
				{
					Name:        "public",
					Description: "Indicates whether this was a public activity",
					Type:        "string",
				},
				{
					Name:        "repo",
					Description: "The repository that something was pushed to",
					Type:        "string",
				},
				{
					Name:        "repo_url",
					Description: "The repository's URL",
					Type:        "url",
				},
				{
					Name:        "username",
					Description: "Username that pushed the commits",
					Type:        "string",
				},
				{
					Name:        "url",
					Description: "URL to the diff-view of the changes",
					Type:        "url",
				},
				{
					Name:        "event_id",
					Description: "ID of the GitHub event",
					Type:        "string",
				},
			},
		},
		{
			Namespace:   factory.Name(),
			Name:        "commit",
			Description: "New commit in a repository",
			Options: []bees.PlaceholderDescriptor{
				{
					Name:        "public",
					Description: "Indicates whether this was a public activity",
					Type:        "string",
				},
				{
					Name:        "repo",
					Description: "The repository that something was committed to",
					Type:        "string",
				},
				{
					Name:        "repo_url",
					Description: "The repository's URL",
					Type:        "url",
				},
				{
					Name:        "username",
					Description: "Username that authored the commit",
					Type:        "string",
				},
				{
					Name:        "url",
					Description: "URL to the diff-view of the commit",
					Type:        "url",
				},
				{
					Name:        "id",
					Description: "ID of the commit",
					Type:        "string",
				},
				{
					Name:        "message",
					Description: "Commit message",
					Type:        "string",
				},
			},
		},
		{
			Namespace:   factory.Name(),
			Name:        "star",
			Description: "Someone starred a repository",
			Options: []bees.PlaceholderDescriptor{
				{
					Name:        "public",
					Description: "Indicates whether this was a public activity",
					Type:        "string",
				},
				{
					Name:        "repo",
					Description: "The repository that was starred",
					Type:        "string",
				},
				{
					Name:        "repo_url",
					Description: "The repository's URL",
					Type:        "url",
				},
				{
					Name:        "username",
					Description: "Username that starred the repository",
					Type:        "string",
				},
				{
					Name:        "event_id",
					Description: "ID of the GitHub event",
					Type:        "string",
				},
			},
		},
		{
			Namespace:   factory.Name(),
			Name:        "fork",
			Description: "Someone forked a repository",
			Options: []bees.PlaceholderDescriptor{
				{
					Name:        "public",
					Description: "Indicates whether this was a public activity",
					Type:        "string",
				},
				{
					Name:        "repo",
					Description: "The repository that was forked",
					Type:        "string",
				},
				{
					Name:        "repo_url",
					Description: "The repository's URL",
					Type:        "url",
				},
				{
					Name:        "username",
					Description: "Username that forked the repository",
					Type:        "string",
				},
				{
					Name:        "event_id",
					Description: "ID of the GitHub event",
					Type:        "string",
				},
			},
		},
		{
			Namespace:   factory.Name(),
			Name:        "issue_open",
			Description: "An issue was opened",
			Options: []bees.PlaceholderDescriptor{
				{
					Name:        "public",
					Description: "Indicates whether this was a public activity",
					Type:        "string",
				},
				{
					Name:        "repo",
					Description: "The repository that the issue belongs to",
					Type:        "string",
				},
				{
					Name:        "repo_url",
					Description: "The repository's URL",
					Type:        "url",
				},
				{
					Name:        "username",
					Description: "Username that created the issue",
					Type:        "string",
				},
				{
					Name:        "url",
					Description: "URL to the issue on GitHub",
					Type:        "url",
				},
				{
					Name:        "id",
					Description: "ID of the issue",
					Type:        "int",
				},
				{
					Name:        "title",
					Description: "Issue title",
					Type:        "string",
				},
				{
					Name:        "text",
					Description: "Issue text",
					Type:        "string",
				},
			},
		},
		{
			Namespace:   factory.Name(),
			Name:        "issue_close",
			Description: "An issue was closed",
			Options: []bees.PlaceholderDescriptor{
				{
					Name:        "public",
					Description: "Indicates whether this was a public activity",
					Type:        "string",
				},
				{
					Name:        "repo",
					Description: "The repository that the issue belongs to",
					Type:        "string",
				},
				{
					Name:        "repo_url",
					Description: "The repository's URL",
					Type:        "url",
				},
				{
					Name:        "username",
					Description: "Username that closed the issue",
					Type:        "string",
				},
				{
					Name:        "url",
					Description: "URL to the issue on GitHub",
					Type:        "url",
				},
				{
					Name:        "id",
					Description: "ID of the issue",
					Type:        "int",
				},
				{
					Name:        "title",
					Description: "Issue title",
					Type:        "string",
				},
				{
					Name:        "text",
					Description: "Issue text",
					Type:        "string",
				},
			},
		},
		{
			Namespace:   factory.Name(),
			Name:        "issue_comment",
			Description: "An issue was commented on",
			Options: []bees.PlaceholderDescriptor{
				{
					Name:        "public",
					Description: "Indicates whether this was a public activity",
					Type:        "string",
				},
				{
					Name:        "repo",
					Description: "The repository that the issue belongs to",
					Type:        "string",
				},
				{
					Name:        "repo_url",
					Description: "The repository's URL",
					Type:        "url",
				},
				{
					Name:        "username",
					Description: "Username that commented on the issue",
					Type:        "string",
				},
				{
					Name:        "url",
					Description: "URL to the comment on GitHub",
					Type:        "url",
				},
				{
					Name:        "id",
					Description: "ID of the issue",
					Type:        "int",
				},
				{
					Name:        "text",
					Description: "Issue text",
					Type:        "string",
				},
			},
		},
		{
			Namespace:   factory.Name(),
			Name:        "pullrequest_open",
			Description: "A Pull Request was created",
			Options: []bees.PlaceholderDescriptor{
				{
					Name:        "public",
					Description: "Indicates whether this was a public activity",
					Type:        "string",
				},
				{
					Name:        "repo",
					Description: "The repository that the Pull Request belongs to",
					Type:        "string",
				},
				{
					Name:        "repo_url",
					Description: "The repository's URL",
					Type:        "url",
				},
				{
					Name:        "username",
					Description: "Username that opened the Pull Request",
					Type:        "string",
				},
				{
					Name:        "url",
					Description: "URL to the Pull Request on GitHub",
					Type:        "url",
				},
				{
					Name:        "id",
					Description: "ID of the Pull Request",
					Type:        "int",
				},
				{
					Name:        "title",
					Description: "Pull Request title",
					Type:        "string",
				},
				{
					Name:        "text",
					Description: "Pull Request text",
					Type:        "string",
				},
			},
		},
		{
			Namespace:   factory.Name(),
			Name:        "pullrequest_close",
			Description: "A Pull Request was closed",
			Options: []bees.PlaceholderDescriptor{
				{
					Name:        "public",
					Description: "Indicates whether this was a public activity",
					Type:        "string",
				},
				{
					Name:        "repo",
					Description: "The repository that the Pull Request belongs to",
					Type:        "string",
				},
				{
					Name:        "repo_url",
					Description: "The repository's URL",
					Type:        "url",
				},
				{
					Name:        "username",
					Description: "Username that closed the Pull Request",
					Type:        "string",
				},
				{
					Name:        "url",
					Description: "URL to the Pull Request on GitHub",
					Type:        "url",
				},
				{
					Name:        "id",
					Description: "ID of the Pull Request",
					Type:        "int",
				},
				{
					Name:        "title",
					Description: "Pull Request title",
					Type:        "string",
				},
				{
					Name:        "text",
					Description: "Pull Request text",
					Type:        "string",
				},
			},
		},
		{
			Namespace:   factory.Name(),
			Name:        "pullrequest_review_comment",
			Description: "A Pull Request commit was commented on",
			Options: []bees.PlaceholderDescriptor{
				{
					Name:        "public",
					Description: "Indicates whether this was a public activity",
					Type:        "string",
				},
				{
					Name:        "repo",
					Description: "The repository that the Pull Request belongs to",
					Type:        "string",
				},
				{
					Name:        "repo_url",
					Description: "The repository's URL",
					Type:        "url",
				},
				{
					Name:        "username",
					Description: "Username that commented on the Pull Request commit",
					Type:        "string",
				},
				{
					Name:        "url",
					Description: "URL to the comment on GitHub",
					Type:        "url",
				},
				{
					Name:        "id",
					Description: "ID of the Pull Request",
					Type:        "int",
				},
				{
					Name:        "text",
					Description: "Review text",
					Type:        "string",
				},
			},
		},
	}
	return events
}

// Actions describes the available actions provided by this Bee.
func (factory *GitHubBeeFactory) Actions() []bees.ActionDescriptor {
	actions := []bees.ActionDescriptor{
		{
			Namespace:   factory.Name(),
			Name:        "follow",
			Description: "Follow a user",
			Options: []bees.PlaceholderDescriptor{
				{
					Name:        "username",
					Description: "Username to follow",
					Type:        "string",
					Mandatory:   true,
				},
			},
		},
		{
			Namespace:   factory.Name(),
			Name:        "unfollow",
			Description: "Unfollow a user",
			Options: []bees.PlaceholderDescriptor{
				{
					Name:        "username",
					Description: "Username to unfollow",
					Type:        "string",
					Mandatory:   true,
				},
			},
		},
		{
			Namespace:   factory.Name(),
			Name:        "star",
			Description: "Star a repository",
			Options: []bees.PlaceholderDescriptor{
				{
					Name:        "owner",
					Description: "Owner of the repository",
					Type:        "string",
					Mandatory:   true,
				},
				{
					Name:        "repo",
					Description: "Repository to star",
					Type:        "string",
					Mandatory:   true,
				},
			},
		},
		{
			Namespace:   factory.Name(),
			Name:        "unstar",
			Description: "Unstar a repository",
			Options: []bees.PlaceholderDescriptor{
				{
					Name:        "owner",
					Description: "Owner of the repository",
					Type:        "string",
					Mandatory:   true,
				},
				{
					Name:        "repo",
					Description: "Repository to unstar",
					Type:        "string",
					Mandatory:   true,
				},
			},
		},
	}

	return actions
}

func init() {
	f := GitHubBeeFactory{}
	bees.RegisterFactory(&f)
}

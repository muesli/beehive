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
 *      Christian Muehlhaeuser <muesli@gmail.com>
 */

// Package githubbee is a Bee that can interface with GitHub
package githubbee

import (
	"encoding/json"
	"fmt"

	"github.com/google/go-github/github"

	"github.com/muesli/beehive/bees"
)

func (mod *GitHubBee) handleReleaseEvent(event *github.Event) {
	var b github.ReleaseEvent
	json.Unmarshal(*event.RawPayload, &b)

	ev := bees.Event{
		Bee:  mod.Name(),
		Name: "release",
		Options: []bees.Placeholder{
			{
				Name:  "public",
				Type:  "bool",
				Value: *event.Public,
			},
			{
				Name:  "repo",
				Type:  "string",
				Value: *event.Repo.Name,
			},
			{
				Name:  "repo_url",
				Type:  "url",
				Value: "https://github.com/" + *event.Repo.Name,
			},
			{
				Name:  "username",
				Type:  "string",
				Value: *event.Actor.Login,
			},
			{
				Name:  "event_id",
				Type:  "string",
				Value: *event.ID,
			},
			{
				Name:  "release_title",
				Type:  "string",
				Value: *b.Release.Name,
			},
			{
				Name:  "release_tag_version",
				Type:  "string",
				Value: *b.Release.TagName,
			},
			{
				Name:  "release_description",
				Type:  "string",
				Value: *b.Release.Body,
			},
		},
	}
	mod.eventChan <- ev
}

func (mod *GitHubBee) handlePushEvent(event *github.Event) {
	var b github.PushEvent
	json.Unmarshal(*event.RawPayload, &b)

	ev := bees.Event{
		Bee:  mod.Name(),
		Name: "push",
		Options: []bees.Placeholder{
			{
				Name:  "public",
				Type:  "bool",
				Value: *event.Public,
			},
			{
				Name:  "repo",
				Type:  "string",
				Value: *event.Repo.Name,
			},
			{
				Name:  "repo_url",
				Type:  "url",
				Value: "https://github.com/" + *event.Repo.Name,
			},
			{
				Name:  "username",
				Type:  "string",
				Value: *event.Actor.Login,
			},
			{
				Name:  "url",
				Type:  "url",
				Value: fmt.Sprintf("https://github.com/%s/compare/%s...%s", *event.Repo.Name, (*b.Before)[0:12], (*b.Head)[0:12]),
			},
			{
				Name:  "event_id",
				Type:  "string",
				Value: *event.ID,
			},
		},
	}
	mod.eventChan <- ev

	for _, commit := range b.Commits {
		ev := bees.Event{
			Bee:  mod.Name(),
			Name: "commit",
			Options: []bees.Placeholder{
				{
					Name:  "public",
					Type:  "bool",
					Value: *event.Public,
				},
				{
					Name:  "repo",
					Type:  "string",
					Value: *event.Repo.Name,
				},
				{
					Name:  "repo_url",
					Type:  "url",
					Value: "https://github.com/" + *event.Repo.Name,
				},
				{
					Name:  "username",
					Type:  "string",
					Value: *commit.Author.Name,
				},
				{
					Name:  "url",
					Type:  "url",
					Value: fmt.Sprintf("https://github.com/%s/commit/%s", *event.Repo.Name, (*commit.SHA)[0:12]),
				},
				{
					Name:  "id",
					Type:  "string",
					Value: (*commit.SHA)[0:12],
				},
				{
					Name:  "message",
					Type:  "string",
					Value: *commit.Message,
				},
			},
		}
		mod.eventChan <- ev
	}
}

func (mod *GitHubBee) handleWatchEvent(event *github.Event) {
	var b github.WatchEvent
	json.Unmarshal(*event.RawPayload, &b)

	ev := bees.Event{
		Bee:  mod.Name(),
		Name: "star",
		Options: []bees.Placeholder{
			{
				Name:  "public",
				Type:  "bool",
				Value: *event.Public,
			},
			{
				Name:  "repo",
				Type:  "string",
				Value: *event.Repo.Name,
			},
			{
				Name:  "repo_url",
				Type:  "url",
				Value: "https://github.com/" + *event.Repo.Name,
			},
			{
				Name:  "username",
				Type:  "string",
				Value: *event.Actor.Login,
			},
			{
				Name:  "event_id",
				Type:  "string",
				Value: *event.ID,
			},
		},
	}
	mod.eventChan <- ev
}

func (mod *GitHubBee) handleForkEvent(event *github.Event) {
	var b github.ForkEvent
	json.Unmarshal(*event.RawPayload, &b)

	ev := bees.Event{
		Bee:  mod.Name(),
		Name: "fork",
		Options: []bees.Placeholder{
			{
				Name:  "public",
				Type:  "bool",
				Value: *event.Public,
			},
			{
				Name:  "repo",
				Type:  "string",
				Value: *event.Repo.Name,
			},
			{
				Name:  "repo_url",
				Type:  "url",
				Value: "https://github.com/" + *event.Repo.Name,
			},
			{
				Name:  "username",
				Type:  "string",
				Value: *event.Actor.Login,
			},
			{
				Name:  "event_id",
				Type:  "string",
				Value: *event.ID,
			},
		},
	}
	mod.eventChan <- ev
}

func (mod *GitHubBee) handleIssuesEvent(event *github.Event) {
	var b github.IssuesEvent
	json.Unmarshal(*event.RawPayload, &b)

	var t string

	switch *b.Action {
	case "closed":
		t = "close"
		fallthrough
	case "opened":
		if t == "" {
			t = "open"
		}
		ev := bees.Event{
			Bee:  mod.Name(),
			Name: "issue_" + t,
			Options: []bees.Placeholder{
				{
					Name:  "public",
					Type:  "bool",
					Value: *event.Public,
				},
				{
					Name:  "repo",
					Type:  "string",
					Value: *event.Repo.Name,
				},
				{
					Name:  "repo_url",
					Type:  "url",
					Value: "https://github.com/" + *event.Repo.Name,
				},
				{
					Name:  "username",
					Type:  "string",
					Value: *event.Actor.Login,
				},
				{
					Name:  "id",
					Type:  "int",
					Value: *b.Issue.Number,
				},
				{
					Name:  "url",
					Type:  "url",
					Value: *b.Issue.HTMLURL,
				},
				{
					Name:  "title",
					Type:  "string",
					Value: *b.Issue.Title,
				},
				{
					Name:  "text",
					Type:  "string",
					Value: *b.Issue.Body,
				},
			},
		}
		mod.eventChan <- ev

	default:
		mod.LogErrorf("Unhandled issues event: %s", *b.Action)
	}
}

func (mod *GitHubBee) handleIssueCommentEvent(event *github.Event) {
	var b github.IssueCommentEvent
	json.Unmarshal(*event.RawPayload, &b)

	switch *b.Action {
	case "created":
		ev := bees.Event{
			Bee:  mod.Name(),
			Name: "issue_comment",
			Options: []bees.Placeholder{
				{
					Name:  "public",
					Type:  "bool",
					Value: *event.Public,
				},
				{
					Name:  "repo",
					Type:  "string",
					Value: *event.Repo.Name,
				},
				{
					Name:  "repo_url",
					Type:  "url",
					Value: "https://github.com/" + *event.Repo.Name,
				},
				{
					Name:  "username",
					Type:  "string",
					Value: *event.Actor.Login,
				},
				{
					Name:  "id",
					Type:  "int",
					Value: *b.Issue.Number,
				},
				{
					Name:  "title",
					Type:  "string",
					Value: *b.Issue.Title,
				},
				{
					Name:  "url",
					Type:  "url",
					Value: *b.Comment.HTMLURL,
				},
				{
					Name:  "text",
					Type:  "string",
					Value: *b.Comment.Body,
				},
			},
		}
		mod.eventChan <- ev

	default:
		mod.LogErrorf("Unhandled issue comment event: %s", *b.Action)
	}
}

func (mod *GitHubBee) handlePullRequestEvent(event *github.Event) {
	var b github.PullRequestEvent
	json.Unmarshal(*event.RawPayload, &b)

	var t string

	switch *b.Action {
	case "closed":
		t = "close"
		fallthrough
	case "opened":
		if t == "" {
			t = "open"
		}
		ev := bees.Event{
			Bee:  mod.Name(),
			Name: "pullrequest_" + t,
			Options: []bees.Placeholder{
				{
					Name:  "public",
					Type:  "bool",
					Value: *event.Public,
				},
				{
					Name:  "repo",
					Type:  "string",
					Value: *event.Repo.Name,
				},
				{
					Name:  "repo_url",
					Type:  "url",
					Value: "https://github.com/" + *event.Repo.Name,
				},
				{
					Name:  "username",
					Type:  "string",
					Value: *event.Actor.Login,
				},
				{
					Name:  "id",
					Type:  "int",
					Value: *b.PullRequest.Number,
				},
				{
					Name:  "url",
					Type:  "url",
					Value: *b.PullRequest.HTMLURL,
				},
				{
					Name:  "title",
					Type:  "string",
					Value: *b.PullRequest.Title,
				},
				{
					Name:  "text",
					Type:  "string",
					Value: *b.PullRequest.Body,
				},
			},
		}

		mod.eventChan <- ev

	default:
		mod.LogErrorf("Unhandled pullrequest event: %s", *b.Action)
	}
}

func (mod *GitHubBee) handlePullRequestReviewCommentEvent(event *github.Event) {
	var b github.PullRequestReviewCommentEvent
	json.Unmarshal(*event.RawPayload, &b)

	switch *b.Action {
	case "created":
		ev := bees.Event{
			Bee:  mod.Name(),
			Name: "pullrequest_review_comment",
			Options: []bees.Placeholder{
				{
					Name:  "public",
					Type:  "bool",
					Value: *event.Public,
				},
				{
					Name:  "repo",
					Type:  "string",
					Value: *event.Repo.Name,
				},
				{
					Name:  "repo_url",
					Type:  "url",
					Value: "https://github.com/" + *event.Repo.Name,
				},
				{
					Name:  "username",
					Type:  "string",
					Value: *event.Actor.Login,
				},
				{
					Name:  "id",
					Type:  "int",
					Value: *b.PullRequest.Number,
				},
				{
					Name:  "title",
					Type:  "string",
					Value: *b.PullRequest.Title,
				},
				{
					Name:  "url",
					Type:  "url",
					Value: *b.Comment.HTMLURL,
				},
				{
					Name:  "text",
					Type:  "string",
					Value: *b.Comment.Body,
				},
			},
		}
		mod.eventChan <- ev

	default:
		mod.LogErrorf("Unhandled pullrequest review comment event: %s", *b.Action)
	}
}

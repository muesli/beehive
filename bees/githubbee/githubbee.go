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
	"time"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"

	"github.com/muesli/beehive/bees"
)

// GitHubBee is a Bee that can interface with GitHub
type GitHubBee struct {
	bees.Bee

	eventChan chan bees.Event
	client    *github.Client

	accessToken string
	owner       string
	repository  string
}

// Action triggers the actions passed to it.
func (mod *GitHubBee) Action(action bees.Action) []bees.Placeholder {
	outs := []bees.Placeholder{}
	switch action.Name {
	case "follow":
		var user string
		action.Options.Bind("username", &user)

		if _, err := mod.client.Users.Follow(user); err != nil {
			mod.LogErrorf("Failed to follow user: %v", err)
		}

	case "unfollow":
		var user string
		action.Options.Bind("username", &user)

		if _, err := mod.client.Users.Unfollow(user); err != nil {
			mod.LogErrorf("Failed to follow user: %v", err)
		}

	case "star":
		var user string
		var repo string
		action.Options.Bind("owner", &user)
		action.Options.Bind("repository", &repo)

		if _, err := mod.client.Activity.Star(user, repo); err != nil {
			mod.LogErrorf("Failed to star repository: %v", err)
		}

	case "unstar":
		var user string
		var repo string
		action.Options.Bind("owner", &user)
		action.Options.Bind("repository", &repo)

		if _, err := mod.client.Activity.Unstar(user, repo); err != nil {
			mod.LogErrorf("Failed to unstar repository: %v", err)
		}

	default:
		panic("Unknown action triggered in " + mod.Name() + ": " + action.Name)
	}

	return outs
}

// Run executes the Bee's event loop.
func (mod *GitHubBee) Run(eventChan chan bees.Event) {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: mod.accessToken},
	)
	tc := oauth2.NewClient(oauth2.NoContext, ts)

	mod.eventChan = eventChan
	mod.client = github.NewClient(tc)

	since := time.Now() // .Add(-time.Duration(24 * time.Hour))
	timeout := time.Duration(time.Second * 10)
	for {
		select {
		case <-mod.SigChan:
			return
		case <-time.After(timeout):
			mod.getRepositoryEvents(mod.owner, mod.repository, since)

		}
		since = time.Now()
		timeout = time.Duration(time.Minute)
	}
}

func (mod *GitHubBee) getRepositoryEvents(owner, repo string, since time.Time) {
	for page := 1; ; page++ {
		opts := &github.ListOptions{
			Page: page,
		}
		events, _, err := mod.client.Activity.ListRepositoryEvents(owner, repo, opts)
		if err != nil {
			mod.LogErrorf("Failed to fetch events: %v", err)
			return
		}
		if len(events) == 0 {
			mod.LogErrorf("No more events found")
			return
		}

		for _, v := range events {
			if since.After(*v.CreatedAt) {
				return
			}
			switch *v.Type {
			case "PushEvent":
				mod.handlePushEvent(v)
			case "WatchEvent":
				mod.handleWatchEvent(v)
			case "ForkEvent":
				mod.handleForkEvent(v)
			case "IssuesEvent":
				mod.handleIssuesEvent(v)
			case "IssueCommentEvent":
				mod.handleIssueCommentEvent(v)
			case "PullRequestEvent":
				mod.handlePullRequestEvent(v)
			case "PullRequestReviewCommentEvent":
				mod.handlePullRequestReviewCommentEvent(v)

			default:
				mod.LogErrorf("Unhandled event: %s", *v.Type)
			}
		}
	}
}

/*
func (mod *GitHubBee) getNotifications() {
	opts := &github.NotificationListOptions{
		All: true,
		// Participating: true,
		ListOptions: github.ListOptions{
			PerPage: 10,
		},
	}
	notif, _, err := mod.client.Activity.ListNotifications(opts)
	if err != nil {
		mod.LogErrorf("Failed to fetch notifications: %v", err)
	}
	for _, v := range notif {
		ev := bees.Event{
			Bee:  mod.Name(),
			Name: "notification",
			Options: []bees.Placeholder{
				{
					Name:  "subject_title",
					Type:  "string",
					Value: *v.Subject.Title,
				},
				{
					Name:  "subject_type",
					Type:  "string",
					Value: *v.Subject.Type,
				},
				{
					Name:  "subject_url",
					Type:  "url",
					Value: *v.Subject.URL,
				},
				{
					Name:  "reason",
					Type:  "string",
					Value: *v.Reason,
				},
				{
					Name:  "id",
					Type:  "string",
					Value: *v.ID,
				},
				{
					Name:  "url",
					Type:  "url",
					Value: *v.URL,
				},
			},
		}
		mod.eventChan <- ev
	}
}
*/

// ReloadOptions parses the config options and initializes the Bee.
func (mod *GitHubBee) ReloadOptions(options bees.BeeOptions) {
	mod.SetOptions(options)

	options.Bind("accesstoken", &mod.accessToken)
	options.Bind("owner", &mod.owner)
	options.Bind("repository", &mod.repository)
}

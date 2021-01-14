/*
 *    Copyright (C) 2021 Anthony Corrado
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
	"encoding/json"
	"fmt"

	"github.com/andygrunwald/go-jira"
	"github.com/muesli/beehive/bees"
)

// JiraEvent represents a Jira Event
type JiraEvent struct {
	WebhookEvent     string                 `json:"webhookEvent"`
	Issue            *jira.Issue            `json:"issue"`
	ChangelogHistory *jira.ChangelogHistory `json:"changelog"`
}

func (mod *JiraBee) handleJiraEvent(data []byte) (*JiraEvent, error) {
	jiraEvent := &JiraEvent{}
	err := json.Unmarshal(data, &jiraEvent)
	if err != nil {
		return nil, fmt.Errorf("Error during JiraEvent Unmarshal: %v", err)
	}

	switch jiraEvent.WebhookEvent {

	case "jira:issue_created":
		err = mod.handleIssueCreatedEvent(jiraEvent)
		if err != nil {
			return nil, err
		}
	case "jira:issue_updated":
		err = mod.handleIssueUpdatedEvent(jiraEvent)
		if err != nil {
			return nil, err
		}
	default:
		return jiraEvent, fmt.Errorf("Unhandled event: %s", jiraEvent.WebhookEvent)
	}

	return jiraEvent, nil
}

func (mod *JiraBee) handleIssueCreatedEvent(data *JiraEvent) error {

	if data.Issue == nil {
		return fmt.Errorf("Error occured during handleIssueCreatedEvent, Issue field was empty")
	}

	key := data.Issue.Key
	summary := data.Issue.Fields.Summary
	description := data.Issue.Fields.Description

	ev := bees.Event{
		Bee:  mod.Name(),
		Name: "issue_created",
		Options: []bees.Placeholder{
			{
				Name:  "key",
				Type:  "string",
				Value: key,
			},
			{
				Name:  "title",
				Type:  "string",
				Value: summary,
			},
			{
				Name:  "description",
				Type:  "string",
				Value: description,
			},
		},
	}

	mod.eventChan <- ev
	return nil
}

func (mod *JiraBee) handleIssueUpdatedEvent(jiraEvent *JiraEvent) error {

	if jiraEvent.Issue == nil {
		return fmt.Errorf("Issue field is empty, impossible to identify the issue key")
	}

	if jiraEvent.ChangelogHistory == nil {
		return fmt.Errorf("Changelog field is empty, impossible to identify the change type")
	}

	for _, item := range jiraEvent.ChangelogHistory.Items {
		if item.Field == "status" {
			mod.handleIssueStatusUpdatedEvent(jiraEvent.Issue.Key, item.FromString, item.ToString)
		}
	}

	return nil
}

func (mod *JiraBee) handleIssueStatusUpdatedEvent(key string, fromStatus string, toStatus string) error {

	ev := bees.Event{
		Bee:  mod.Name(),
		Name: "issue_status_updated",
		Options: []bees.Placeholder{
			{
				Name:  "key",
				Type:  "string",
				Value: key,
			},
			{
				Name:  "fromStatus",
				Type:  "string",
				Value: fromStatus,
			},
			{
				Name:  "toStatus",
				Type:  "string",
				Value: toStatus,
			},
		},
	}

	mod.eventChan <- ev
	return nil
}

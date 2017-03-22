package cricketbee

/*
 *    Copyright (C) 2017 Akash Shinde
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
 *      Akash Shinde <akashshinde159@gmail.com>
 */

import (
	"github.com/muesli/beehive/bees"
)

// Cricketbeefactory is a factory for CricketBee.
type CricketBeeFactory struct {
	bees.BeeFactory
}

// New returns a new Bee instance configured with the supplied options.
func (factory *CricketBeeFactory) New(name, description string, options bees.BeeOptions) bees.BeeInterface {
	bee := CricketBee{
		Bee: bees.NewBee(name, factory.ID(), description, options),
	}
	bee.ReloadOptions(options)
	return &bee
}

// ID returns the ID of this Bee.
func (factory *CricketBeeFactory) ID() string {
	return "cricketbee"
}

// Name returns the name of this Bee.
func (factory *CricketBeeFactory) Name() string {
	return "Cricket"
}

// Description returns the description of this Bee.
func (factory *CricketBeeFactory) Description() string {
	return "Triggers Events on cricket score updates"
}

// Image returns the filename of an image for this Bee.
func (factory *CricketBeeFactory) Image() string {
	return factory.ID() + ".png"
}

// LogoColor returns the preferred logo background color (used by the admin interface).
func (factory *CricketBeeFactory) LogoColor() string {
	return "#9571d6"
}

// Options returns the options available to configure this Bee.
func (factory *CricketBeeFactory) Options() []bees.BeeOptionDescriptor {
	opts := []bees.BeeOptionDescriptor{
		{
			Name:        "favourite_team",
			Default:     "Ind",
			Description: "Please mention for which team you want to get updates",
			Type:        "string",
			Mandatory:   true,
		},
	}
	return opts
}

// Events describes the available events provided by this Bee.
func (factory *CricketBeeFactory) Events() []bees.EventDescriptor {
	events := []bees.EventDescriptor{
		{
			Namespace:   factory.Name(),
			Name:        "out",
			Description: "Event gets triggered on fall of wickets",
			Options: []bees.PlaceholderDescriptor{
				{
					Name:        "batting_team",
					Description: "Name of currently batting team",
					Type:        "string",
				},
				{
					Name:        "score",
					Description: "Runs",
					Type:        "string",
				},
				{
					Name:        "wickets",
					Description: "Wickets fallen",
					Type:        "string",
				},
				{
					Name:        "overs",
					Description: "Total number of overs",
					Type:        "string",
				},
			},
		},
		{
			Namespace:   factory.Name(),
			Name:        "over_changed",
			Description: "Event will get triggered when overs get completed",
			Options: []bees.PlaceholderDescriptor{
				{
					Name:        "batting_team",
					Description: "Name of currently batting team",
					Type:        "string",
				},
				{
					Name:        "score",
					Description: "Runs",
					Type:        "string",
				},
				{
					Name:        "wickets",
					Description: "Wickets fallen",
					Type:        "string",
				},
				{
					Name:        "overs",
					Description: "Total number of overs",
					Type:        "string",
				},
			},
		},
		{
			Namespace:   factory.Name(),
			Name:        "run_change",
			Description: "Event gets triggered when there is any change in score",
			Options: []bees.PlaceholderDescriptor{
				{
					Name:        "batting_team",
					Description: "Name of currently batting team",
					Type:        "string",
				},
				{
					Name:        "score",
					Description: "Runs",
					Type:        "string",
				},
				{
					Name:        "wickets",
					Description: "Wickets fallen",
					Type:        "string",
				},
				{
					Name:        "overs",
					Description: "Total number of overs",
					Type:        "string",
				},
			},
		},
	}
	return events
}

// Actions describes the available actions provided by this Bee.
func (factory *CricketBeeFactory) Actions() []bees.ActionDescriptor {
	actions := []bees.ActionDescriptor{}
	return actions
}

func init() {
	f := CricketBeeFactory{}
	bees.RegisterFactory(&f)
}

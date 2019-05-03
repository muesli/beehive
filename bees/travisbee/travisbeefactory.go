/*
 *    Copyright (C) 2019 CalmBit
 *                  2014-2019 Christian Muehlhaeuser
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
 *      CalmBit <calmbit@posteo.net>
 *      Christian Muehlhaeuser <muesli@gmail.com>
 */

package travisbee

import (
	"github.com/muesli/beehive/bees"
)

// TravisBeeFactory is a factory for TravisBees.
type TravisBeeFactory struct {
	bees.BeeFactory

	lastBuilds map[string]string
}

// New returns a new Bee instance configured with the supplied options.
func (factory *TravisBeeFactory) New(name, description string, options bees.BeeOptions) bees.BeeInterface {
	bee := TravisBee{
		Bee: bees.NewBee(name, factory.ID(), description, options),
	}
	bee.ReloadOptions(options)

	return &bee
}

// ID returns the ID of this Bee.
func (factory *TravisBeeFactory) ID() string {
	return "travisbee"
}

// Name returns the name of this Bee.
func (factory *TravisBeeFactory) Name() string {
	return "Travis CI"
}

// Description returns the description of this Bee.
func (factory *TravisBeeFactory) Description() string {
	return "Allows for monitoring Travis CI jobs"
}

// Image returns the filename of an image for this Bee.
func (factory *TravisBeeFactory) Image() string {
	return factory.ID() + ".png"
}

// LogoColor returns the preferred logo background color (used by the admin interface).
func (factory *TravisBeeFactory) LogoColor() string {
	return "#3EAAAF"
}

// Options returns the options available to configure this Bee.
func (factory *TravisBeeFactory) Options() []bees.BeeOptionDescriptor {
	opts := []bees.BeeOptionDescriptor{
		{
			Name:        "api_key",
			Description: "Your travis-ci.org API key",
			Type:        "string",
			Mandatory:   true,
		},
	}
	return opts
}

// Events describes the available events provided by this Bee.
func (factory *TravisBeeFactory) Events() []bees.EventDescriptor {
	events := []bees.EventDescriptor{
		{
			Namespace:   factory.Name(),
			Name:        "build_started",
			Description: "is triggered when a build has been started",
			Options: []bees.PlaceholderDescriptor{
				{
					Name:        "id",
					Description: "The id of the build",
					Type:        "uint",
				},
				{
					Name:        "state",
					Description: "The current state of the build",
					Type:        "string",
				},
				{
					Name:        "repo_slug",
					Description: "The slug of the repo being built",
					Type:        "string",
				},
				{
					Name:        "duration",
					Description: "The current duration of the build in seconds",
					Type:        "uint",
				},
			},
		},
		{
			Namespace:   factory.Name(),
			Name:        "build_status_change",
			Description: "is triggered when a currently active build changes status",
			Options: []bees.PlaceholderDescriptor{
				{
					Name:        "id",
					Description: "The id of the build",
					Type:        "uint",
				},
				{
					Name:        "state",
					Description: "The current state of the build",
					Type:        "string",
				},
				{
					Name:        "last_state",
					Description: "The last state of the build",
					Type:        "string",
				},
				{
					Name:        "repo_slug",
					Description: "The slug of the repo being built",
					Type:        "string",
				},
				{
					Name:        "duration",
					Description: "The current duration of the build in seconds",
					Type:        "uint",
				},
			},
		},
		{
			Namespace:   factory.Name(),
			Name:        "build_finished",
			Description: "is triggered when a build enters any terminal state (canceled/errored/passed/failed)",
			Options: []bees.PlaceholderDescriptor{
				{
					Name:        "id",
					Description: "The id of the build",
					Type:        "uint",
				},
				{
					Name:        "state",
					Description: "The final state of the build",
					Type:        "string",
				},
				{
					Name:        "repo_slug",
					Description: "The slug of the repo being built",
					Type:        "string",
				},
				{
					Name:        "duration",
					Description: "The duration of the build in seconds",
					Type:        "uint",
				},
			},
		},
	}
	return events
}

func init() {
	f := TravisBeeFactory{}
	bees.RegisterFactory(&f)
}

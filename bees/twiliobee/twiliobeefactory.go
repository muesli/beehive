/*
 *    Copyright (C) 2017 James Vaughan
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
 *      James Vaughan <james@jamesbvaughan.com>
 */

package twiliobee

import (
	"github.com/muesli/beehive/bees"
)

// TwilioBeeFactory is a factory for TwilioBees.
type TwilioBeeFactory struct {
	bees.BeeFactory
}

// New returns a new Bee instance configured with the supplied options.
func (factory *TwilioBeeFactory) New(name, description string, options bees.BeeOptions) bees.BeeInterface {
	bee := TwilioBee{
		Bee: bees.NewBee(name, factory.ID(), description, options),
	}
	bee.ReloadOptions(options)

	return &bee
}

// ID returns the ID of this Bee.
func (factory *TwilioBeeFactory) ID() string {
	return "twiliobee"
}

// Name returns the name of this Bee.
func (factory *TwilioBeeFactory) Name() string {
	return "Twilio"
}

// Description returns the description of this Bee.
func (factory *TwilioBeeFactory) Description() string {
	return "Sends SMS messages"
}

// Image returns the filename of an image for this Bee.
func (factory *TwilioBeeFactory) Image() string {
	return factory.ID() + ".png"
}

// LogoColor returns the preferred logo background color (used by the admin interface).
func (factory *TwilioBeeFactory) LogoColor() string {
	return "#0d122b"
}

// Options returns the options available to configure this Bee.
func (factory *TwilioBeeFactory) Options() []bees.BeeOptionDescriptor {
	opts := []bees.BeeOptionDescriptor{
		{
			Name:        "account_sid",
			Description: "Twilio account SID",
			Type:        "string",
			Mandatory:   true,
		},
		{
			Name:        "auth_token",
			Description: "Twilio auth token",
			Type:        "string",
			Mandatory:   true,
		},
		{
			Name:        "from_number",
      Description: "Phone number to send SMS messages from (ex: \"+15558675309\")",
			Type:        "string",
			Mandatory:   true,
		},
		{
			Name:        "to_number",
			Description: "Phone number to send SMS messages to (ex: \"+15554815162\")",
			Type:        "string",
			Mandatory:   true,
		},
	}
	return opts
}

// Actions describes the available actions provided by this Bee.
func (factory *TwilioBeeFactory) Actions() []bees.ActionDescriptor {
	actions := []bees.ActionDescriptor{
		{
			Namespace:   factory.Name(),
			Name:        "send",
			Description: "Sends an SMS message",
			Options: []bees.PlaceholderDescriptor{
				{
					Name:        "body",
					Description: "Message body",
					Type:        "string",
					Mandatory:   true,
				},
			},
		},
	}
	return actions
}

func init() {
	f := TwilioBeeFactory{}
	bees.RegisterFactory(&f)
}

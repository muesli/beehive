/*
 *    Copyright (C) 2016 Gonzalo Izquierdo
 *                  2017 Christian Muehlhaeuser
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
 *      Gonzalo Izquierdo <lalotone@gmail.com>
 *      Christian Muehlhaeuser <muesli@gmail.com>
 */

package telegrambee

import "github.com/muesli/beehive/bees"

// TelegramBeeFactory is a factory for TelegramBees.
type TelegramBeeFactory struct {
	bees.BeeFactory
}

// New returns a new Bee instance configured with the supplied options.
func (factory *TelegramBeeFactory) New(name, description string, options bees.BeeOptions) bees.BeeInterface {
	bee := TelegramBee{
		Bee: bees.NewBee(name, factory.ID(), description, options),
	}
	bee.ReloadOptions(options)

	return &bee
}

// ID returns the ID of this Bee.
func (factory *TelegramBeeFactory) ID() string {
	return "telegrambee"
}

// Name returns the name of this Bee.
func (factory *TelegramBeeFactory) Name() string {
	return "Telegram"
}

// Description returns the description of this Bee.
func (factory *TelegramBeeFactory) Description() string {
	return "Connects to Telegram"
}

// Image returns the filename of an image for this Bee.
func (factory *TelegramBeeFactory) Image() string {
	return factory.ID() + ".png"
}

// LogoColor returns the preferred logo background color (used by the admin interface).
func (factory *TelegramBeeFactory) LogoColor() string {
	return "#003b66"
}

// Options returns the options available to configure this Bee.
func (factory *TelegramBeeFactory) Options() []bees.BeeOptionDescriptor {
	opts := []bees.BeeOptionDescriptor{
		{
			Name:        "api_key",
			Description: "Telegram bot API key",
			Type:        "string",
			Mandatory:   true,
		},
		{
			Name:        "formatting_enabled",
			Description: "Enable HTML text formatting",
			Type:        "bool",
			Default:     false,
			Mandatory:   false,
		},
	}
	return opts
}

// Events describes the available events provided by this Bee.
func (factory *TelegramBeeFactory) Events() []bees.EventDescriptor {
	events := []bees.EventDescriptor{
		{
			Namespace:   factory.Name(),
			Name:        "message",
			Description: "A message received via Telegram bot",
			Options: []bees.PlaceholderDescriptor{
				{
					Name:        "text",
					Description: "The message that was received",
					Type:        "string",
				}, {
					Name:        "chat_id",
					Description: "Telegram's chat ID",
					Type:        "string",
				},
				{
					Name:        "user_id",
					Description: "User ID sending the message",
					Type:        "string",
				},
				{
					Name:        "timestamp",
					Description: "Timestamp",
					Type:        "timestamp",
				},
			},
		},
	}

	return events
}

// Actions describes the available actions provided by this Bee.
func (factory *TelegramBeeFactory) Actions() []bees.ActionDescriptor {
	actions := []bees.ActionDescriptor{{
		Namespace:   factory.Name(),
		Name:        "send",
		Description: "Sends a message to a Telegram chat or group",
		Options: []bees.PlaceholderDescriptor{
			{
				Name:        "chat_id",
				Description: "Telegram chat/group to send the message to",
				Type:        "string",
				Mandatory:   true,
			},
			{
				Name:        "text",
				Description: "Content of the message",
				Type:        "string",
				Mandatory:   true,
			},
		},
	}}
	return actions
}

func init() {
	f := TelegramBeeFactory{}
	bees.RegisterFactory(&f)
}

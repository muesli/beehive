/*
 *    Copyright (C) 2016 Gonzalo Izquierdo
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
 */

package transmissionbee

import (
	"github.com/kr/pretty"
	"github.com/muesli/beehive/bees"
	"github.com/odwrtw/transmission"
)

type TransmissionBeeFactory struct {
	bees.BeeFactory
}

func (factory *TransmissionBeeFactory) New(name, description string, options bees.BeeOptions) bees.BeeInterface {
	serverURL := options.GetValue("serverURL").(string)
	username := options.GetValue("username").(string)
	password := options.GetValue("password").(string)

	conf := transmission.Config{
		Address:  serverURL,
		User:     username,
		Password: password,
	}

	t, err := transmission.New(conf)
	if err != nil {
		pretty.Println(err)
	}
	bee := TransmissionBee{
		Bee:    bees.NewBee(name, factory.Name(), description),
		client: t,
	}

	return &bee
}

func (factory *TransmissionBeeFactory) Options() []bees.BeeOptionDescriptor {
	opts := []bees.BeeOptionDescriptor{
		bees.BeeOptionDescriptor{
			Name:        "serverURL",
			Description: "Transmission server URL",
			Type:        "url",
			Mandatory:   true,
		},
		bees.BeeOptionDescriptor{
			Name:        "username",
			Description: "Transmission server username",
			Type:        "string",
			Mandatory:   true,
		},
		bees.BeeOptionDescriptor{
			Name:        "password",
			Description: "Transmission server password",
			Type:        "string",
			Mandatory:   true,
		},
	}
	return opts
}

func (factory *TransmissionBeeFactory) Actions() []bees.ActionDescriptor {
	actions := []bees.ActionDescriptor{bees.ActionDescriptor{
		Namespace:   factory.Name(),
		Name:        "add-torrent",
		Description: "Torrent URL or magnet",
		Options: []bees.PlaceholderDescriptor{
			bees.PlaceholderDescriptor{
				Name:        "torrent",
				Description: "Telegram chat/group to send the message to",
				Type:        "string",
			},
			bees.PlaceholderDescriptor{
				Name:        "commandPrefix",
				Description: "String that precedes the torrent URL/magnet (will be removed)",
				Type:        "string",
			},
		},
	}}
	return actions
}

func (factory *TransmissionBeeFactory) Name() string {
	return "transmissionbee"
}

func (factory *TransmissionBeeFactory) Image() string {
	return factory.Name() + ".png"
}

func (factory *TransmissionBeeFactory) Description() string {
	return "A bee for adding torrents to a transmission server"
}

func init() {
	f := TransmissionBeeFactory{}
	bees.RegisterFactory(&f)
}

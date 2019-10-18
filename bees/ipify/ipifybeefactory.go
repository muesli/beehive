/*
 *    Copyright (C) 2019 Sergio Rubio
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
 *      Sergio Rubio <sergio@rubio.im>
 */

package ipify

import (
	"github.com/muesli/beehive/bees"
)

const defaultUpdateInterval int = 60

// IpifyBeeFactory takes care of initializing IpifyBee
type IpifyBeeFactory struct {
	bees.BeeFactory
}

// New returns a new Bee instance configured with the supplied options.
func (factory *IpifyBeeFactory) New(name, description string, options bees.BeeOptions) bees.BeeInterface {
	bee := IpifyBee{
		Bee: bees.NewBee(name, factory.ID(), description, options),
	}
	bee.ReloadOptions(options)

	return &bee
}

// ID returns the ID of this Bee.
func (factory *IpifyBeeFactory) ID() string {
	return "ipify"
}

// Name returns the name of this Bee.
func (factory *IpifyBeeFactory) Name() string {
	return "ipify"
}

// Description returns the description of this Bee.
func (factory *IpifyBeeFactory) Description() string {
	return "Monitor your public IP address via ipify.org and notify when the IP changes"
}

// Image returns the asset name of this Bee (in the assets/bees folder)
func (factory *IpifyBeeFactory) Image() string {
	return factory.Name() + ".png"
}

// Events describes the available events provided by this Bee.
func (factory *IpifyBeeFactory) Events() []bees.EventDescriptor {
	events := []bees.EventDescriptor{
		{
			Namespace:   factory.Name(),
			Name:        "ip",
			Description: "The public IP retrieved from ipify.org",
			Options: []bees.PlaceholderDescriptor{
				{
					Name:        "ip",
					Description: "IP address string",
					Type:        "string",
				},
			},
		},
	}
	return events
}

// Options returns the options available to configure this Bee.
func (factory *IpifyBeeFactory) Options() []bees.BeeOptionDescriptor {
	opts := []bees.BeeOptionDescriptor{
		{
			Name:        "interval",
			Description: "Interval in minutes to query ipify.org (60 minutes by default)",
			Type:        "int",
			Default:     defaultUpdateInterval,
			Mandatory:   false,
		},
	}
	return opts
}

func init() {
	f := IpifyBeeFactory{}
	bees.RegisterFactory(&f)
}

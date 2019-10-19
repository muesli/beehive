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

package cfddns

import (
	"github.com/muesli/beehive/bees"
)

// CFDDNSBeeFactory takes care of initializing CFDDNSBee
type CFDDNSBeeFactory struct {
	bees.BeeFactory
}

// New returns a new Bee instance configured with the supplied options.
func (factory *CFDDNSBeeFactory) New(name, description string, options bees.BeeOptions) bees.BeeInterface {
	bee := CFDDNSBee{
		Bee: bees.NewBee(name, factory.ID(), description, options),
	}
	bee.ReloadOptions(options)

	return &bee
}

// ID returns the ID of this Bee.
func (factory *CFDDNSBeeFactory) ID() string {
	return "cfddns"
}

// Name returns the name of this Bee.
func (factory *CFDDNSBeeFactory) Name() string {
	return "cfddns"
}

// Description returns the description of this Bee.
func (factory *CFDDNSBeeFactory) Description() string {
	return "Updates the IP address of a Cloudflare domain name"
}

// Image returns the asset name of this Bee (in the assets/bees folder)
func (factory *CFDDNSBeeFactory) Image() string {
	return factory.Name() + ".png"
}

// Events describes the available events provided by this Bee.
func (factory *CFDDNSBeeFactory) Events() []bees.EventDescriptor {
	return []bees.EventDescriptor{}
}

// Actions describes the available actions provided by this Bee.
func (factory *CFDDNSBeeFactory) Actions() []bees.ActionDescriptor {
	actions := []bees.ActionDescriptor{
		{
			Namespace:   factory.Name(),
			Name:        "update_domain",
			Description: "Updates the DNS name",
			Options: []bees.PlaceholderDescriptor{
				{
					Name:        "address",
					Description: "IP addresss",
					Type:        "string",
					Mandatory:   true,
				},
			},
		},
	}
	return actions
}

// Options returns the options available to configure this Bee.
func (factory *CFDDNSBeeFactory) Options() []bees.BeeOptionDescriptor {
	opts := []bees.BeeOptionDescriptor{
		{
			Name:        "email",
			Description: "CF API user email",
			Type:        "string",
			Mandatory:   true,
		},
		{
			Name:        "key",
			Description: "CF API Key",
			Type:        "string",
			Mandatory:   true,
		},
		{
			Name:        "domain",
			Description: "CF domain to update",
			Type:        "string",
			Mandatory:   true,
		},
	}
	return opts
}

func init() {
	f := CFDDNSBeeFactory{}
	bees.RegisterFactory(&f)
}

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
	"regexp"

	"github.com/cloudflare/cloudflare-go"
	"github.com/muesli/beehive/bees"
)

var domainSplitter = regexp.MustCompile(".+\\.(.+\\..+)")

// CFDDNSBee updates a Cloudflare domain name
type CFDDNSBee struct {
	bees.Bee
	client *cloudflare.API
	domain string
}

// Run executes the Bee's event loop.
func (mod *CFDDNSBee) Run(eventChan chan bees.Event) {
}

// Action triggers the action passed to it.
func (mod *CFDDNSBee) Action(action bees.Action) []bees.Placeholder {
	outs := []bees.Placeholder{}

	switch action.Name {
	case "update_domain":
		address := action.Options.Value("address").(string)
		domain := mod.domain
		zoneName := domainSplitter.FindStringSubmatch(domain)[1]
		// Resolve the zone and record id for the host
		zone, err := mod.client.ZoneIDByName(zoneName)
		if err != nil {
			mod.LogErrorf("zone id resolution failed: %v", err)
			return outs
		}
		recs, err := mod.client.DNSRecords(zone, cloudflare.DNSRecord{Name: domain, Type: "A"})
		if err != nil {
			mod.LogErrorf("record id resolution failed: %v", err)
			return outs
		}
		if len(recs) != 1 {
			mod.LogErrorf("invalid number of DNS records found: %+v", recs)
			return outs
		}
		record := recs[0]

		// Post the Cloudflare dns update
		record.Content = address

		if err := mod.client.UpdateDNSRecord(zone, record.ID, record); err != nil {
			mod.LogErrorf("dns record update for %s failed: %v", domain, err)
		}
	default:
		mod.LogDebugf("Unknown action triggered in %s: %s", mod.Name(), action.Name)
	}

	return outs
}

// ReloadOptions parses the config options and initializes the Bee.
func (mod *CFDDNSBee) ReloadOptions(options bees.BeeOptions) {
	mod.SetOptions(options)
	key := options.Value("key").(string)
	user := options.Value("email").(string)
	mod.domain = options.Value("domain").(string)
	c, err := cloudflare.New(key, user)
	if err != nil {
		mod.LogFatal(err)
	}
	mod.client = c
}

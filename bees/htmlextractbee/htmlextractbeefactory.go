/*
 *    Copyright (C) 2014-2017 Christian Muehlhaeuser
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
 *      Christian Muehlhaeuser <muesli@gmail.com>
 */

package htmlextractbee

import (
	"github.com/muesli/beehive/bees"
)

type HtmlExtractBeeFactory struct {
	bees.BeeFactory
}

// Interface impl

func (factory *HtmlExtractBeeFactory) New(name, description string, options bees.BeeOptions) bees.BeeInterface {
	bee := HtmlExtractBee{
		Bee: bees.NewBee(name, factory.Name(), description, options),
	}
	bee.ReloadOptions(options)

	return &bee
}

func (factory *HtmlExtractBeeFactory) Name() string {
	return "htmlextractbee"
}

func (factory *HtmlExtractBeeFactory) Description() string {
	return "A bee that extracts information from an arbitrary web page"
}

func (factory *HtmlExtractBeeFactory) Image() string {
	return factory.Name() + ".png"
}

func (factory *HtmlExtractBeeFactory) Actions() []bees.ActionDescriptor {
	actions := []bees.ActionDescriptor{
		{
			Namespace:   factory.Name(),
			Name:        "extract",
			Description: "Extract information from a web page",
			Options: []bees.PlaceholderDescriptor{
				{
					Name:        "url",
					Description: "The web page you want to extract information from",
					Type:        "url",
				},
			},
		},
	}
	return actions
}

func (factory *HtmlExtractBeeFactory) Events() []bees.EventDescriptor {
	events := []bees.EventDescriptor{
		{
			Namespace:   factory.Name(),
			Name:        "info_extracted",
			Description: "Information has been extracted from the web page",
			Options: []bees.PlaceholderDescriptor{
				{
					Name:        "title",
					Description: "HTML title tag",
					Type:        "string",
				},
				{
					Name:        "domain",
					Description: "Domain",
					Type:        "string",
				},
				{
					Name:        "topimage",
					Description: "The top image for the page",
					Type:        "url",
				},
				{
					Name:        "finalurl",
					Description: "Eventual URL after potentially being redirected",
					Type:        "url",
				},
				{
					Name:        "meta_description",
					Description: "HTML meta description",
					Type:        "string",
				},
				{
					Name:        "meta_keywords",
					Description: "HTML meta keywords",
					Type:        "string",
				},
			},
		},
	}
	return events
}

func init() {
	f := HtmlExtractBeeFactory{}
	bees.RegisterFactory(&f)
}

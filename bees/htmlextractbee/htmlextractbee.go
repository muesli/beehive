/*
 *    Copyright (C) 2014 Christian Muehlhaeuser
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

// beehive's Html Extraction module.
package htmlextractbee

import (
	"github.com/muesli/beehive/bees"
	"github.com/advancedlogic/GoOse"
	"strings"
)

type HtmlExtractBee struct {
	bees.Bee

	url string

	evchan chan bees.Event
}

func (mod *HtmlExtractBee) Action(action bees.Action) []bees.Placeholder {
	outs := []bees.Placeholder{}

	switch action.Name {
	case "extract":
		var url string
		for _, opt := range action.Options {
			if opt.Name == "url" {
				url = opt.Value.(string)
				if start := strings.Index(url, "http"); start >= 0 {
					url = url[start:]
					if end := strings.Index(url, " "); end >= 0 {
						url = url[:end]
					}
				}
			}
		}

		g := goose.New()
    	article := g.ExtractFromUrl(url)
    	if len(strings.TrimSpace(article.Title)) > 0 {
	    	ev := bees.Event{
				Bee:  mod.Name(),
				Name: "info_extracted",
				Options: []bees.Placeholder{
					bees.Placeholder{
						Name:  "title",
						Type:  "string",
						Value: article.Title,
					},
					bees.Placeholder{
						Name:  "domain",
						Type:  "string",
						Value: article.Domain,
					},
					bees.Placeholder{
						Name:  "topimage",
						Type:  "url",
						Value: article.TopImage,
					},
					bees.Placeholder{
						Name:  "finalurl",
						Type:  "url",
						Value: article.FinalUrl,
					},
					bees.Placeholder{
						Name:  "meta_description",
						Type:  "string",
						Value: article.MetaDescription,
					},
					bees.Placeholder{
						Name:  "meta_keywords",
						Type:  "string",
						Value: article.MetaKeywords,
					},
				},
			}
			mod.evchan <- ev
		}

	default:
		panic("Unknown action triggered in " +mod.Name()+": "+action.Name)
	}

	return outs
}

func (mod *HtmlExtractBee) Run(eventChan chan bees.Event) {
	mod.evchan = eventChan
}

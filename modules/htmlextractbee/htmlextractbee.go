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
	"github.com/muesli/beehive/modules"
	"github.com/advancedlogic/GoOse"
	"strings"
)

type HtmlExtractBee struct {
	modules.Module

	url string

	evchan chan modules.Event
}

func (mod *HtmlExtractBee) Action(action modules.Action) []modules.Placeholder {
	outs := []modules.Placeholder{}

	switch action.Name {
	case "extract":
		var url string
		for _, opt := range action.Options {
			if opt.Name == "url" {
				url = opt.Value.(string)
				if start := strings.Index(url, "http://"); start >= 0 {
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
	    	ev := modules.Event{
				Bee:  mod.Name(),
				Name: "info_extracted",
				Options: []modules.Placeholder{
					modules.Placeholder{
						Name:  "title",
						Type:  "string",
						Value: article.Title,
					},
					modules.Placeholder{
						Name:  "domain",
						Type:  "string",
						Value: article.Domain,
					},
					modules.Placeholder{
						Name:  "topimage",
						Type:  "url",
						Value: article.TopImage,
					},
					modules.Placeholder{
						Name:  "finalurl",
						Type:  "url",
						Value: article.FinalUrl,
					},
					modules.Placeholder{
						Name:  "meta_description",
						Type:  "string",
						Value: article.MetaDescription,
					},
					modules.Placeholder{
						Name:  "meta_keywords",
						Type:  "string",
						Value: article.MetaKeywords,
					},
				},
			}
			mod.evchan <- ev
		}

	default:
		// unknown action
		return outs
	}

	return outs
}

func (mod *HtmlExtractBee) Run(eventChan chan modules.Event) {
	mod.evchan = eventChan
}

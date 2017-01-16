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
	"strings"

	"github.com/advancedlogic/GoOse"
	"github.com/muesli/beehive/bees"
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
		action.Options.Bind("url", &url)
		if start := strings.Index(url, "http"); start >= 0 {
			url = url[start:]
			if end := strings.Index(url, " "); end >= 0 {
				url = url[:end]
			}
		}

		g := goose.New()
		article, _ := g.ExtractFromURL(url)
		article.Title = strings.TrimSpace(strings.Replace(article.Title, "\n", " ", -1))
		if strings.HasPrefix(article.TopImage, "http://data:image") {
			article.TopImage = ""
		}
		if len(article.Title) > 0 {
			ev := bees.Event{
				Bee:  mod.Name(),
				Name: "info_extracted",
				Options: []bees.Placeholder{
					{
						Name:  "title",
						Type:  "string",
						Value: article.Title,
					},
					{
						Name:  "domain",
						Type:  "string",
						Value: article.Domain,
					},
					{
						Name:  "topimage",
						Type:  "url",
						Value: article.TopImage,
					},
					{
						Name:  "finalurl",
						Type:  "url",
						Value: article.FinalURL,
					},
					{
						Name:  "meta_description",
						Type:  "string",
						Value: article.MetaDescription,
					},
					{
						Name:  "meta_keywords",
						Type:  "string",
						Value: article.MetaKeywords,
					},
				},
			}
			mod.evchan <- ev
		}

	default:
		panic("Unknown action triggered in " + mod.Name() + ": " + action.Name)
	}

	return outs
}

func (mod *HtmlExtractBee) Run(eventChan chan bees.Event) {
	mod.evchan = eventChan
}

func (mod *HtmlExtractBee) SetOptions(options bees.BeeOptions) {
	//FIXME: implement this
}

/*
 *    Copyright (C) 2017 Christian Muehlhaeuser
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

package facebookbee

import (
	"fmt"
	"log"
	"net/url"
	"path"
	"strings"

	"golang.org/x/oauth2"
	oauth2fb "golang.org/x/oauth2/facebook"

	"github.com/muesli/beehive/api"
	"github.com/muesli/beehive/bees"
)

// FacebookBeeFactory is a factory for FacebookBees.
type FacebookBeeFactory struct {
	bees.BeeFactory
}

// New returns a new Bee instance configured with the supplied options.
func (factory *FacebookBeeFactory) New(name, description string, options bees.BeeOptions) bees.BeeInterface {
	bee := FacebookBee{
		Bee: bees.NewBee(name, factory.ID(), description, options),
	}
	bee.ReloadOptions(options)

	return &bee
}

// ID returns the ID of this Bee.
func (factory *FacebookBeeFactory) ID() string {
	return "facebookbee"
}

// Name returns the name of this Bee.
func (factory *FacebookBeeFactory) Name() string {
	return "Facebook"
}

// Description returns the description of this Bee.
func (factory *FacebookBeeFactory) Description() string {
	return "Posts and reacts to events in your Facebook timeline"
}

// Image returns the filename of an image for this Bee.
func (factory *FacebookBeeFactory) Image() string {
	return factory.ID() + ".png"
}

// LogoColor returns the preferred logo background color (used by the admin interface).
func (factory *FacebookBeeFactory) LogoColor() string {
	return "#3a5b9b"
}

// OAuth2AccessToken returns the oauth2 access token.
func (factory *FacebookBeeFactory) OAuth2AccessToken(id, secret, code string) (*oauth2.Token, error) {
	// Get facebook access token.
	conf := &oauth2.Config{
		ClientID:     id,
		ClientSecret: secret,
		RedirectURL:  api.CanonicalURL().String() + "/" + path.Join("oauth2", factory.ID(), id, secret),
		Scopes:       []string{"public_profile", "pages_manage_posts", "publish_to_groups", "pages_read_engagement"},
		Endpoint:     oauth2fb.Endpoint,
	}
	token, err := conf.Exchange(oauth2.NoContext, code)
	if err != nil {
		fmt.Println("Error:", err)
		return nil, err
	}

	return token, nil
}

// Options returns the options available to configure this Bee.
func (factory *FacebookBeeFactory) Options() []bees.BeeOptionDescriptor {
	conf := &oauth2.Config{
		ClientID:     "__client_id__",
		ClientSecret: "__client_secret__",
		RedirectURL:  api.CanonicalURL().String() + "/" + path.Join("oauth2", factory.ID(), "__client_id__", "__client_secret__"),
		Scopes:       []string{"public_profile", "pages_manage_posts", "publish_to_groups", "pages_read_engagement"},
		Endpoint:     oauth2fb.Endpoint,
	}
	u, err := url.Parse(conf.Endpoint.AuthURL)
	if err != nil {
		log.Fatal("Parse:", err)
	}
	parameters := url.Values{}
	parameters.Add("client_id", conf.ClientID)
	parameters.Add("scope", strings.Join(conf.Scopes, " "))
	parameters.Add("redirect_uri", conf.RedirectURL)
	parameters.Add("response_type", "code")
	parameters.Add("state", "beehive")
	u.RawQuery = parameters.Encode()

	opts := []bees.BeeOptionDescriptor{
		{
			Name:        "client_id",
			Description: "App ID for the Facebook API",
			Type:        "string",
		},
		{
			Name:        "client_secret",
			Description: "App Secret for the Facebook API",
			Type:        "string",
		},
		{
			Name:        "access_token",
			Description: "Access token for the Facebook API",
			Type:        "oauth2:" + u.String(),
		},
		{
			Name:        "page_id",
			Description: "Page ID of your Facebook page (see wiki)",
			Type:        "string",
		},
		{
			Name:        "page_access_token",
			Description: "Page access token for the Facebook API (see wiki)",
			Type:        "oauth2:" + u.String(),
		},
	}
	return opts
}

// Actions describes the available actions provided by this Bee.
func (factory *FacebookBeeFactory) Actions() []bees.ActionDescriptor {
	actions := []bees.ActionDescriptor{
		{
			Namespace:   factory.Name(),
			Name:        "post",
			Description: "Submit a new post",
			Options: []bees.PlaceholderDescriptor{
				{
					Name:        "text",
					Description: "Text to post",
					Type:        "string",
					Mandatory:   true,
				},
			},
		},
	}
	return actions
}

// Events describes the available events provided by this Bee.
func (factory *FacebookBeeFactory) Events() []bees.EventDescriptor {
	events := []bees.EventDescriptor{
		{
			Namespace:   factory.Name(),
			Name:        "posted",
			Description: "is triggered when you posted something",
			Options: []bees.PlaceholderDescriptor{
				{
					Name:        "text",
					Description: "text content of the post",
					Type:        "string",
				},
				{
					Name:        "id",
					Description: "ID of the post",
					Type:        "string",
				},
			},
		},
	}
	return events
}

func init() {
	f := FacebookBeeFactory{}
	bees.RegisterFactory(&f)
}

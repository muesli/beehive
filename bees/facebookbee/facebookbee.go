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

// Package facebookbee is a Bee that can interface with Facebook.
package facebookbee

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"jaytaylor.com/html2text"
	"net/http"
	"net/url"
	"path"
	"strconv"
	"time"

	"github.com/huandu/facebook"
	"golang.org/x/oauth2"
	oauth2fb "golang.org/x/oauth2/facebook"

	"github.com/muesli/beehive/api"
	"github.com/muesli/beehive/bees"
)

// FacebookBee is a Bee that can interface with Facebook.
type FacebookBee struct {
	bees.Bee

	clientID        string
	clientSecret    string
	accessToken     string
	pageID          string
	pageAccessToken string

	session *facebook.Session

	evchan chan bees.Event
}

// Run executes the Bee's event loop.
func (mod *FacebookBee) Run(eventChan chan bees.Event) {
	mod.handlePermanentPageToken()
	mod.evchan = eventChan

	since := strconv.FormatInt(time.Now().UTC().Unix(), 10)
	timeout := time.Duration(10 * time.Second)

	for {
		select {
		case <-mod.SigChan:
			return
		case <-time.After(time.Duration(timeout)):
			var err error
			since, err = mod.handleStream(since)
			if err != nil {
				panic(err)
			}
		}

		timeout = time.Minute
	}
}

func (mod *FacebookBee) handlePermanentPageToken() {
	if mod.pageAccessToken != "" {
		return
	}

	mod.Logf("Attempting to fetch long lived user access token")

	longToken := mod.fetchLongLivedUserAccessToken()

	if longToken == "" {
		mod.LogErrorf("No long lived user access token!")
		return
	}

	// mod.Logf("Long lived user access token: \"%s\"", longToken)
	accountID := mod.fetchAccountID(longToken)

	if accountID == "" {
		mod.LogErrorf("No account id!")
		return
	}

	pageToken := mod.fetchPermanentPageAccessToken(accountID, longToken)
	mod.Logf("Permanent pageToken: \"%s\"", pageToken)

	setRes := mod.SetOption("page_access_token", pageToken)

	if !setRes {
		mod.LogErrorf("Permanent pageToken could not be stored!")
	}
}

func (mod *FacebookBee) fetchLongLivedUserAccessToken() string {
	// See https://developers.facebook.com/docs/pages/access-tokens/#get-a-long-lived-user-access-token
	baseURL := "https://graph.facebook.com/oauth/access_token"
	v := url.Values{}
	v.Set("grant_type", "fb_exchange_token")
	v.Set("client_id", mod.clientID)
	v.Set("client_secret", mod.clientSecret)
	v.Set("fb_exchange_token", mod.accessToken)
	graphUrl := baseURL + "?" + v.Encode()

	res, err := http.Get(graphUrl)

	if err != nil || res == nil {
		mod.LogErrorf("Fetching long lived user access token failed: %v", err)
		return ""
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)

	if err != nil {
		mod.LogErrorf("Reading content while fetching long lived user access token failed: %v", err)
		return ""
	}

	// mod.Logf("Long lived user access token result: \"%s\"", body)

	type RequestResult struct {
		AccessToken string `json:"access_token"`
	}

	var tokenRes RequestResult
	err = json.Unmarshal(body, &tokenRes)

	if err != nil {
		mod.LogErrorf("Parsing result while fetching long lived user access token failed: %v", err)
		return ""
	}

	return tokenRes.AccessToken
}

func (mod *FacebookBee) fetchAccountID(accessToken string) string {
	baseURL := "https://graph.facebook.com/v8.0/me"
	v := url.Values{}
	v.Set("access_token", accessToken)
	v.Set("fields", "id")
	graphUrl := baseURL + "?" + v.Encode()

	res, err := http.Get(graphUrl)

	if err != nil || res == nil {
		mod.LogErrorf("Fetching user id failed: %v", err)
		return ""
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)

	if err != nil {
		mod.LogErrorf("Reading content while fetching user id failed: %v", err)
		return ""
	}

	// mod.Logf("user id result: \"%s\"", body)

	type RequestResult struct {
		ID string `json:"id"`
	}

	var tokenRes RequestResult
	err = json.Unmarshal(body, &tokenRes)

	if err != nil {
		mod.LogErrorf("Parsing result while fetching user id failed: %v", err)
		return ""
	}

	return tokenRes.ID
}

func (mod *FacebookBee) fetchPermanentPageAccessToken(accountID string, accessToken string) string {
	// the method in https://developers.facebook.com/docs/pages/access-tokens/#get-a-page-access-token doesn't work!
	// https://github.com/Bnjis/Facebook-permanent-token-generator/blob/master/src/components/Form.js helped a lot
	baseURL := "https://graph.facebook.com/v8.0/" + accountID + "/accounts"
	v := url.Values{}
	v.Set("access_token", accessToken)
	graphUrl := baseURL + "?" + v.Encode()

	// var buf io.ReadWriter
	res, err := http.Get(graphUrl)

	if err != nil || res == nil {
		mod.LogErrorf("Fetching page token failed: %v", err)
		return ""
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)

	if err != nil {
		mod.LogErrorf("Reading content while fetching page token failed: %v", err)
		return ""
	}

	// mod.Logf("Page token result: \"%s\"", body)

	type RequestResultInner struct {
		AccessToken string `json:"access_token"`
	}

	type RequestResult struct {
		Data []RequestResultInner `json:"data"`
	}

	var tokenRes RequestResult
	err = json.Unmarshal(body, &tokenRes)

	if err != nil {
		mod.LogErrorf("Parsing result while fetching page token failed: %v", err)
		return ""
	}

	return tokenRes.Data[0].AccessToken
}

// Action triggers the action passed to it.
func (mod *FacebookBee) Action(action bees.Action) []bees.Placeholder {
	outs := []bees.Placeholder{}
	switch action.Name {
	case "post":
		var text string
		action.Options.Bind("text", &text)
		mod.Logf("Attempting to post \"%s\" to Facebook Page \"%s\"", text, mod.pageID)

		// transform possible html in the text
		textNoHtml, err := html2text.FromString(text, html2text.Options{PrettyTables: true})

		if err == nil {
			text = textNoHtml
		}

		// See https://developers.facebook.com/docs/pages/publishing#before-you-start
		baseURL := "https://graph.facebook.com/" + mod.pageID + "/feed"
		v := url.Values{}
		v.Set("message", text)
		v.Set("access_token", mod.pageAccessToken)
		graphUrl := baseURL + "?" + v.Encode()

		var buf io.ReadWriter
		res, err := http.Post(graphUrl, "", buf)

		if err != nil || res == nil {
			mod.LogErrorf("Posting to Facebook Page failed: %v", err)
			return outs
		}

		defer res.Body.Close()
		body, err := ioutil.ReadAll(res.Body)

		if err != nil {
			mod.LogErrorf("Reading content from post request to Facebook Page failed: %v", err)
			return outs
		}

		mod.Logf("Facebook Page post: \"%s\"", body)

	default:
		panic("Unknown action triggered in " + mod.Name() + ": " + action.Name)
	}

	return outs
}

func (mod *FacebookBee) handleStreamEvent(item map[string]interface{}) {
	mod.Logf("Event: %+v", item)

	if msg, ok := item["message"]; ok {
		ev := bees.Event{
			Bee:  mod.Name(),
			Name: "posted",
			Options: []bees.Placeholder{
				{
					Name:  "text",
					Type:  "string",
					Value: msg,
				},
				{
					Name:  "id",
					Type:  "string",
					Value: item["id"],
				},
			},
		}

		mod.evchan <- ev
	}
}

func (mod *FacebookBee) handleStream(since string) (string, error) {
	conf := &oauth2.Config{
		ClientID:     mod.clientID,
		ClientSecret: mod.clientSecret,
		RedirectURL:  api.CanonicalURL().String() + "/" + path.Join("oauth2", mod.Namespace(), mod.clientID, mod.clientSecret),
		Scopes:       []string{"public_profile", "pages_manage_posts", "publish_to_groups", "pages_read_engagement"},
		Endpoint:     oauth2fb.Endpoint,
	}

	token := oauth2.Token{
		AccessToken: mod.accessToken,
	}

	// Create a client to manage access token life cycle.
	client := conf.Client(oauth2.NoContext, &token)

	// Use OAuth2 client with session.
	mod.session = &facebook.Session{
		Version:    "v2.4",
		HttpClient: client,
	}

	// Use session.
	params := facebook.Params{}
	if since != "" {
		params["since"] = since
	}

	res, err := mod.session.Get("/me/feed", params)
	if err != nil {
		// err can be an facebook API error.
		// if so, the Error struct contains error details.
		if e, ok := err.(*facebook.Error); ok {
			mod.LogErrorf("Error: [message:%v] [type:%v] [code:%v] [subcode:%v]",
				e.Message, e.Type, e.Code, e.ErrorSubcode)
		}
		return since, err
	}

	// mod.Logln("Result:", res)
	// process feed
	events := res.Get("data").([]interface{})
	for _, e := range events {
		mod.handleStreamEvent(e.(map[string]interface{}))
	}

	if res.Get("paging.previous") != nil {
		su, _ := url.Parse(res.Get("paging.previous").(string))
		s := su.Query().Get("since")
		if s != "" {
			return s, nil
		}
	}
	return since, nil
}

// ReloadOptions parses the config options and initializes the Bee.
func (mod *FacebookBee) ReloadOptions(options bees.BeeOptions) {
	mod.SetOptions(options)

	options.Bind("client_id", &mod.clientID)
	options.Bind("client_secret", &mod.clientSecret)
	options.Bind("access_token", &mod.accessToken)
	options.Bind("page_id", &mod.pageID)
	options.Bind("page_access_token", &mod.pageAccessToken)
}

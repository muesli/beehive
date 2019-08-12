/*
 *    Copyright (C) 2016 Sergio Rubio
 *                  2017 Christian Muehlhaeuser
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
 *      Sergio Rubio <rubiojr@frameos.org>
 *      Christian Muehlhaeuser <muesli@gmail.com>
 */

// Package rocketchatbee is a Bee that can connect to Rocketchat.
package rocketchatbee

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type Client struct {
	httpClient *http.Client
	url        string
	userID     string
	authToken  string
}

// Message defines a chat message
type Message struct {
	Text      string `json:"msg"`
	ChannelID string `json:"rid"`
	Alias     string `json:"alias,omitempty"`
}

// MessageContainer is sent to Rocket.Chat to create a chat message
type MessageContainer struct {
	Message Message `json:"message"`
}

// Room is a Channel in Rocket.Chat
type Room struct {
	ID string `json:"_id"`
}

// RoomResponseContainer is used to receive the ID of a room/channel
type RoomResponseContainer struct {
	Room Room `json:"room"`
}

// ErrorResponse
type ErrorResponse struct {
	Error string `json:"error"`
}

func NewClient(url, userID, authToken string) *Client {
	return &Client{
		&http.Client{
			Timeout: 30 * time.Second,
		},
		url,
		userID,
		authToken,
	}
}

func (c *Client) newRequest(method, path string, body interface{}) (*http.Request, error) {
	url := fmt.Sprintf("%s/%s", c.url, path)

	var buf io.ReadWriter
	if body != nil {
		buf = new(bytes.Buffer)
		err := json.NewEncoder(buf).Encode(body)
		if err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequest(method, url, buf)
	if err != nil {
		return nil, err
	}

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("X-User-Id", c.userID)
	req.Header.Set("X-Auth-Token", c.authToken)

	return req, nil
}

func (c *Client) do(req *http.Request, v interface{}) (*http.Response, error) {
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode > 399 {
		var errorResponse ErrorResponse
		err = json.NewDecoder(resp.Body).Decode(&errorResponse)
		if err != nil {
			return nil, fmt.Errorf("request to '%s' failed with http %d: %s", resp.Request.URL, resp.StatusCode, err)
		}
		return nil, fmt.Errorf("request to '%s' failed with http code %d: %s", resp.Request.URL, resp.StatusCode, errorResponse.Error)
	}

	// nil for empty responses
	if v != nil {
		err = json.NewDecoder(resp.Body).Decode(v)
	}
	return resp, err
}

func (c *Client) SendMessage(channel, text, alias string) error {
	var (
		channelID string
	)

	channelID, err := c.GetChannelID(channel)
	if err != nil {
		return err
	}

	msg := &MessageContainer{
		Message{
			Text:      text,
			ChannelID: channelID,
			Alias:     alias,
		},
	}

	req, err := c.newRequest("POST", "/api/v1/chat.sendMessage", msg)
	if err != nil {
		return err
	}

	// errors are handeled in do
	_, err = c.do(req, nil)
	return err
}

// getChannelID returns the ID of a channel
func (c *Client) GetChannelID(name string) (string, error) {

	req, err := c.newRequest("GET", "/api/v1/rooms.info", nil)

	// add query parameter roomName
	q := req.URL.Query()
	q.Add("roomName", name)
	req.URL.RawQuery = q.Encode()

	rc := &RoomResponseContainer{}
	_, err = c.do(req, rc)
	if err != nil {
		return "", err
	}
	return rc.Room.ID, nil
}

func (c *Client) TestConnection() error {

	req, err := c.newRequest("GET", "/api/v1/me", nil)
	if err != nil {
		return err
	}

	_, err = c.do(req, nil)
	if err != nil {
		return err
	}

	return nil
}

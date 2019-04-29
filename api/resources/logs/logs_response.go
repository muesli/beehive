/*
 *    Copyright (C) 2015-2019 Christian Muehlhaeuser
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

package logs

import (
	"time"

	"github.com/muesli/beehive/bees"

	"github.com/muesli/smolder"
)

// LogResponse is the common response to 'log' requests
type LogResponse struct {
	smolder.Response

	Logs []logInfoResponse `json:"logs,omitempty"`
	logs []*bees.LogMessage
}

type logInfoResponse struct {
	ID        string    `json:"id"`
	Bee       string    `json:"bee"`
	Level     int64     `json:"level"`
	Message   string    `json:"message"`
	Timestamp time.Time `json:"timestamp"`
}

// Init a new response
func (r *LogResponse) Init(context smolder.APIContext) {
	r.Parent = r
	r.Context = context

	r.Logs = []logInfoResponse{}
}

// AddLog adds a log to the response
func (r *LogResponse) AddLog(log *bees.LogMessage) {
	r.logs = append(r.logs, log)
	r.Logs = append(r.Logs, prepareLogResponse(r.Context, log))
}

// EmptyResponse returns an empty API response for this endpoint if there's no data to respond with
func (r *LogResponse) EmptyResponse() interface{} {
	if len(r.logs) == 0 {
		var out struct {
			Logs interface{} `json:"logs"`
		}
		out.Logs = []logInfoResponse{}
		return out
	}
	return nil
}

func prepareLogResponse(context smolder.APIContext, log *bees.LogMessage) logInfoResponse {
	//	ctx := context.(*context.APIContext)
	resp := logInfoResponse{
		ID:        (*log).ID,
		Bee:       (*log).Bee,
		Level:     int64((*log).MessageType),
		Message:   (*log).Message,
		Timestamp: (*log).Timestamp,
	}

	return resp
}

/*
 *    Copyright (C) 2017-2019 Christian Muehlhaeuser
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

// Package templatehelper provides a func-map of common template functions
package templatehelper

import (
	"encoding/json"
	"errors"
	htmlTemplate "html/template"
	"regexp"
	"strings"
	"text/template"
	"time"
)

// FuncMap contains a few convenient template helpers
var (
	FuncMap = template.FuncMap{
		"JSON": func(values ...interface{}) htmlTemplate.JS {
			json, _ := json.Marshal(values)
			return htmlTemplate.JS(json)
		},
		"Left": func(s string, n int) string {
			if n > len(s) {
				n = len(s)
			}
			return s[:n]
		},
		"Matches": func(s string, pattern string) (bool, error) {
			return regexp.MatchString(pattern, s)
		},
		"Mid": func(s string, left int, values ...int) string {
			if left > len(s) {
				left = len(s)
			}

			if len(values) != 0 {
				right := values[0]
				if right > len(s) {
					right = len(s)
				}
				return s[left:right]
			}
			return s[left:]
		},
		"Right": func(s string, right int) string {
			left := len(s) - right
			if left < 0 {
				left = 0
			}
			return s[left:]
		},
		"Last": func(items []string) (string, error) {
			if len(items) == 0 {
				return "", errors.New("cannot get last element from empty slice")
			}
			return items[len(items)-1], nil
		},
		"ContainsAny": func(target, subs string, other ...string) bool {
			if len(other) == 0 {
				return strings.ContainsAny(target, subs)
			}
			for _, another := range other {
				if strings.Contains(target, another) {
					return true
				}
			}
			return strings.Contains(target, subs)
		},
		// strings functions
		"Compare":      strings.Compare, // 1.5+ only
		"Contains":     strings.Contains,
		"Count":        strings.Count,
		"EqualFold":    strings.EqualFold,
		"HasPrefix":    strings.HasPrefix,
		"HasSuffix":    strings.HasSuffix,
		"Index":        strings.Index,
		"IndexAny":     strings.IndexAny,
		"Join":         strings.Join,
		"LastIndex":    strings.LastIndex,
		"LastIndexAny": strings.LastIndexAny,
		"Repeat":       strings.Repeat,
		"Replace":      strings.Replace,
		"Split":        strings.Split,
		"SplitAfter":   strings.SplitAfter,
		"SplitAfterN":  strings.SplitAfterN,
		"SplitN":       strings.SplitN,
		"TimeNow":      time.Now,
		"Title":        strings.Title,
		"ToLower":      strings.ToLower,
		"ToTitle":      strings.ToTitle,
		"ToUpper":      strings.ToUpper,
		"Trim":         strings.Trim,
		"TrimLeft":     strings.TrimLeft,
		"TrimPrefix":   strings.TrimPrefix,
		"TrimRight":    strings.TrimRight,
		"TrimSpace":    strings.TrimSpace,
		"TrimSuffix":   strings.TrimSuffix,
	}
)

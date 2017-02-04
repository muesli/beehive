package templatehelper

import (
	"strings"
	"text/template"
)

// FuncMap contains all the common string helpers
var (
	FuncMap = template.FuncMap{
		"Left": func(values ...interface{}) string {
			return values[0].(string)[:values[1].(int)]
		},
		"Mid": func(values ...interface{}) string {
			if len(values) > 2 {
				return values[0].(string)[values[1].(int):values[2].(int)]
			}
			return values[0].(string)[values[1].(int):]
		},
		"Right": func(values ...interface{}) string {
			return values[0].(string)[len(values[0].(string))-values[1].(int):]
		},
		"Last": func(values ...interface{}) string {
			return values[0].([]string)[len(values[0].([]string))-1]
		},
		// strings functions
		// "Compare":      strings.Compare, // 1.5+ only
		"Contains":     strings.Contains,
		"ContainsAny":  strings.ContainsAny,
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

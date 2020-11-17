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

// Package templatefilter provides a template-based filter.
package templatefilter

import (
	"bytes"
	"strings"
	"text/template"

	"github.com/muesli/beehive/filters"
	"github.com/muesli/beehive/templatehelper"
)

// TemplateFilter is a template-based filter.
type TemplateFilter struct {
}

// Name returns the name of this Filter.
func (filter *TemplateFilter) Name() string {
	return "template"
}

// Description returns the description of this Filter.
func (filter *TemplateFilter) Description() string {
	return "This filter passes when a template-if returns true"
}

// Passes returns true when the Filter matched the data.
func (filter *TemplateFilter) Passes(data map[string]interface{}, v string) bool {
	var res bytes.Buffer

	if strings.Contains(v, "{{test") {
		v = strings.Replace(v, "{{test", "{{if", -1)
		v += "true{{end}}"
	}

	tmpl, err := template.New("_" + v).Funcs(templatehelper.FuncMap).Parse(v)
	if err == nil {
		err = tmpl.Execute(&res, data)
	}
	if err != nil {
		panic(err)
	}

	return strings.TrimSpace(res.String()) == "true"
}

func init() {
	f := TemplateFilter{}
	filters.RegisterFilter(&f)
}

/*
 *    Copyright (C) 2014 Daniel 'grindhold' Brendle
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
 *      Daniel 'grindhold' Brendle <grindhold@skarphed.org>
 */

package booleanfilter

import (
	"github.com/muesli/beehive/filters"
	_"github.com/muesli/beehive/modules"
	"log"
    "errors"
)

const (
	OR  = "or"
	AND = "and"
	XOR = "xor"
	NOT = "not"
)

type BooleanFilter struct {
    filters.Filter
}

func (filter *BooleanFilter) SetOptions(options []filters.FilterOption) {
    filter.Options = options
}

func (filter *BooleanFilter) GetOptions() []filters.FilterOption{
    return filter.Options
}

func (filter *BooleanFilter) GetOption (opt string) (*filters.FilterOption, error) {
    for x := range(filter.GetOptions()){
        if filter.GetOptions()[x].Name == opt {
            return &filter.GetOptions()[x], nil
        }
    }
    return nil, errors.New("option not found")
}

func (filter *BooleanFilter) Name() string {
	return "boolean"
}

func (filter *BooleanFilter) Description() string {
	return "This filter allows boolean operations on 1 (NOT) or more (AND, OR, XOR) filters"
}

func (filter *BooleanFilter) unmarshalFilterOption (d map[string]interface{}) filters.FilterOption {
    fo := filters.FilterOption{
                Name: d["Name"].(string),
                Type: d["Type"].(string),
                Value: d["Value"],
    }
    return fo
}

func (filter *BooleanFilter) unmarshalFilter(d map[string]interface{}) filters.Filter{
    var options []filters.FilterOption

    for _, y := range(d["Options"].([]map[string]interface{})){
        options = append(options, filter.unmarshalFilterOption(y.(map[string]interface{})))
    }

    f := filters.Filter{
        Name: d["Name"].(string),
        Options: options,
    }
    return f
}

func (filter *BooleanFilter) unmarshalFilters(d []map[string]interface{}) []filters.Filter{
    var filters []filters.Filter
    for x := range(d){
        filters = append(filters, filter.unmarshalFilter(d[x]))
    }
    return filters
}

func (filter *BooleanFilter) Passes(data map[string]interface{}) bool {
    fopt, err := filter.GetOption("filters")
    if err != nil {
        log.Println(err)
        return false
    }
    log.Println("here")
    log.Println(fopt.Value)
    rfl := filter.unmarshalFilters(fopt.Value.([]map[string]interface{}))
    log.Println("there")
    if len(rfl) < 1 {
        log.Println("No filters added -> Returning false")
        return false
    }

    var filterList []filters.FilterInterface
    for f := range rfl {
        filterList = append(filterList, (*filters.GetFilter(rfl[f].Name)))
        filterList[f].SetOptions(rfl[f].Options)
    }
    fopt, err = filter.GetOption("operation")
    if err != nil {
        log.Println(err)
        return false
    }
    operation := fopt.Value.(string)
	switch operation {
	case OR:
		for f := range filterList {
			if filterList[f].Passes(data) {
				return true
			}
		}
		return false
	case AND:
		for f := range filterList {
			if !filterList[f].Passes(data) {
				return false
			}
		}
		return true
	case XOR:
		true_found := 0
		for f := range filterList {
			if filterList[f].Passes(data) {
				true_found++
			}
			if true_found > 1 {
				return false
			}
		}
		return true_found == 1
	case NOT:
		return !filterList[0].Passes(data)
	default:
		log.Println("Cannot join filters: No valid operation!")
	}
	return false
}

func init() {
	f := BooleanFilter{}
	filters.RegisterFilter(&f)
}

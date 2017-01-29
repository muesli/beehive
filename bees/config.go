/*
 *    Copyright (C) 2014-2017 Christian Muehlhaeuser
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

// Package bees is Beehive's central module system.
package bees

import "errors"

// BeeConfig contains all settings for a single Bee.
type BeeConfig struct {
	Name        string
	Class       string
	Description string
	Options     BeeOptions
}

// NewBeeConfig validates a configuration and sets up a new BeeConfig
func NewBeeConfig(name, class, description string, options BeeOptions) (BeeConfig, error) {
	if len(name) == 0 {
		return BeeConfig{}, errors.New("A Bee's name can't be empty")
	}

	b := GetBee(name)
	if b != nil {
		return BeeConfig{}, errors.New("A Bee with that name already exists")
	}

	f := GetFactory(class)
	if f == nil {
		return BeeConfig{}, errors.New("Invalid class specified")
	}

	return BeeConfig{
		Name:        name,
		Class:       class,
		Description: description,
		Options:     options,
	}, nil
}

// BeeConfigs returns configs for all Bees.
func BeeConfigs() []BeeConfig {
	bs := []BeeConfig{}
	for _, b := range bees {
		bs = append(bs, (*b).Config())
	}

	return bs
}

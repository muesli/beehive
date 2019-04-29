/*
 *    Copyright (C) 2019 Christian Muehlhaeuser
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

var (
	ctx = NewContext()
)

type Context struct {
	state map[*Bee]map[string]interface{}
}

func NewContext() *Context {
	return &Context{
		state: make(map[*Bee]map[string]interface{}),
	}
}

func (c *Context) Set(bee *Bee, key string, value interface{}) {
	if _, ok := c.state[bee]; !ok {
		c.state[bee] = make(map[string]interface{})
	}
	c.state[bee][key] = value
}

func (c *Context) Value(bee *Bee, key string) interface{} {
	return c.state[bee][key]
}

func (c *Context) FillMap(m map[string]interface{}) {
	cd := make(map[string]interface{})
	for bee, d := range c.state {
		cd[bee.Name()] = d
	}
	m["context"] = cd
}

func (bee *Bee) ContextSet(key string, value interface{}) {
	ctx.Set(bee, key, value)
}

func (bee *Bee) ContextValue(key string) interface{} {
	return ctx.Value(bee, key)
}

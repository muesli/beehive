/*
 *    Copyright (C) 2014 Christian Muehlhaeuser
 *                  2014 Michael Wendland
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
 *      Michael Wendland <michiwend@michiwend.com>
 */

package main

import (
	"log"

	"github.com/muesli/beehive/app"

	"github.com/muesli/beehive/modules"
	_ "github.com/muesli/beehive/modules/hellobee"
	_ "github.com/muesli/beehive/modules/ircbee"
	_ "github.com/muesli/beehive/modules/webbee"
)

func main() {
	// Parse command-line args for all registered modules
	app.Run()

	log.Println("Beehive is buzzing...")

	// Initialize modules
	modules.StartModules()

	// Keep app alive
	ch := make(chan bool)
	<-ch
}

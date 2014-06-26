/*
 *    Copyright (C) 2014 Christian Muehlhaeuser
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

package containsfilter

import (
	"testing"
)

var (
	haystack = "foobar"
	needle   = "oba"
	fail     = "abo"
)

func TestContainsFilter(t *testing.T) {
	f := ContainsFilter{}

	if !f.Passes(haystack, needle) {
		t.Error("ContainsFilter fails on string comparison")
	}
	if f.Passes(haystack, fail) {
		t.Error("ContainsFilter fails on string comparison")
	}
}

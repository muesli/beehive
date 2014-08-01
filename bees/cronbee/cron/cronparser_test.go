/*
 *        Copyright (C) 2014 Stefan 'glaxx' Luecke
 *
 *        This program is free software: you can redistribute it and/or modify
 *        it under the terms of the GNU Affero General Public License as published
 *        by the Free Software Foundation, either version 3 of the License, or
 *        (at your option) any later version.
 *
 *        This program is distributed in the hope that it will be useful,
 *        but WITHOUT ANY WARRANTY; without even the implied warranty of
 *        MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 *        GNU Affero General Public License for more details.
 *
 *        You should have received a copy of the GNU Affero General Public License
 *        along with this program.      If not, see <http://www.gnu.org/licenses/>.
 *
 *        Authors:	Stefan Luecke <glaxx@glaxx.net>
 */

package cron

 import (
	"fmt"
	"testing"
)

 func Test_add(t *testing.T) {
	if add(40, 20, 60) != 0 {
		t.Error("crontime: add()-Test failed, 1")
	} else if add(13, 13, 60) != 26 {
		t.Error("crontime: add()-Test failed, 2")
	} else if add(13, 14, 60) != 27 {
		t.Error("crontime: add()-Test failed, 3")
	} else if add(20, 40, 60) != 0 {
		t.Error("crontime: add()-Test failed, 4")
	} else if add(60, 60, 60) != 0 {
		t.Error("crontime: add()-Test failed, 5")
	} else {
		t.Log("crontime: add()-Test passed")
	}
}

func Test_absolute_over_breakpoint(t *testing.T) {
	if absolute_over_breakpoint(13, 17, 60) != 4 {
		t.Error("crontime: absolute_over_breakpoint()-Test failed, 1")
	} else if absolute_over_breakpoint(17, 13, 60) != 56 {
		t.Error("crontime: absolute_over_breakpoint()-Test failed, 2")
	} else if absolute_over_breakpoint(20, 1, 24) != 5 {
		t.Error("crontime: absolute_over_breakpoint()-Test failed, 3")
	} else if absolute_over_breakpoint(0, 0, 24) != 0 {
		t.Error("crontime: absolute_over_breakpoint()-Test failed, 4")
	} else if absolute_over_breakpoint(13, 13, 60) != 0 {
		t.Error("crontime: absolute_over_breakpoint()-Test failed, 5")
	} else {
		t.Log("crontime: absolute_over_breakpoint()-Test passed")
	}
}

func Test_value_range(t *testing.T) {
	if !IntsEquals(value_range(2, 5, 60), []int{2, 3, 4, 5}) {
		t.Error("crontime: value_range()-Test failed, 1")
		fmt.Println(value_range(2, 5, 60), []int{2, 3, 4, 5})
	} else if !IntsEquals(value_range(59, 1, 60), []int{59, 0, 1}) {
		t.Error("crontime: value_range()-Test failed, 2s")
	} else {
		t.Log("crontime: value_range()-Test passed")
	}
}

func Test_periodic(t *testing.T) {
	if !IntsEquals(periodic(20, 0, 60), []int{0, 20, 40}) {
		t.Error("crontime: periodic()-Test failed, 1")
	} else if !IntsEquals(periodic(10, 0, 60), []int{0, 10, 20, 30, 40, 50}){
		t.Error("crontime: periodic()-Test failed, 2")
	} else if !IntsEquals(periodic(23, 0, 24), []int{0, 23, 22, 21, 20, 19, 18, 17, 16, 15, 14, 13, 12, 11, 10, 9, 8, 7, 6, 5, 4, 3, 2, 1}){
		t.Error("crontime: periodic()-Test failed, 3")
	} else if !IntsEquals(periodic(10, 15, 60), []int{15, 25, 35, 45, 55, 5}){
		t.Error("crontime: periodic()-Test failed, 4")
	} else if !IntsEquals(periodic(10, 0, 60), []int{0, 10, 20, 30, 40, 50}){
		t.Error("crontime: periodic()-Test failed, 5")
	} else if !IntsEquals(periodic(10, 0, 60), []int{0, 10, 20, 30, 40, 50}){
		t.Error("crontime: periodic()-Test failed, 6")
	} else if !IntsEquals(periodic(10, 0, 60), []int{0, 10, 20, 30, 40, 50}){
		t.Error("crontime: periodic()-Test failed, 7")
	} else {
		t.Log("crontime: periodic()-Test passed")
	}
}

func Test_check_syntax(t *testing.T) {
	if !check_syntax("13") {
		t.Error("crontime: check_syntax()-Test failed, 1")
	} else if !check_syntax("*/05") {
		t.Error("crontime: check_syntax()-Test failed, 2")
	} else if !check_syntax("13,23") {
		t.Error("crontime: check_syntax()-Test failed, 3")
	} else if !check_syntax("14-15") {
		t.Error("crontime: check_syntax()-Test failed, 4")
	} else if !check_syntax("*") {
		t.Error("crontime: check_syntax()-Test failed, 5")
	} else if check_syntax("13-13-13") {
		t.Error("crontime: check_syntax()-Test failed, 6")
	} else if check_syntax("13,") {
		t.Error("crontime: check_syntax()-Test failed, 7")
	} else if check_syntax("13-") {
		t.Error("crontime: check_syntax()-Test failed, 8")
	} else if check_syntax("5/*") {
		t.Error("crontime: check_syntax()-Test failed, 9")
	} else if check_syntax("**") {
		t.Error("crontime: check_syntax()-Test failed, 10")
	} else if check_syntax("asdf") {
		t.Error("crontime: check_syntax()-Test failed, 11")
	} else if check_syntax(",,") {
		t.Error("crontime: check_syntax()-Test failed, 12")
	} else if check_syntax("*/*/*/*/") {
		t.Error("crontime: check_syntax()-Test failed, 13")
	} else if check_syntax("*/*/*/23") {
		t.Error("crontime: check_syntax()-Test failed, 14")
	} else if check_syntax("") {
		t.Error("crontime: check_syntax()-Test failed, 15")
	} else if check_syntax("*/5") {
		t.Error("crontime: check_syntax()-Test failed, 16")
	} else {
		t.Log("crontime: check_syntax()-Test passed")
	}
}

func Test_parseIRange(t *testing.T) {
	var sectestA crontime
	var sectestB crontime
	sectestA.second = []int{2, 5, 23}
	sectestB.parseIRange("02,05,23", 0)

	var mintestA crontime
	var mintestB crontime
	mintestA.minute = []int{2, 5, 23}
	mintestB.parseIRange("02,05,23", 1)

	var hourtestA crontime
	var hourtestB crontime
	hourtestA.hour = []int{2, 5, 23}
	hourtestB.parseIRange("02,05,23", 2)

	var dowtestA crontime
	var dowtestB crontime
	dowtestA.dow = []int{2, 5, 6}
	dowtestB.parseIRange("02,05,06", 3)

	var domtestA crontime
	var domtestB crontime
	domtestA.dom = []int{2, 5, 23}
	domtestB.parseIRange("02,05,23", 4)

	var monthtestA crontime
	var monthtestB crontime
	monthtestA.month = []int{2, 5, 23}
	monthtestB.parseIRange("02,05,23", 5)

	if !CrontimeEquals(sectestA, sectestB) {
		t.Error("crontime: parseIRange()-Test failed, 1")
	} else if !CrontimeEquals(mintestA, mintestB) {
		t.Error("crontime: parseIRange()-Test failed, 2")
	} else if !CrontimeEquals(hourtestA, hourtestB) {
		t.Error("crontime: parseIRange()-Test failed, 3")
	} else if !CrontimeEquals(dowtestA, dowtestB) {
		t.Error("crontime: parseIRange()-Test failed, 4")
	} else if !CrontimeEquals(domtestA, domtestB) {
		t.Error("crontime: parseIRange()-Test failed, 5")
	} else if !CrontimeEquals(monthtestA, monthtestB) {
		t.Error("crontime: parseIRange()-Test failed, 6")
	} else {
		t.Log("crontime: parseIRange()-Test passed")
	}
}

func Test_parseRange(t *testing.T) {
	var sectestA crontime
	var sectestB crontime
	sectestA.second = []int{59, 0, 1}
	sectestB.parseRange("59-01", 0)

	var mintestA crontime
	var mintestB crontime
	mintestA.minute = []int{59, 0, 1}
	mintestB.parseRange("59-01", 1)

	var hourtestA crontime
	var hourtestB crontime
	hourtestA.hour = []int{23, 0, 1}
	hourtestB.parseRange("23-01", 2)

	var dowtestA crontime
	var dowtestB crontime
	dowtestA.dow = []int{6, 0, 1}
	dowtestB.parseRange("06-01", 3)

	var domtestA crontime
	var domtestB crontime
	domtestA.dom = []int{31, 0, 1}
	domtestB.parseRange("31-01", 4)

	var monthtestA crontime
	var monthtestB crontime
	monthtestA.month = []int{12, 0, 1}
	monthtestB.parseRange("12-01", 5)

	if !CrontimeEquals(sectestA, sectestB) {
		t.Error("crontime: parseRange()-Test failed, 1")
	} else if !CrontimeEquals(mintestA, mintestB) {
		t.Error("crontime: parseRange()-Test failed, 2")
	} else if !CrontimeEquals(hourtestA, hourtestB) {
		t.Error("crontime: parseRange()-Test failed, 3")
	} else if !CrontimeEquals(dowtestA, dowtestB) {
		t.Error("crontime: parseRange()-Test failed, 4")
	} else if !CrontimeEquals(domtestA, domtestB) {
		t.Error("crontime: parseRange()-Test failed, 5")
	} else if !CrontimeEquals(monthtestA, monthtestB) {
		t.Error("crontime: parseRange()-Test failed, 6")
	} else {
		t.Log("crontime: parseRange()-Test passed")
	}
}

func Test_parseIgnore(t *testing.T) {
	var sectestA crontime
	var sectestB crontime
	sectestA.second = []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32, 33, 34, 35, 36, 37, 38, 39, 40, 41, 42, 43, 44, 45, 46, 47, 48, 49, 50, 51, 52, 53, 54, 55, 56, 57, 58, 59}
	sectestB.parseIgnore(0)

	var mintestA crontime
	var mintestB crontime
	mintestA.minute = []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32, 33, 34, 35, 36, 37, 38, 39, 40, 41, 42, 43, 44, 45, 46, 47, 48, 49, 50, 51, 52, 53, 54, 55, 56, 57, 58, 59}
	mintestB.parseIgnore(1)

	var hourtestA crontime
	var hourtestB crontime
	hourtestA.hour = []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23}
	hourtestB.parseIgnore(2)

	var dowtestA crontime
	var dowtestB crontime
	dowtestA.dow = []int{0, 1, 2, 3, 4, 5, 6}
	dowtestB.parseIgnore(3)

	var domtestA crontime
	var domtestB crontime
	domtestA.dom = []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31}
	domtestB.parseIgnore(4)

	var monthtestA crontime
	var monthtestB crontime
	monthtestA.month = []int{1, 2, 3, 4, 5, 6 ,7, 8, 9, 10, 11, 12}
	monthtestB.parseIgnore(5)

	if !CrontimeEquals(sectestA, sectestB) {
		t.Error("crontime: parseIgnore()-Test failed, 1")
	} else if !CrontimeEquals(mintestA, mintestB) {
		t.Error("crontime: parseIgnore()-Test failed, 2")
	} else if !CrontimeEquals(hourtestA, hourtestB) {
		t.Error("crontime: parseIgnore()-Test failed, 3")
	} else if !CrontimeEquals(dowtestA, dowtestB) {
		t.Error("crontime: parseIgnore()-Test failed, 4")
	} else if !CrontimeEquals(domtestA, domtestB) {
		t.Error("crontime: parseIgnore()-Test failed, 5")
	} else if !CrontimeEquals(monthtestA, monthtestB) {
		t.Error("crontime: parseIgnore()-Test failed, 6")
	} else {
		t.Log("crontime: parseIgnore()-Test passed")
	}
}

func CrontimeEquals(a, b crontime) bool {
	if len(a.second) != len(b.second) {
		return false
	} else if len(a.minute) != len(b.minute) {
		return false
	} else if len(a.hour) != len(b.hour) {
		return false
	} else if len(a.dow) != len(b.dow) {
		return false
	} else if len(a.dom) != len(b.dom) {
		return false
	} else if len(a.month) != len(b.month) {
		return false
	} else if !IntsEquals(a.second, b.second) {
		return false
	} else if !IntsEquals(a.minute, b.minute) {
		return false
	} else if !IntsEquals(a.hour, b.hour) {
		return false
	} else if !IntsEquals(a.dow, b.dow) {
		return false
	} else if !IntsEquals(a.dom, b.dom) {
		return false
	} else if !IntsEquals(a.month, b.month) {
		return false
	} else {
		return true
	}
}

func IntsEquals(a, b []int) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}
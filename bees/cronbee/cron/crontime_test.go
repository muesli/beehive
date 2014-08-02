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
	"testing"
//	"fmt"
)
/*
func Test_NextEvent(t *testing.T) {
	testStr := [6]string{"*", "*", "*", "*", "*", "*"}
	testA := ParseInput(testStr)
	if !(int(testA.NextEvent()) != 1) {
		t.Error("crontime: NextEvent failed, 1")
	} else {
	t.Log("passed")
	}
}*/

func Test_IsLeapYear(t *testing.T) {
	if !isLeapYear(2012) {
		t.Error("crontime: isLeapYear failed, 1")
	} else if !isLeapYear(2016) {
		t.Error("crontime: isLeapYear failed, 2")
	} else if isLeapYear(2013) {
		t.Error("crontime: isLeapYear failed, 3")
	} else if isLeapYear(2014) {
		t.Error("crontime: isLeapYear failed, 4")
	} else if isLeapYear(2015) {
		t.Error("crontime: isLeapYear failed, 5")
	} else if !isLeapYear(2400) {
		t.Error("crontime: isLeapYear failed, 6")
	} else {
		t.Log("crontime: isLeapYear passed")
	}
}

func Benchmark_NextEventA(b *testing.B) {
	var ipstr [6]string
	ipstr = [6]string{"*/30", "*", "*", "*", "*", "*"}
	cr := ParseInput(ipstr)
	for i := 0; i < b.N; i++{
		cr.nextEvent()
	}
}

func Benchmark_NextEventB(b *testing.B) {
	var ipstr [6]string
	ipstr = [6]string{"*", "*", "*", "*", "*", "07"}
	cr := ParseInput(ipstr)
	for i := 0; i < b.N; i++{
		cr.nextEvent()
	}
}

func Benchmark_NextEventC(b *testing.B) {
	var ipstr [6]string
	ipstr = [6]string{"*/30", "*", "*", "*", "*", "01-12"}
	cr := ParseInput(ipstr)
	for i := 0; i < b.N; i++{
		cr.nextEvent()
	}
}

func Benchmark_NextEventD(b *testing.B) {
	var ipstr [6]string
	ipstr = [6]string{"*", "*", "*", "*", "*", "*"}
	cr := ParseInput(ipstr)
	for i := 0; i < b.N; i++{
		cr.nextEvent()
	}
}

func Benchmark_NextEventE(b *testing.B) {
	var ipstr [6]string
	ipstr = [6]string{"*", "*", "*", "02", "*", "01-04"}
	cr := ParseInput(ipstr)
	for i := 0; i < b.N; i++{
		cr.nextEvent()
	}
}

func Benchmark_ParseInputA(b *testing.B) {
	var ipstr [6]string
	ipstr = [6]string{"*/30", "*", "*", "*", "*", "07"}
	for i := 0; i < b.N; i++{
		ParseInput(ipstr)
	}
}

func Benchmark_ParseInputB(b *testing.B) {
	var ipstr [6]string
	ipstr = [6]string{"*", "*", "*", "*", "*", "*"}
	for i := 0; i < b.N; i++{
		ParseInput(ipstr)
	}
}

func Benchmark_ParseInputC(b *testing.B) {
	var ipstr [6]string
	ipstr = [6]string{"15-47", "*/04", "02,04", "*", "*", "*"}
	for i := 0; i < b.N; i++{
		ParseInput(ipstr)
	}
}

func Benchmark_parsePeriodic(b *testing.B) {
	var cr crontime
	j := 0
	for i := 0; i < b.N; i++{
		cr.parsePeriodic("*/05", j)
		j = add(j, 1, 6)
	}
}

func Benchmark_check_syntaxA(b *testing.B) {
	for i := 0; i < b.N; i++{
		check_syntax("23,24,00,01,02,03,04,05,06,07")
	}
}

func Benchmark_check_syntaxB(b *testing.B) {
	for i := 0; i < b.N; i++{
		check_syntax("01-23")
	}
}

func Benchmark_check_syntaxC(b *testing.B) {
	for i := 0; i < b.N; i++{
		check_syntax("*/05")
	}
}

func Benchmark_check_syntaxD(b *testing.B) {
	for i := 0; i < b.N; i++{
		check_syntax("*")
	}
}

func Benchmark_check_syntaxE(b *testing.B) {
	for i := 0; i < b.N; i++{
		check_syntax("13")
	}
}

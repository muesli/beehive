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
	"fmt"
	"time"
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

func Test_calculateEventA(t *testing.T) {
	ipstr := [6]string{"*/30", "*", "*", "*", "*", "*"}
	baseTime := time.Date(2014, time.Month(1), 1, 0, 0, 0, 0, time.Local)
	expectedTime := time.Date(2014, time.Month(1), 1, 0, 0, 30, 10000, time.Local)
	cr := ParseInput(ipstr)
	if cr.calculateEvent(baseTime) != expectedTime {
		fmt.Println(cr.calculatedTime, expectedTime)
		t.Error("crontime: calculateEvent failed, 1")
	} else {
		t.Log("crontime: calculateEvent passed")
	}
}

func Test_calculateEventB(t *testing.T) {
	ipstr := [6]string{"*", "*", "*", "*", "*", "*"}
	baseTime := time.Date(2014, time.Month(1), 1, 0, 0, 0, 0, time.Local)
	expectedTime := time.Date(2014, time.Month(1), 1, 0, 0, 1, 10000, time.Local)
	cr := ParseInput(ipstr)
	if cr.calculateEvent(baseTime) != expectedTime {
		fmt.Println(cr.calculatedTime, expectedTime)
		t.Error("crontime: calculateEvent failed, 2")
	} else {
		t.Log("crontime: calculateEvent passed")
	}
}

func Test_calculateEventC(t *testing.T) {
	ipstr := [6]string{"*", "*", "*", "*", "29", "02"}
	baseTime := time.Date(2014, time.Month(1), 1, 0, 0, 0, 0, time.Local)
	expectedTime := time.Date(2016, time.Month(2), 29, 0, 0, 0, 10000, time.Local)
	cr := ParseInput(ipstr)
	if cr.calculateEvent(baseTime) != expectedTime {
		fmt.Println(cr.calculatedTime, expectedTime)
		t.Error("crontime: calculateEvent failed, 3")
	} else {
		t.Log("crontime: calculateEvent passed")
	}
}

func Test_calculateEventD(t *testing.T) {
	ipstr := [6]string{"*", "*", "03", "*", "*", "*"}
	baseTime :=     time.Date(2014, time.Month(1), 1, 4, 0, 0, 0, time.Local)
	expectedTime := time.Date(2014, time.Month(1), 2, 3, 0, 0, 10000, time.Local)
	cr := ParseInput(ipstr)
	if cr.calculateEvent(baseTime) != expectedTime {
		fmt.Println(cr.calculatedTime, expectedTime)
		t.Error("crontime: calculateEvent failed, 4")
	} else {
		t.Log("crontime: calculateEvent passed")
	}
}

func Test_calculateEventE(t *testing.T) {
	ipstr := [6]string{"*", "*", "*", "02", "*", "*"}
	baseTime :=     time.Date(2014, time.Month(1), 1, 0, 0, 0, 0, time.Local)
	expectedTime := time.Date(2014, time.Month(1), 7, 0, 0, 0, 10000, time.Local)
	cr := ParseInput(ipstr)
	if cr.calculateEvent(baseTime) != expectedTime {
		fmt.Println(cr.calculatedTime, expectedTime)
		t.Error("crontime: calculateEvent failed, 5")
	} else {
		t.Log("crontime: calculateEvent passed")
	}
}

func Benchmark_calculateEventA(b *testing.B) {
	var ipstr [6]string
	ipstr = [6]string{"*/30", "*", "*", "*", "*", "*"}
	cr := ParseInput(ipstr)
	for i := 0; i < b.N; i++{
		cr.calculateEvent(time.Now())
	}
}

func Benchmark_calculateEventB(b *testing.B) {
	var ipstr [6]string
	ipstr = [6]string{"*", "*", "*", "*", "*", "07"}
	cr := ParseInput(ipstr)
	for i := 0; i < b.N; i++{
		cr.calculateEvent(time.Now())
	}
}

func Benchmark_calculateEventC(b *testing.B) {
	var ipstr [6]string
	ipstr = [6]string{"*/30", "*", "*", "*", "*", "01-12"}
	cr := ParseInput(ipstr)
	for i := 0; i < b.N; i++{
		cr.calculateEvent(time.Now())
	}
}

func Benchmark_calculateEventD(b *testing.B) {
	var ipstr [6]string
	ipstr = [6]string{"*", "*", "*", "*", "*", "*"}
	cr := ParseInput(ipstr)
	for i := 0; i < b.N; i++{
		cr.calculateEvent(time.Now())
	}
}

func Benchmark_calculateEventE(b *testing.B) {
	var ipstr [6]string
	ipstr = [6]string{"*", "*", "*", "02", "*", "01-04"}
	cr := ParseInput(ipstr)
	for i := 0; i < b.N; i++{
		cr.calculateEvent(time.Now())
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

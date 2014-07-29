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
//	"fmt"
	"strconv"
	"strings"
	"time"
	"sort"
)

//Reads the input and returns a pointer to the generated datastructure.
// input[0]: Second
// input[1]: Minute
// input[2]: Hour
// input[3]: Weekday
// input[4]: Day (of Month)
// input[5]: Month
func ParseInput(input [6]string) *crontime {
	var result crontime
	// Check the syntax of the input
	for i := 0; i != len(input); i++ {
		if check_syntax(input[i]) == false {
			panic("Invalid Config") // TODO be more helpful
		}
	}
	// Parse Input like 23-05
	for i := 0; i != len(input); i++ {
		if strings.Contains(input[i], "-") {
			result.parseRange(input[i], i)
		}
	}
	// Parse Input like 05,23,17
	for i := 0; i != len(input); i++ {
		if strings.Contains(input[i], ",") {
			result.parseIRange(input[i], i)
		}
	}
	// Parse input like */05
	for i := 0; i != len(input); i++ {
		if strings.Contains(input[i], "*/") {
			result.parsePeriodic(input[i], i)
		}
	}
	// Parse input like *
	for i := 0; i != len(input); i++ {
		if input[i] == "*" {
			result.parseIgnore(i)
		}
	}
	//Parse single Numbers (05 or 23)
	for i := 0; i != len(input); i++ {
		if !strings.ContainsAny(input[i], "-*/,") {
			tempary := make([]int, 1)
			temp, err := strconv.ParseInt(input[i], 10, 0)
			tempary[0] = int(temp)
			if err != nil {
				panic(err)
			}
			switch i {
			case 0:
				result.second = tempary
			case 1:
				result.minute = tempary
			case 2:
				result.hour = tempary
			case 3:
				result.dow = tempary
			case 4:
				result.dom = tempary
			case 5:
				result.month = tempary
			}
		}
	}
	// Do a sanity check, eg. there is no 32th Day in any Month
	if !check_values(result.second, 0) { panic("rndpanic") }
	if !check_values(result.minute, 1) { panic("rndpanic") }
	if !check_values(result.hour, 2) { panic("rndpanic") }
	if !check_values(result.dow, 3) { panic("rndpanic") }
	if !check_values(result.dom, 4) { panic("rndpanic") }
	if !check_values(result.month, 5) { panic("rndpanic") }
	// Makes timestamp generation easier
	sort.Ints(result.second)
	sort.Ints(result.minute)
	sort.Ints(result.hour)
	sort.Ints(result.dow)
	sort.Ints(result.dom)
	sort.Ints(result.month)

	return &result
}

func (c *crontime) parsePeriodic(input string, i int) {
	temp64, err := strconv.ParseInt(strings.Split(input, "*/")[1], 10, 0)
	if err != nil {
		panic(err)
	}
	temp := int(temp64)
	switch i {
	case 0:
		cur := int(time.Now().Second())
		c.second = periodic(temp, cur, 60)
	case 1:
		cur := int(time.Now().Minute())
		c.minute = periodic(temp, cur, 60)
	case 2:
		cur := int(time.Now().Hour())
		c.hour = periodic(temp, cur, 24)
	case 3:
		cur := int(time.Now().Weekday())
		c.dow = periodic(temp, cur, 7)
	case 4:
		cur := int(time.Now().Day())
		c.dom = periodic(temp, cur, 31)
	case 5:
		cur := int(time.Now().Month())
		c.dom = periodic(temp, cur, 12)
	}
}

func (c *crontime) parseIRange(input string, i int) {
	tempstr := strings.Split(input, ",")
	tempary64 := make([]int64, len(tempstr))
	tempary := make([]int, len(tempstr))
	var err error
	for j := 0; j != len(tempstr); j++ {
		tempary64[j], err = strconv.ParseInt(tempstr[j], 10, 0)
		tempary[j] = int(tempary64[j])
		if err != nil {
			panic(err)
		}
	}
	switch i {
	case 0:
		c.second = tempary
	case 1:
		c.minute = tempary
	case 2:
		c.hour = tempary
	case 3:
		c.dow = tempary
	case 4:
		c.dom = tempary
	case 5:
		c.month = tempary
	}
}

func (c *crontime) parseRange(input string, i int) {
	tempstr := strings.Split(input, "-")
	a64, aerr := strconv.ParseInt(tempstr[0], 10, 0)
	b64, berr := strconv.ParseInt(tempstr[1], 10, 0)
	a := int(a64)
	b := int(b64)
	if aerr != nil {
		panic(aerr)
	}
	if berr != nil {
		panic(berr)
	}
	switch i {
	case 0:
		c.second = value_range(a, b, 60)
	case 1:
		c.minute = value_range(a, b, 60)
	case 2:
		c.hour = value_range(a, b, 24)
	case 3:
		c.dow = value_range(a, b, 7)
	case 4:
		c.dom = value_range(a, b, 32)
	case 5:
		c.month = value_range(a, b, 13)
	}
}

func (c *crontime) parseIgnore(i int) {
	switch i {
	case 0:
		c.second = make([]int, 60)
		for j := 0; j != len(c.second); j++ {
			c.second[j] = j
		}
	case 1:
		c.minute = make([]int, 60)
		for j := 0; j != len(c.minute); j++ {
			c.minute[j] = j
		}
	case 2:
		c.hour = make([]int, 24)
		for j := 0; j != len(c.hour); j++ {
			c.hour[j] = j
		}
	case 3:
		c.dow = make([]int, 7)
		for j := 0; j != len(c.dow); j++ {
			c.dow[j] = j
		}
	case 4:
		c.dom = make([]int, 31)
		for j := 0; j != len(c.dom); j++ {
			c.dom[j] = j + 1
		}
	case 5:
		c.month = make([]int, 12)
		for j := 0; j != len(c.month); j++ {
			c.month[j] = j + 1
		}
	}
}

func check_syntax(insane string) bool {
	/* State machine; Sane == S3, S1
	~	0\-9		,		-		*		/
	S0	S + 2		Err		Err		S + 3	Err
	S1	Err			Err		Err		Err		Err
	S2	S++			Err		Err		Err		Err
	S3	Err			S - 3	S + 2	Err		S + 2
	S4	S -	3		Err		Err		Err		Err
	S5	S--			Err		Err		Err		Err
	*/

	state := 0
	for i := 0; i != len(insane); i++ {
		if strings.ContainsAny(string(insane[i]), "0123456789") {
			if state == 0 {
				state += 2
			} else if state == 2 {
				state++
			} else if state == 4 {
				state -= 3
			} else if state == 5 {
				state--
			} else {
				return false
			}
		} else if strings.ContainsAny(string(insane[i]), ",") {
			if state == 3 {
				state -= 3
			} else {
				return false
			}
		} else if strings.ContainsAny(string(insane[i]), "-") {
			if state == 3 {
				state += 2
			} else {
				return false
			}
		} else if strings.ContainsAny(string(insane[i]), "*") {
			if state == 0 {
				state += 3
			} else {
				return false
			}
		} else if strings.ContainsAny(string(insane[i]), "/") {
			if state == 3 {
				state += 2
			} else {
				return false
			}
		} else {
			return false
		}
	}
	if state == 3 || state == 1 {
		return true
	} else {
		return false
	}
}

func check_values(a []int, i int) bool {
	for j := 0; j != len(a); j++{
		switch i{
			case 0: if a[j] < 0 || a[j] > 59 { return false }
			case 1: if a[j] < 0 || a[j] > 59 { return false }
			case 2: if a[j] < 0 || a[j] > 23 { return false }
			case 3: if a[j] < 0 || a[j] > 6 { return false }
			case 4: if a[j] < 1 || a[j] > 31 { return false }
			case 5: if a[j] < 1 || a[j] > 12 { return false }
		}
	}
	return true
}

// Add two values and ignore the carry (for calculations in the Sexagesimal 
// system).
func add(a, b, bp int) int {
	return (a + b) % bp
}

// 'Absolute' value (or distance) with an variable base.
// Example: absolute_over_breakpoint(59, 1, 60) returns 2
func absolute_over_breakpoint(a, b, base int) int {
	if a >= base || b >= base {
		panic("Invalid Values")
	}
	if a < b {
		return b - a
	} else if a == b {
		return 0
	} else {
		return (base - a) + b
	}
}

// Returns an array filled with all values between a and b considering
// the base.
func value_range(a, b, base int) []int {
	value := make([]int, absolute_over_breakpoint(a, b, base)+1)
	i := 0
	for ; a != b; a++ {
		if a == base {
			a = 0
		}
		value[i] = a
		i++
	}
	value[i] = a
	return value
}

func periodic(a, cur, bp int) []int {
	ret := make([]int, 60)
	i := 2
	ret[0] = cur
	ret[1] = add(a, cur, bp)
	for ;; i++ {
		ret[i] = add(a, ret[i - 1], bp)
		if ret[i] == cur{
			break
		}
	}
	return ret[0:i]
}

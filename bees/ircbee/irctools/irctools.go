// Package irctools is a collection of convenient IRC styling methods.
package irctools

/*
func PostTobeehive(host string, channel string, val string) {
	addr, err := net.ResolveTCPAddr("tcp", host)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	conn, err := net.DialTCP("tcp", nil, addr)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer conn.Close()

	if len(channel) > 0 {
		val = channel + " " + val
	}
	_, err = conn.Write([]byte(val))
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	_, err = conn.Write([]byte("\n"))
}
*/

// Bold wraps val in bold styling tags.
func Bold(val string) string {
	return "\x02" + val + "\x02"
}

// Colored wraps val in color styling tags.
func Colored(val string, color string) string {
	// 00 white 01 black 02 blue (navy) 03 green 04 red 05 brown (maroon)
	// 06 purple 07 orange (olive) 08 yellow 09 light green (lime)
	// 10 teal (a green/blue cyan) 11 light cyan (cyan) (aqua) 12 light blue (royal)
	// 13 pink (light purple) (fuchsia) 14 grey 15 light grey (silver)
	c := "01"
	switch color {
	case "white":
		c = "00"
	case "black":
		c = "01"
	case "blue":
		c = "02"
	case "green":
		c = "03"
	case "red":
		c = "04"
	case "brown":
		c = "05"
	case "purple":
		c = "06"
	case "orange":
		c = "07"
	case "yellow":
		c = "08"
	case "lime":
		c = "09"
	case "teal":
		c = "10"
	case "cyan":
		c = "11"
	case "lightblue":
		c = "12"
	case "pink":
		c = "13"
	case "grey":
		c = "14"
	case "silver":
		c = "15"
	}

	return "\x03" + c + val + "\x03"
}

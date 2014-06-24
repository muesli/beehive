// beehive's application container. Handles command-line arguments
// parsing.
package app

import (
	"flag"
)

type CliFlag struct {
	V     interface{}
	Name  string
	Value interface{}
	Desc  string
}

var (
	appflags []CliFlag
)

func AddFlags(flags []CliFlag) {
	for _, flag := range flags {
		appflags = append(appflags, flag)
	}
}

func Run() {
	for _, f := range appflags {
		switch f.Value.(type) {
		case string:
			flag.StringVar((f.V).(*string), f.Name, f.Value.(string), f.Desc)
		case bool:
			flag.BoolVar((f.V).(*bool), f.Name, f.Value.(bool), f.Desc)
		}
	}

	flag.Parse()
}

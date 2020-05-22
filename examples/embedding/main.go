package main

import (
	"github.com/muesli/beehive/bees"
	_ "github.com/muesli/beehive/bees/hellobee"
	_ "github.com/muesli/beehive/bees/timebee"
	"github.com/muesli/beehive/cfg"
	"github.com/muesli/beehive/reactor"
)

func main() {
	config, err := cfg.New("test.cfg")
	// We could also use the memory backend
	// config, err := cfg.New("mem://")
	if err != nil {
		panic(err)
	}

	config.Bees = []bees.BeeConfig{
		newTimeBee(),
		newHelloBee(),
	}

	// Create the Action and add it to the config
	action := bees.Action{}
	action.ID = "123"
	action.Bee = "hello" // this is the name we gave to the bee in newHelloBee, not the bee ID
	action.Name = "say_hello"
	action.Options = bees.Placeholders{}
	config.Actions = []bees.Action{action}

	// Create the event t
	event := bees.Event{}
	event.Name = "time"
	event.Bee = "timer"

	// Create the chain and add it to the config
	chain := bees.Chain{}
	chain.Name = "foochain"
	chain.Description = "this is a test chain that will say hello every second"
	chain.Actions = []string{"123"} // Action ID we create above
	chain.Event = &event
	chain.Filters = []string{}
	config.Chains = []bees.Chain{chain}

	// Debugging level, prints debug messages from bees
	reactor.SetLogLevel(5)
	reactor.Run(config)

	// Optional, saves the config to disk (if we didn't use the mem backend)
	config.Save()
}

// Create a new bee that says hello
func newHelloBee() bees.BeeConfig {
	options := bees.BeeOptions{}
	bc, err := bees.NewBeeConfig("hello", "hellobee", "test", options)
	if err != nil {
		panic(err)
	}

	return bc
}

// Create a new bee that triggers events every second
func newTimeBee() bees.BeeConfig {
	options := bees.BeeOptions{
		bees.BeeOption{Name: "second", Value: "-1"},
		bees.BeeOption{Name: "minute", Value: "-1"},
		bees.BeeOption{Name: "hour", Value: "-1"},
		bees.BeeOption{Name: "day_of_week", Value: "-1"},
		bees.BeeOption{Name: "day_of_month", Value: "-1"},
		bees.BeeOption{Name: "month", Value: "-1"},
		bees.BeeOption{Name: "year", Value: "-1"},
	}
	bc, err := bees.NewBeeConfig("timer", "timebee", "test", options)
	if err != nil {
		panic(err)
	}

	return bc
}

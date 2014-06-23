package main

import (
	"fmt"
	_ "log"

	"github.com/muesli/beehive/app"

	/* "github.com/muesli/beehive/commands"
	_ "github.com/muesli/beehive/commands/send" */

	"github.com/muesli/beehive/modules"
	_ "github.com/muesli/beehive/modules/ircbee"
)

func main() {
	// Parse command-line args for all registered sub modules
	app.Run()

	fmt.Println("Beehive is buzzing...")

	// Initialize commands and messaging sub-systems
	// commands.StartCommands()
	modules.StartModules()

	// Keep app alive
	ch := make(chan bool)
	<-ch
}

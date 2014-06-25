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

package reactor

import (
	"os"
	"os/signal"
	"syscall"

	log "github.com/sirupsen/logrus"

	"github.com/muesli/beehive/bees"
	"github.com/muesli/beehive/cfg"
)

// Reactor loops and handles signals.
func Run(config *cfg.Config) {
	// Load actions from config
	bees.SetActions(config.Actions)
	// Load chains from config
	bees.SetChains(config.Chains)
	// Initialize bees
	bees.StartBees(config.Bees)

	// Wait for signals
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt, syscall.SIGHUP, syscall.SIGTERM, syscall.SIGKILL)

	log.Debug("Starting Beehive's reactor")
	for s := range ch {
		log.Println("Got signal:", s)

		abort := false
		switch s {
		case syscall.SIGHUP:
			err := config.Load()
			if err != nil {
				log.Panicf("Error loading config from %s: %v", config.URL(), err)
			}
			bees.StopBees()
			bees.SetActions(config.Actions)
			bees.SetChains(config.Chains)
			bees.StartBees(config.Bees)

		case syscall.SIGTERM:
			fallthrough
		case syscall.SIGKILL:
			fallthrough
		case os.Interrupt:
			abort = true
			break
		}

		if abort {
			break
		}
	}
}

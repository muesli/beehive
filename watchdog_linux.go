package main

import (
	"net/http"
	"time"

	"github.com/coreos/go-systemd/daemon"
	"github.com/muesli/beehive/api"
	log "github.com/sirupsen/logrus"
)

// Notifies Systemd's watchdog every WatchdogSec/3 seconds when running
// under Systemd and the watchdog feature has been enabled in
// the service unit.
//
// This will no-op when not running under Systemd.
//
// See http://0pointer.de/blog/projects/watchdog.html
// and https://www.freedesktop.org/software/systemd/man/systemd.service.html
//
func init() {
	// returns the configured WatchdogSec in the service unit as time.Duration
	interval, err := daemon.SdWatchdogEnabled(false)
	if err != nil || interval == 0 {
		log.Printf("Systemd watchdog not enabled")
		return
	}

	// We want to notify the watchdog every WatchdogSec/3, that is, if WatchdogSec is
	// set to 30 seconds, we'll send a notification to systemd every 10 seconds.
	runEvery := interval / 3
	log.Printf("Systemd watchdog notifications every %.2f seconds", runEvery.Seconds())

	go func() {
		for {
			select {
			case <-time.After(runEvery):
				resp, err := http.Get(api.CanonicalURL().String())
				if err == nil {
					resp.Body.Close()
					log.Debugf("Systemd watchdog notify")
					daemon.SdNotify(false, daemon.SdNotifyWatchdog)
				}
			}
		}
	}()
}

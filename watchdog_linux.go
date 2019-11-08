package main

import (
	"net/http"
	"time"

	"github.com/coreos/go-systemd/daemon"
	"github.com/muesli/beehive/api"
	log "github.com/sirupsen/logrus"
)

// Notifies Systemd's watchdog every WatchdogSec/3 when running
// under Systemd and the watchdog feature has been enabled in
// the service unit.
//
// This will no-op when not running under Systemd.
//
// See http://0pointer.de/blog/projects/watchdog.html
// and https://www.freedesktop.org/software/systemd/man/systemd.service.html
//
func init() {
	interval, err := daemon.SdWatchdogEnabled(false)
	// systemd's service unit interval in microseconds
	if err != nil || interval == 0 {
		log.Printf("Systemd watchdog not enabled")
		return
	}
	interval = interval / 3 / 1000000000
	log.Printf("Systemd watchdog interval: %d", interval)

	go func() {
		for {
			select {
			case <-time.After(time.Duration(interval) * time.Second):
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

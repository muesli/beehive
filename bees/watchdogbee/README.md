# Systemd's Watchdog bee

Integrate Beehive with the Systemd's watchdog.

See http://0pointer.de/blog/projects/watchdog.html

## Configuration

Note the "Systemd Watchdog" bee needs to be added to Beehive's config before we start Beehive from systemd with a watchdog enabled.
The easiest way is to do it using Beehive's admin web UI.

* Start beehive
* Add a new "Systemd Watchdog" bee (no config required)
* Add the following Systemd's service unit to `/etc/systemd/system`

```
[Unit]
Description=Beehive with Systemd's watchdog

[Service]
Type=simple
ExecStart=/home/rubiojr/beehive/beehive --config /path/to/beehive.conf
Restart=on-failure
WatchdogSec=30s

[Install]
WantedBy=multi-user.target
```
*Note: change `/path/to/beehive.conf` to a real path pointing to Beehive's config*

* Kill Beehive's existing process you started to configure the watchdog bee
* Enable the new service: `systemctl enable beehive`
* Start the service: `systemctl start beehive`

Beehive will now notify Systemd's watchdog every WatchdogSec/3 seconds, 10 seconds in this particular case.
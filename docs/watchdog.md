# Systemd's Watchdog integration

Integrates Beehive with the Systemd's watchdog.

See http://0pointer.de/blog/projects/watchdog.html

## Configuration

As a system service:

* Add the service unit to `/etc/systemd/system/beehive.conf`

Sample service unit:

```
[Unit]
Description=Beehive with Systemd's watchdog
[Service]
Type=simple
ExecStart=/usr/bin/beehive --config /path/to/beehive.conf
Restart=on-failure
WatchdogSec=30s

[Install]
WantedBy=multi-user.target
```
*Note: change `/path/to/beehive.conf` to a real path pointing to Beehive's config*

* Enable the new service: `systemctl enable beehive`
* Start the service: `systemctl start beehive`

Beehive will automatically detect it's running under Systemd and notify Systemd's watchdog every WatchdogSec/3 seconds (10 seconds in this particular case).
#! /bin/sh

### BEGIN INIT INFO
# Provides:             beehive
# Required-Start:       $syslog
# Required-Stop:        $syslog
# Default-Start:        2 3 4 5
# Default-Stop:
# Short-Description:    beehive daemon
### END INIT INFO

set -e

# /etc/init.d/beehive: start and stop the beehive daemon

umask 022

. /lib/lsb/init-functions

export GOROOT="/usr/local/opt/go"
export PATH="$GOROOT/bin:$PATH"
export GOPATH="/home/beehive/go"

BINARY="$GOPATH/bin/beehive"
PIDFILE="/var/run/beehive.pid"
CONFIG="/home/beehive/beehive.conf"

test -x $BINARY || exit 0

case "$1" in
  start)
        log_daemon_msg "Starting beehive daemon" "beehive" || true
        if start-stop-daemon --start -b --quiet --oknodo -m --pidfile $PIDFILE --exec $BINARY -- -config="$CONFIG"; then
            log_end_msg 0 || true
        else
            log_end_msg 1 || true
        fi
        ;;
  stop)
        log_daemon_msg "Stopping beehive daemon" "beehive" || true
        if start-stop-daemon --stop --quiet --oknodo --pidfile $PIDFILE; then
            log_end_msg 0 || true
        else
            log_end_msg 1 || true
        fi
        ;;

  reload|force-reload)
        log_daemon_msg "Reloading beehive daemon's configuration" "beehive" || true
        if start-stop-daemon --stop --signal 1 --quiet --oknodo --pidfile $PIDFILE --exec $BINARY ; then
            log_end_msg 0 || true
        else
            log_end_msg 1 || true
        fi
        ;;

  restart)
        log_daemon_msg "Restarting beehive daemon" "beehive" || true
        start-stop-daemon --stop --quiet --oknodo --retry 30 --pidfile $PIDFILE
        if start-stop-daemon --start -b --quiet --oknodo -m --pidfile $PIDFILE --exec $BINARY -- -config="$CONFIG"; then
            log_end_msg 0 || true
        else
            log_end_msg 1 || true
        fi
        ;;

  status)
        status_of_proc -p $PIDFILE $BINARY beehive && exit 0 || exit $?
        ;;

  *)
        log_action_msg "Usage: /etc/init.d/beehive {start|stop|reload|restart|status}" || true
        exit 1
esac

exit 0

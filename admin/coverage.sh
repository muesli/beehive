#!/bin/sh

go test -coverprofile /tmp/beehive_coverage1.part ./filters/startswith
go test -coverprofile /tmp/beehive_coverage2.part ./filters/endswith
go test -coverprofile /tmp/beehive_coverage3.part ./filters/contains
go test -coverprofile /tmp/beehive_coverage4.part ./filters/equals

go test -coverprofile /tmp/beehive_coverage4.part ./bees/cronbee/cron

echo "mode: set" > /tmp/beehive_coverage.out
grep -h -v "mode: set" /tmp/beehive_coverage*.part >> /tmp/beehive_coverage.out
go tool cover -html=/tmp/beehive_coverage.out

rm /tmp/beehive_coverage.out /tmp/beehive_coverage*.part

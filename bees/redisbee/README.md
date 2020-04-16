# Redis bee

Store arbitrary key/value strings in a Redis server.

## Configuration

* host: Redis host (defaults to localhost)
* port: Redis port (defaults to 6379)
* password: Redis server password (defaults to passwordless auth)
* db: Redis database to use (defaults to 0)
* channel: Redis channel to subscribe to (pubsub disabled if not specified)

## Ideas

* Use it in combination with the ipify hive to store your public IP history
* Store URLs sent via POST to the HTTP server hive
* Store all messages received in a Slack channel using the Slack hive

See [beehive-youtube-dl](https://github.com/rubiojr/beehive-youtube-dl), a sample project that combines multiple hives, including this one.

## Credits

[Redis Icon](https://iconscout.com/icon/redis-4) by Icon Mafia on [Iconscout]](https://iconscout.com)

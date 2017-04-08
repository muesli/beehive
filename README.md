Beehive
=======

Beehive is an event and agent system, which allows you to create your own
agents that perform automated tasks triggered by events and filters. It is
modular, flexible and really easy to extend for anyone. It has modules
(we call them *Hives*), so it can interface with, talk to, or retrieve
information from Twitter, Tumblr, Email, IRC, Jabber, RSS, Jenkins, Hue - to
name just a few. Check out the full list of [available Hives](https://github.com/muesli/beehive/wiki/Available-Hives)
in our Wiki.

Connecting those modules with each other lets you create immensly useful agents.

#### Here are just a few examples of things Beehive could do for you:
* Re-post tweets on your Tumblr blog
* Forward incoming chat messages to your email account
* Turn on the heating system if the temperature drops below a certain value
* Run your own IRC bot that lets you trigger builds on a Jenkins CI
* Control your Hue lighting system
* Notify you when a stock's price drops below a certain value

![beehive's Logo](/assets/logo_256.png?raw=true)

## Build it

Make sure you have a working Go environment. See the [install instructions](http://golang.org/doc/install.html).

To install beehive, simply run:

    go get github.com/muesli/beehive

To compile it from source:

    cd $GOPATH/src/github.com/muesli/beehive
    go submodule update --init
    go get -u
    go build

Run ```beehive --help``` to see a full list of options.

## Run it

### Locally

Start `beehive` and open http://localhost:8181 in your browser. Note that Beehive will create a config file beehive.conf in its current working directory, unless you specify a different file with the `-config` option.

### Remotely (Linux)

You can run beehive on a remote Linux server (Raspberry Pi servers are great!). To do this, we'll need to copy the `assets` and `config` directories from the beehive source directory to the remote server.

1. Create a directory on the remote server to store everything (`/home/myuser/beehive` for example).
2. After following the building instrutions, copy the `beehive` binary, the `assets` and `config` directories to the directory in the remote server you created in the previous step.

Step 2 is important, otherwise beehive won't be able to access the assets and you'll see a blank web page or missing images when trying to configure it.

### Option 1: bind to localhost and forward port via SSH

Login to the remote server and run beehive:

```
remote-server$ cd /home/myuser/beehive
remote-server$ ./beehive --config beehive.conf
```

Access the admin interface remotely forwarding the connections via SSH:

```
ssh -L8181:localhost:8181 -N myuser@my-server-ip
```

Fire up a browser and go to http://localhost:8181

### Option 2: use a proxy server like Caddy to add basic auth (no TLS, insecure)

You can use something like [Caddy](https://caddyserver.com) to easily proxy beehive and add basic authentication.

Login to the remote server and run beehive:

```
remote-server$ cd /home/myuser/beehive
remote-server$ ./beehive --bind 127.0.0.1:8182 --canonicalurl http://my-server:8182 --config beehive.conf
```

Run Caddy with the following Caddy file:

```
# Caddyfile for beehive
:8181
basicauth / admin secret
proxy / localhost:8182
```

```
caddy --conf Caddyfile
```

Fire up a browser and go to http://my-server:8181, it'll ask you for the username and password configured in the Caddyfile (user: admin, password: secret).

If your server has a public IP, you could easily benefit from Caddy's [automatic HTTPS feature](https://caddyserver.com/docs/automatic-https).

NGINX, HAProxy and Apache are other (albeit more complex options).

## Configuration

Think of Hives as little plugins, extending Beehive's abilities with events you
can react on and actions you can execute.

Just as examples, there's a Twitter plugin that can
 - react to someone you follow posting a tweet (an event)
 - post a new tweet for you (an action)
 - ...

or an RSS plugin that lets you
 - monitor RSS feeds and react on new feed items (another event)

or an email plugin that gives you the ability to
 - send emails (another action)

Each Hive lets you spawn one or multiple Bees in it, all working independently
from another. That allows you to create separate plugin instances, e.g. one
email-Bee for your private mail account, and another one for your work email.

### Creating Bees

Sounds complicated? It's not! Just for fun, let's setup Beehive to send us an
email whenever an RSS feed gets updated. Start ```beehive``` and open <http://localhost:8181/>
in your browser. Note that Beehive will create a config file ```beehive.conf```
in its current working directory, unless you specify a different file with the
```-config``` option.

Note: You currently have to start ```beehive``` from within $GOPATH/src/github.com/muesli/beehive
in order for it to find all the resources for the admin interface. Also see the
Troubleshooting & Notes section of this README.

The admin interface will present you with a list of available Hives. We will
need to create two Bees here, one for the RSS feed and one for your email
account.

![New Bees](https://github.com/muesli/beehive-docs/raw/master/screencaps/new_bees.gif)

### Setting up a Chain

Now we will have to create a new Chain, which will wire up the two Bees we just
created. First we pick the Bee & Event we want to react on, then we pick the
Bee we want to execute an Action with. The RSS-Bee's event gives us a whole set
of parameters we can work with: the feed item's title, its links and
description among others. You can manipulate and combine these parameters with
a full templating language at your disposal. For example we can set the email's
content to something like:

```
Title: {{.title}} - Link: {{index .links 0}}
```

Whenever this action gets executed, Beehive will replace ```{{.title}}``` with
the RSS event's ```title``` parameter, which is the title of the feed item it
retrieved. In the same manner ```{{index .links 0}}``` becomes the first URL of
this event's ```links``` array.

![New Chain](https://github.com/muesli/beehive-docs/raw/master/screencaps/new_chain.gif)

That's it. Whenever the RSS-feed gets updated, Beehive will now send you an
email! It's really easy to make various Bees work together seamlessly and do
clever things for you. Try it yourself!

You can find more information on how to configure beehive and examples [in our Wiki](https://github.com/muesli/beehive/wiki/Configuration).

## Troubleshooting & Notes

The web interface and other resources aren't currently embedded in the binary.
Beehive tries to find those files in its current working directory, so it's
currently recommended to start Beehive from within its git repository, if you
plan to use the web interface.

Should you still not be able to reach the web interface, check if the ```config```
directory in the git repository is empty. If that's the case, make sure the
git submodules get initialized by running ```git submodule update --init```.

The web interface does *not* require authentication yet. Beehive currently
accepts all connections from the loopback device *only*.

## Development

Need help? Want to hack on your own Hives? Join us on IRC (irc://freenode.net/#beehive) or [Gitter](https://gitter.im/the_beehive/Lobby). Follow the bees on [Twitter](https://twitter.com/beehive_app)!

API docs can be found [here](http://godoc.org/github.com/muesli/beehive).

[![Build Status](https://secure.travis-ci.org/muesli/beehive.png)](http://travis-ci.org/muesli/beehive)
[![Go ReportCard](http://goreportcard.com/badge/muesli/beehive)](http://goreportcard.com/report/muesli/beehive)

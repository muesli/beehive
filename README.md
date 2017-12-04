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

## Installation

Beehive requires Go 1.7 or higher. Make sure you have a working Go environment. See the [install instructions](http://golang.org/doc/install.html).

### From source

The recommended way is to fetch the sources and run make:

    go get github.com/muesli/beehive
    cd $GOPATH/src/github.com/muesli/beehive
    make

You can build and install the `beehive` binary like other Go binaries out there (`go get -u`)
but you'll need to make sure beehive can find the assets (images, javascript, css, etc).
See the Troubleshooting/Notes section for additional details.

Run `beehive --help` to see a full list of options.

### Deployment Tools
 - [Dockerfile](docker)
 - [Ansible](https://github.com/morbidick/ansible-role-beehive)

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
email whenever an RSS feed gets updated. Start `beehive` and open <http://localhost:8181/>
in your browser. Note that Beehive will create a config file `beehive.conf`
in its current working directory, unless you specify a different file with the
`-config` option.

Note: You currently have to start `beehive` from within $GOPATH/src/github.com/muesli/beehive
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

Whenever this action gets executed, Beehive will replace `{{.title}}` with
the RSS event's `title` parameter, which is the title of the feed item it
retrieved. In the same manner `{{index .links 0}}` becomes the first URL of
this event's `links` array.

![New Chain](https://github.com/muesli/beehive-docs/raw/master/screencaps/new_chain.gif)

That's it. Whenever the RSS-feed gets updated, Beehive will now send you an
email! It's really easy to make various Bees work together seamlessly and do
clever things for you. Try it yourself!

You can find more information on how to configure beehive and examples [in our Wiki](https://github.com/muesli/beehive/wiki/Configuration).

## Troubleshooting & Notes

The web interface and other resources are embedded in the binary by default.
When using `make noembed`, Beehive tries to find those files
in its current working directory, so it's currently recommended to start Beehive from
within its git repository, if you plan to use the web interface.

Should you still not be able to reach the web interface, check if the `config`
directory in the git repository is empty. If that's the case, make sure the
git submodules get initialized by running `git submodule update --init`.

The web interface does *not* require authentication yet. Beehive currently
accepts all connections from the loopback device *only*.

## Development

Need help? Want to hack on your own Hives? Join us on IRC (irc://freenode.net/#beehive) or [Gitter](https://gitter.im/the_beehive/Lobby). Follow the bees on [Twitter](https://twitter.com/beehive_app)!

[![GoDoc](https://godoc.org/github.com/golang/gddo?status.svg)](https://godoc.org/github.com/muesli/beehive)
[![Build Status](https://travis-ci.org/muesli/beehive.svg?branch=master)](https://travis-ci.org/muesli/beehive)
[![Go ReportCard](http://goreportcard.com/badge/muesli/beehive)](http://goreportcard.com/report/muesli/beehive)

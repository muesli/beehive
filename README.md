beehive
=======

beehive is an event and agent system, which allows you to create your own
agents that perform automated tasks triggered by events and filters. It is
modular, flexible and really easy to extend - for anyone. It has modules
(we call them *Bees*), so it can interface with, talk to, or retrieve
information from Twitter, Tumblr, Email, IRC, Jabber, RSS, Jenkins, Hue - to name
just a few. Check out the full list of [available Bees](https://github.com/muesli/beehive/wiki/Available-Bees)
in our Wiki.

Connecting those Bees with each other let's you create immensly useful agents.

#### Here are just a few examples of things Beehive could do for you:
* Re-post tweets on your Tumblr blog
* Forward incoming chat messages to your email account
* Turn on the heating system if the temperature drops below a certain value
* Run your own IRC bot that let's you trigger builds on a Jenkins CI
* Control your Hue lighting system
* Notify you when a stock's price drops below a certain value

![beehive's Logo](/assets/logo_256.png?raw=true)

## Installation

Make sure you have a working Go environment. See the [install instructions](http://golang.org/doc/install.html).

To install beehive, simply run:

    go get github.com/muesli/beehive

To compile it from source:

    cd $GOPATH/src/github.com/muesli/beehive
    go get -u -v
    go build

Run beehive --help to see a full list of options.

## Configuration

Start Beehive and open <http://localhost:8181/> in your browser. Note that Beehive will create a config file named "beehive.conf" in its current working directory, unless you specify a different file with the -config option.

### Creating Bees

As an example, let's wire up Beehive to send us an email whenever an RSS feed gets updated. The admin interface will present you with a list of available Hives. We will need to create two Bees here, one for the RSS feed and one for your email account.

![New Bees](https://github.com/muesli/beehive-docs/raw/master/screencaps/new_bees.gif)

### Setting up a Chain

Now we will have to create a new Chain, which will wire up the two Bees we just created. First we pick the Bee & Event we want to react on, then we pick the Bee we want to execute an Action with. The RSS event provides us a whole set of parameters we can work with: the feed item's title, its links and description among others. You can manipulate and combine these values with a full templating language at your disposal. For example we can set the email's content to something like:

```
Title: {{.title}} - Link: {{index .links 0}}
```

Whenever this action gets executed Beehive will replace ```{{.title}}``` with the RSS event's ```title``` value, which is the title of the feed item it retrieved. In the same manner ```{{index .links 0}}``` becomes the first URL of this event's ```links``` array.

![New Chain](https://github.com/muesli/beehive-docs/raw/master/screencaps/new_chain.gif)

You can find more information on how to configure beehive and examples [in our Wiki](https://github.com/muesli/beehive/wiki/Configuration).

## Development

API docs can be found [here](http://godoc.org/github.com/muesli/beehive).

[![Build Status](https://secure.travis-ci.org/muesli/beehive.png)](http://travis-ci.org/muesli/beehive)
[![Go ReportCard](http://goreportcard.com/badge/muesli/beehive)](http://goreportcard.com/report/muesli/beehive)

beehive
=======

beehive is an event and agent system, which allows you to create your own
agents that perform automated tasks triggered by events and filters. It is
modular, flexible and really easy to extend - for anyone. It has modules
(we call them *bees*), so it can interface with, talk to, or retrieve
information from Twitter, Tumblr, Email, IRC, Jabber, RSS, Jenkins, Hue - to name
just a few. Check out the full list of [available bees](https://github.com/muesli/beehive/wiki/Available-Bees)
in our Wiki.

Connecting those bees with each other let's you create immensly useful agents.

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

Start Beehive and open http://localhost:8181/config/ in your browser. Note that Beehive will create a config file named "beehive.conf" in its current working directory, unless you specify a different file with the -config option.

As an example, let's wire up Beehive to send us an email whenever an RSS feed gets updated. The admin interface will present you with a list of available Hives. We will need to create two bees here, one for the RSS feed and one for your email account.

![New Bees](https://github.com/muesli/beehive-docs/raw/master/screencaps/new_bees.gif)

Now we will have to create a new Chain, which will wire up the two Bees we just created.

## Manual Configuration

The following paragraph is slightly outdated due to recent changes. We strongly recommend using the admin web-interface to configure Beehive.

To be updated:

#### TL;DR
Take a look at a [few chain recipes here](https://github.com/muesli/beehive/tree/master/recipes).
Pick one, edit it to your needs and store it as 'beehive.conf'. beehive looks for this
configuration file in its current working directory. Alternatively, you can specify a different config file using the -config option.

#### Detailed
The configuration file, beehive.conf by default, is a JSON file that consists of two parts: Bees and Chains.

Bees are pieces of Go code that can communicate with some service.

* you could have a Bee that sits in your IRC channel and responds to certain commands.
* you could have a Bee that hooks into your Instagram account and auto-tweets all of your photos.
* you could have a Bee that hooks into cron, and sends an email with a command's status every N hours of the day.

Bees consist of four parts:

1. "Name"        : the name of the Bee that you are creating
2. "Class"       : the Golang class that communicates with the service you want
3. "Description" : a description of what this Bee will do
4. "Options"     : an array of options to send to the Golang class

Example:

    "Bees":[
       {
          "Name":"ircbee_freenode",
          "Class":"ircbee",
          "Description":"ircbee connected to freenode, channel #beehive",
          "Options":[
             {
                "Name":"server",
                "Value":"irc.freenode.net"
             },
             {
                "Name":"nick",
                "Value":"beehive"
             },
             {
                "Name":"channels",
                "Value":["#beehive"]
             }
          ]
       }
    ],

The above definition creates a Bee that can communicate with IRC channels and can respond to IRC events:

* "ircbee_freenode" is the name of this Bee.
* "ircbee" is the name of Golang class that can communicate using the IRC protocol.
* "Options" are the options to pass to the Golang class "ircbee":
    * The "server" option tells "ircbee" to connect to the IRC server at "irc.freenode.net"
    * The "nick" option tells "ircbee" to use the nickname "beehive"
    * The "channel" option tells "ircbee" to join IRC channel "#beehive".

Chains define what your Bees do, your configuration can have one or more Chains.

Chains consist of four parts:

1. "Name"        : the name of your Chain
2. "Description" : a description of what this Chain does
3. "Event"       : the Event to "listen" for
4. "Elements"    : an array of Filters to apply to, and Actions to take when your Event occurs

Example:

    "Chains":[
       {
          "Name": "filter_chain",
          "Description": "A chain containing various filters",
          "Event":{
             "Bee":"ircbee_freenode",
             "Name":"message"
          },
          "Elements":[
             {
                "Filter":{
                   "Name":"contains",
                   "Options":[
                      {
                         "Name":"text",
                         "CaseInsensitive":true,
                         "Value":"muesli"
                      }
                   ]
                }
             },
             {
                "Action":{
                   "Bee":"ircbee_freenode",
                   "Name":"send",
                   "Options":[
                      {
                         "Name":"channel",
                         "Value":"muesli"
                      },
                      {
                         "Name":"text",
                         "Value":"{{.user}} in {{.channel}} said: {{.text}}"
                      }
                   ]
                }
             }
          ]
       }
    ]

The above definition defines a Chain that responds to the "message" event of the ircbee_freenode Bee.
Whenever a "message" event occurs, all of the Filters defined in the Elements array are applied to the data generated by the "message" event.
If all of the Filters return true, the "send" Action is executed. Bees _generate events_ and can _execute actions_.

Each Action specifies its own options, so it would be difficult to cover them all here.
In this case however, the "send" action messages the user _muesli_ with the message "_{{.user}}_ in _{{.channel}}_ said: _{{.text}}_". _{{.user}}_ is replaced by the user who generated the message that contained "muesli", _{{.channel}}_ is replaced by the channel that the user was in when they generated the message, and _{{.text}}_ is replaced by the "message" that was sent.


## Development

API docs can be found [here](http://godoc.org/github.com/muesli/beehive).

[![Build Status](https://secure.travis-ci.org/muesli/beehive.png)](http://travis-ci.org/muesli/beehive)
[![Go ReportCard](http://goreportcard.com/badge/muesli/beehive)](http://goreportcard.com/report/muesli/beehive)

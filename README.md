beehive
=======

beehive is a flexible event system, which allows you to create your own action-
chains and filters. It is modular and easy to extend - for anyone. It can be
your IRC bot, it can activate your heating system or make you a cup of coffee.
Yes, finally!

![beehive's Logo](/assets/logo.png?raw=true)

## Installation

Make sure you have a working Go environment. See the [install instructions](http://golang.org/doc/install.html).

First of all you need to checkout the source code:

    git clone git://github.com/muesli/beehive.git
    cd beehive

Now we need to get the required dependencies:

    go get -v

Let's build beehive:

    go build

Run beehive -help to see a full list of options!

## Configuration

To be written by someone! :-)

Until then you can find a [few chain recipes here](https://github.com/muesli/beehive/tree/master/recipes).
Pick one, edit it to your needs and store it as 'beehive.conf'. beehive looks for this
configuration file in its current working directory.

## Development

API docs can be found [here](http://godoc.org/github.com/muesli/beehive).

Continuous integration: [![Build Status](https://secure.travis-ci.org/muesli/beehive.png)](http://travis-ci.org/muesli/beehive)

beehive
=======

## Installation

Make sure you have a working Go environment. See the [install instructions](http://golang.org/doc/install.html).

First we need to get the required dependencies. beehive itself is part of that
list so the main executable can depend on our sub-packages:

    go get github.com/fluffle/goirc/client
    go get github.com/hoisie/web
    go get github.com/muesli/beehive

Now we can build beehive:

    git clone git://github.com/muesli/beehive.git
    cd beehive
    go build

## Development

API docs can be found [here](http://godoc.org/github.com/muesli/beehive).

Continous integration: [![Build Status](https://secure.travis-ci.org/muesli/beehive.png)](http://travis-ci.org/muesli/beehive)

// beehive's IRC module.
package ircbee

import (
	irc "github.com/fluffle/goirc/client"
	"github.com/muesli/beehive/app"
	"github.com/muesli/beehive/modules"
	"log"
	"strings"
	"time"
)

type IrcBee struct {
	// channel signaling irc connection status
	ConnectedState chan bool

	// setup IRC client:
	client *irc.Conn

	irchost     string
	ircnick     string
	ircpassword string
	ircssl      bool
	ircchannel  string

	channels []string
}

// Interface impl

func (sys *IrcBee) Name() string {
	return "ircbee"
}

func (sys *IrcBee) Description() string {
	return "An IRC module for beehive"
}

func (sys *IrcBee) Events() []modules.Event {
	events := []modules.Event{
		modules.Event{
			Namespace:	 sys.Name(),
			Name:        "message",
			Description: "A message was received over IRC, either in a channel or a private query",
			Options: []modules.Placeholder{
				modules.Placeholder{
					Name:        "text",
					Description: "The message that was received",
					Type:        "string",
				},
				modules.Placeholder{
					Name:        "channel",
					Description: "The channel the message was received in",
					Type:        "string",
				},
				modules.Placeholder{
					Name:        "user",
					Description: "The user that sent the message",
					Type:        "string",
				},
			},
		},
	}
	return events
}

func (sys *IrcBee) Actions() []modules.Action {
	actions := []modules.Action{
		modules.Action{
			Namespace:	 sys.Name(),
			Name:        "send",
			Description: "Sends a message to a channel or a private query",
			Options: []modules.Placeholder{
				modules.Placeholder{
					Name:        "channel",
					Description: "Which channel to send the message to",
					Type:        "string",
				},
				modules.Placeholder{
					Name:        "text",
					Description: "Content of the message",
					Type:        "string",
				},
			},
		},
		modules.Action{
			Namespace:	 sys.Name(),
			Name:        "join",
			Description: "Joins a channel",
			Options: []modules.Placeholder{
				modules.Placeholder{
					Name:        "channel",
					Description: "Channel to join",
					Type:        "string",
				},
			},
		},
		modules.Action{
			Namespace:	 sys.Name(),
			Name:        "part",
			Description: "Parts a channel",
			Options: []modules.Placeholder{
				modules.Placeholder{
					Name:        "channel",
					Description: "Channel to part",
					Type:        "string",
				},
			},
		},
	}
	return actions
}

func (sys *IrcBee) Action(action modules.Action) []modules.Placeholder {
	outs := []modules.Placeholder{}
	tos := []string{}
	text := ""

	switch action.Name {
	case "send":
		for _, opt := range action.Options {
			if opt.Name == "channel" {
				tos = append(tos, opt.Value.(string))
			}
			if opt.Name == "text" {
				text = opt.Value.(string)
			}
		}
	default:
		// unknown action
		return outs
	}

	for _, recv := range tos {
		if recv == "*" {
			// special: send to all joined channels
			for _, to := range sys.channels {
				sys.client.Privmsg(to, text)
			}
		} else {
			// needs stripping hostname when sending to user!host
			if strings.Index(recv, "!") > 0 {
				recv = recv[0:strings.Index(recv, "!")]
			}

			sys.client.Privmsg(recv, text)
		}
	}

	return outs
}

// ircbee specific impl

func (sys *IrcBee) Rejoin() {
	for _, channel := range sys.channels {
		sys.client.Join(channel)
	}
}

func (sys *IrcBee) Join(channel string) {
	channel = strings.TrimSpace(channel)
	sys.client.Join(channel)

	sys.channels = append(sys.channels, channel)
}

func (sys *IrcBee) Part(channel string) {
	channel = strings.TrimSpace(channel)
	sys.client.Part(channel)

	for k, v := range sys.channels {
		if v == channel {
			sys.channels = append(sys.channels[:k], sys.channels[k+1:]...)
			return
		}
	}
}

func (sys *IrcBee) Run(channelIn chan modules.Event) {
	if len(sys.irchost) == 0 {
		return
	}

	// channel signaling irc connection status
	sys.ConnectedState = make(chan bool)

	// setup IRC client:
	sys.client = irc.SimpleClient(sys.ircnick, "beehive", "beehive")
	sys.client.SSL = sys.ircssl

	sys.client.AddHandler(irc.CONNECTED, func(conn *irc.Conn, line *irc.Line) {
		sys.ConnectedState <- true
	})
	sys.client.AddHandler(irc.DISCONNECTED, func(conn *irc.Conn, line *irc.Line) {
		sys.ConnectedState <- false
	})
	sys.client.AddHandler("PRIVMSG", func(conn *irc.Conn, line *irc.Line) {
		channel := line.Args[0]
		if channel == sys.client.Me.Nick {
			channel = line.Src // replies go via PM too.
		}

		ev := modules.Event{
			Namespace: sys.Name(),
			Name: "message",
			Options: []modules.Placeholder{
				modules.Placeholder{
					Name:  "channel",
					Type:  "string",
					Value: channel,
				},
				modules.Placeholder{
					Name:  "user",
					Type:  "string",
					Value: line.Src,
				},
				modules.Placeholder{
					Name:  "text",
					Type:  "string",
					Value: line.Args[1],
				},
			},
		}
		channelIn <- ev
	})

	// loop on IRC dis/connected events
	go func() {
		for {
			log.Println("Connecting to IRC:", sys.irchost)
			err := sys.client.Connect(sys.irchost, sys.ircpassword)
			if err != nil {
				log.Println("Failed to connect to IRC:", sys.irchost)
				log.Println(err)
				continue
			}
			for {
				status := <-sys.ConnectedState
				if status {
					log.Println("Connected to IRC:", sys.irchost)

					if len(sys.channels) == 0 {
						// join default channel
						sys.Join(sys.ircchannel)
					} else {
						// we must have been disconnected, rejoin channels
						sys.Rejoin()
					}
				} else {
					log.Println("Disconnected from IRC:", sys.irchost)
					break
				}
			}
			time.Sleep(5 * time.Second)
		}
	}()
}

func init() {
	irc := IrcBee{}

	app.AddFlags([]app.CliFlag{
		app.CliFlag{&irc.irchost, "irchost", "", "Hostname of IRC server, eg: irc.example.org:6667"},
		app.CliFlag{&irc.ircnick, "ircnick", "beehive", "Nickname to use for IRC"},
		app.CliFlag{&irc.ircpassword, "ircpassword", "", "Password to use to connect to IRC server"},
		app.CliFlag{&irc.ircchannel, "ircchannel", "#beehivetest", "Which channel to join"},
		app.CliFlag{&irc.ircssl, "ircssl", false, "Use SSL for IRC connection"},
	})

	modules.RegisterModule(&irc)
}

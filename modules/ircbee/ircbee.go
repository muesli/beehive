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

func (mod *IrcBee) Name() string {
	return "ircbee"
}

func (mod *IrcBee) Description() string {
	return "An IRC module for beehive"
}

func (mod *IrcBee) Events() []modules.EventDescriptor {
	events := []modules.EventDescriptor{
		modules.EventDescriptor{
			Namespace:   mod.Name(),
			Name:        "message",
			Description: "A message was received over IRC, either in a channel or a private query",
			Options: []modules.PlaceholderDescriptor{
				modules.PlaceholderDescriptor{
					Name:        "text",
					Description: "The message that was received",
					Type:        "string",
				},
				modules.PlaceholderDescriptor{
					Name:        "channel",
					Description: "The channel the message was received in",
					Type:        "string",
				},
				modules.PlaceholderDescriptor{
					Name:        "user",
					Description: "The user that sent the message",
					Type:        "string",
				},
			},
		},
	}
	return events
}

func (mod *IrcBee) Actions() []modules.ActionDescriptor {
	actions := []modules.ActionDescriptor{
		modules.ActionDescriptor{
			Namespace:   mod.Name(),
			Name:        "send",
			Description: "Sends a message to a channel or a private query",
			Options: []modules.PlaceholderDescriptor{
				modules.PlaceholderDescriptor{
					Name:        "channel",
					Description: "Which channel to send the message to",
					Type:        "string",
				},
				modules.PlaceholderDescriptor{
					Name:        "text",
					Description: "Content of the message",
					Type:        "string",
				},
			},
		},
		modules.ActionDescriptor{
			Namespace:   mod.Name(),
			Name:        "join",
			Description: "Joins a channel",
			Options: []modules.PlaceholderDescriptor{
				modules.PlaceholderDescriptor{
					Name:        "channel",
					Description: "Channel to join",
					Type:        "string",
				},
			},
		},
		modules.ActionDescriptor{
			Namespace:   mod.Name(),
			Name:        "part",
			Description: "Parts a channel",
			Options: []modules.PlaceholderDescriptor{
				modules.PlaceholderDescriptor{
					Name:        "channel",
					Description: "Channel to part",
					Type:        "string",
				},
			},
		},
	}
	return actions
}

func (mod *IrcBee) Action(action modules.Action) []modules.Placeholder {
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
			for _, to := range mod.channels {
				mod.client.Privmsg(to, text)
			}
		} else {
			// needs stripping hostname when sending to user!host
			if strings.Index(recv, "!") > 0 {
				recv = recv[0:strings.Index(recv, "!")]
			}

			mod.client.Privmsg(recv, text)
		}
	}

	return outs
}

// ircbee specific impl

func (mod *IrcBee) Rejoin() {
	for _, channel := range mod.channels {
		mod.client.Join(channel)
	}
}

func (mod *IrcBee) Join(channel string) {
	channel = strings.TrimSpace(channel)
	mod.client.Join(channel)

	mod.channels = append(mod.channels, channel)
}

func (mod *IrcBee) Part(channel string) {
	channel = strings.TrimSpace(channel)
	mod.client.Part(channel)

	for k, v := range mod.channels {
		if v == channel {
			mod.channels = append(mod.channels[:k], mod.channels[k+1:]...)
			return
		}
	}
}

func (mod *IrcBee) Run(channelIn chan modules.Event) {
	if len(mod.irchost) == 0 {
		return
	}

	// channel signaling irc connection status
	mod.ConnectedState = make(chan bool)

	// setup IRC client:
	mod.client = irc.SimpleClient(mod.ircnick, "beehive", "beehive")
	mod.client.SSL = mod.ircssl

	mod.client.AddHandler(irc.CONNECTED, func(conn *irc.Conn, line *irc.Line) {
		mod.ConnectedState <- true
	})
	mod.client.AddHandler(irc.DISCONNECTED, func(conn *irc.Conn, line *irc.Line) {
		mod.ConnectedState <- false
	})
	mod.client.AddHandler("PRIVMSG", func(conn *irc.Conn, line *irc.Line) {
		channel := line.Args[0]
		if channel == mod.client.Me.Nick {
			channel = line.Src // replies go via PM too.
		}

		ev := modules.Event{
			Namespace: mod.Name(),
			Name:      "message",
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
			log.Println("Connecting to IRC:", mod.irchost)
			err := mod.client.Connect(mod.irchost, mod.ircpassword)
			if err != nil {
				log.Println("Failed to connect to IRC:", mod.irchost)
				log.Println(err)
				continue
			}
			for {
				status := <-mod.ConnectedState
				if status {
					log.Println("Connected to IRC:", mod.irchost)

					if len(mod.channels) == 0 {
						// join default channel
						mod.Join(mod.ircchannel)
					} else {
						// we must have been disconnected, rejoin channels
						mod.Rejoin()
					}
				} else {
					log.Println("Disconnected from IRC:", mod.irchost)
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

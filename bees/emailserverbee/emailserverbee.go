/*
 *    Copyright (C) 2014-2017 Christian Muehlhaeuser
 *
 *    This program is free software: you can redistribute it and/or modify
 *    it under the terms of the GNU Affero General Public License as published
 *    by the Free Software Foundation, either version 3 of the License, or
 *    (at your option) any later version.
 *
 *    This program is distributed in the hope that it will be useful,
 *    but WITHOUT ANY WARRANTY; without even the implied warranty of
 *    MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 *    GNU Affero General Public License for more details.
 *
 *    You should have received a copy of the GNU Affero General Public License
 *    along with this program.  If not, see <http://www.gnu.org/licenses/>.
 *
 *    Authors:
 *      Matthias Krauser <matthias@krauser.eu>
 */

package emailserverbee

import (
	"github.com/flashmob/go-guerrilla"
	"github.com/flashmob/go-guerrilla/backends"
	"github.com/flashmob/go-guerrilla/mail"

	//"github.com/flashmob/go-guerrilla/response"

	"github.com/muesli/beehive/bees"
)

// EmailServerBee is a Bee that starts an SMTP server and fires events for incoming
// Emails.
type EmailServerBee struct {
	bees.Bee

	address        string
	allowedHosts   []string
	hostname       string
	maxSize        int64
	timeout        int
	maxClients     int
	startTLSOn     bool
	tlsAlwaysOn    bool
	privateKeyFile string
	publicKeyFile  string

	eventChan chan bees.Event
}

func (mod *EmailServerBee) mailProcessor() func() backends.Decorator {
	return func() backends.Decorator {
		return func(p backends.Processor) backends.Processor {
			return backends.ProcessWith(
				func(e *mail.Envelope, task backends.SelectTask) (backends.Result, error) {
					if task == backends.TaskSaveMail {

						recipients := make([]string, len(e.RcptTo))

						for index, rcpt := range e.RcptTo {
							recipients[index] = rcpt.String()
						}

						/*headers := make([]string, len(e.Header))

						i := 0
						for name, value := range e.Header {
							headers[i] = name + ": " + value
							i += 1
						}*/

						// create events and send it to cin
						ev := bees.Event{
							Bee:  mod.Name(),
							Name: "received",
							Options: []bees.Placeholder{
								{
									Name:  "sender",
									Type:  "string",
									Value: e.MailFrom.String(),
								},
								{
									Name:  "remote_ip",
									Type:  "string",
									Value: e.RemoteIP,
								},
								{
									Name:  "recipients",
									Type:  "[]string",
									Value: recipients,
								},
								{
									Name:  "subject",
									Type:  "string",
									Value: e.Subject,
								},
								{
									Name:  "tls",
									Type:  "boolean",
									Value: e.TLS,
								},
								{
									Name:  "headers",
									Type:  "[]string",
									Value: e.Header,
								},
								{
									Name:  "body",
									Type:  "string",
									Value: e.Data.String(),
								},
							},
						}

						mod.eventChan <- ev
					}

					return p.Process(e, task)
				},
			)
		}
	}
}

// Run executes the Bee's event loop.
func (mod *EmailServerBee) Run(cin chan bees.Event) {
	mod.eventChan = cin

	// see https://github.com/flashmob/go-guerrilla/wiki/Using-as-a-package

	cfg := &guerrilla.AppConfig{
		AllowedHosts: mod.allowedHosts,
		// LogLevel controls the lowest level we log.
		// "info", "debug", "error", "panic". Default "info"
		LogLevel: "info",
	}

	tc := guerrilla.ServerTLSConfig{
		StartTLSOn:     mod.startTLSOn,
		AlwaysOn:       mod.tlsAlwaysOn,
		PrivateKeyFile: mod.privateKeyFile,
		PublicKeyFile:  mod.publicKeyFile,
	}
	sc := guerrilla.ServerConfig{
		IsEnabled:       true,
		Hostname:        mod.hostname,
		ListenInterface: mod.address,
		MaxSize:         mod.maxSize,
		Timeout:         mod.timeout,
		MaxClients:      mod.maxClients,
		TLS:             tc,
	}
	cfg.Servers = append(cfg.Servers, sc)

	bcfg := backends.BackendConfig{
		"save_process":       "HeadersParser|Header|Hasher|Beehive",
		"log_received_mails": true,
	}
	cfg.BackendConfig = bcfg

	d := guerrilla.Daemon{Config: cfg}
	d.AddProcessor("Beehive", mod.mailProcessor())

	err := d.Start()
	if err != nil {
		mod.LogFatal("Error starting SMTP-Server", err)
	}

	select {
	case <-mod.SigChan:
		d.Shutdown()
		return
	}
}

// ReloadOptions parses the config options and initializes the Bee.
func (mod *EmailServerBee) ReloadOptions(options bees.BeeOptions) {
	mod.SetOptions(options)

	options.Bind("address", &mod.address)
	options.Bind("allowedHosts", &mod.allowedHosts)
	options.Bind("startTLSOn", &mod.startTLSOn)
	options.Bind("tlsAlwaysOn", &mod.tlsAlwaysOn)
	options.Bind("privateKeyFile", &mod.privateKeyFile)
	options.Bind("publicKeyFile", &mod.publicKeyFile)
	options.Bind("hostname", &mod.hostname)
	options.Bind("maxSize", &mod.maxSize)
	options.Bind("timeout", &mod.timeout)
	options.Bind("maxClient", &mod.maxClients)
}

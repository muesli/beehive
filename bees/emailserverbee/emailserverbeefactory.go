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
	"os"

	"github.com/muesli/beehive/bees"
)

// EmailServerBeeFactory is a factory for EmailServerBees.
type EmailServerBeeFactory struct {
	bees.BeeFactory
}

// New returns a new Bee instance configured with the supplied options.
func (factory *EmailServerBeeFactory) New(name, description string, options bees.BeeOptions) bees.BeeInterface {
	bee := EmailServerBee{
		Bee: bees.NewBee(name, factory.ID(), description, options),
	}
	bee.ReloadOptions(options)

	return &bee
}

// ID returns the ID of this Bee.
func (factory *EmailServerBeeFactory) ID() string {
	return "emailserverbee"
}

// Name returns the name of this Bee.
func (factory *EmailServerBeeFactory) Name() string {
	return "Email Server"
}

// Description returns the description of this Bee.
func (factory *EmailServerBeeFactory) Description() string {
	return "Receives incoming Emails via SMTP(S) and reacts to them"
}

// Image returns the filename of an image for this Bee.
func (factory *EmailServerBeeFactory) Image() string {
	return factory.ID() + ".png"
}

// LogoColor returns the preferred logo background color (used by the admin interface).
func (factory *EmailServerBeeFactory) LogoColor() string {
	return "#00bbff"
}

// Options returns the options available to configure this Bee.
func (factory *EmailServerBeeFactory) Options() []bees.BeeOptionDescriptor {
	hostname, _ := os.Hostname()

	opts := []bees.BeeOptionDescriptor{
		{
			Name:        "address",
			Description: "Which addr to listen on, eg: 0.0.0.0:25",
			Type:        "address",
			Mandatory:   true,
		},
		{
			Name:        "allowedHosts",
			Description: "AllowedHosts lists which hosts to accept email for.",
			Type:        "[]string",
			Default:     hostname,
			Mandatory:   false,
		},
		{
			Name:        "hostname",
			Description: "Hostname will be used in the server's reply to HELO/EHLO. If TLS enabled make sure that the Hostname matches the cert.",
			Default:     hostname,
			Type:        "string",
			Mandatory:   false,
		},
		{
			Name:        "maxSize",
			Description: "MaxSize is the maximum size of an email that will be accepted for delivery.",
			Default:     1024 * 1024 * 10,
			Type:        "int64",
			Mandatory:   false,
		},
		{
			Name:        "timeout",
			Description: "Timeout specifies the connection timeout in seconds. Defaults to 30",
			Default:     30,
			Type:        "int",
			Mandatory:   false,
		},
		{
			Name:        "maxClients",
			Description: "MaxClients controls how many maxiumum clients we can handle at once.",
			Default:     10,
			Type:        "int",
			Mandatory:   false,
		},
		{
			Name:        "startTLSOn",
			Description: "StartTLSOn should we offer STARTTLS command. Cert must be valid.",
			Default:     false,
			Type:        "bool",
			Mandatory:   false,
		},
		{
			Name:        "tlsAlwaysOn",
			Description: "TLSAlwaysOn run this server as a pure TLS server, i.e. SMTPS",
			Default:     false,
			Type:        "bool",
			Mandatory:   false,
		},
		{
			Name:        "privateKeyFile",
			Description: "PrivateKeyFile path to cert private key in PEM format",
			Default:     "",
			Type:        "string",
			Mandatory:   false,
		},
		{
			Name:        "publicKeyFile",
			Description: "PublicKeyFile path to cert (public key) chain in PEM format.",
			Default:     "",
			Type:        "string",
			Mandatory:   false,
		},
	}
	return opts
}

// Events describes the available events provided by this Bee.
func (factory *EmailServerBeeFactory) Events() []bees.EventDescriptor {
	events := []bees.EventDescriptor{
		{
			Namespace:   factory.Name(),
			Name:        "received",
			Description: "An Email was received",
			Options: []bees.PlaceholderDescriptor{
				{
					Name:        "sender",
					Description: "The Email sender",
					Type:        "string",
				},
				{
					Name:        "remote_ip",
					Description: "The IP from the remote-server",
					Type:        "string",
				},
				{
					Name:        "recipients",
					Description: "The recipients",
					Type:        "[]string",
				},
				{
					Name:        "subject",
					Description: "The Email subject",
					Type:        "string",
				},
				{
					Name:        "tls",
					Description: "Indicates, if the mail was transmitted using tls",
					Type:        "bool",
				},
				{
					Name:        "headers",
					Description: "The Email headers",
					Type:        "[]string",
				},
				{
					Name:        "body",
					Description: "The raw Email body",
					Type:        "string",
				},
			},
		},
	}
	return events
}

func init() {
	f := EmailServerBeeFactory{}
	bees.RegisterFactory(&f)
}

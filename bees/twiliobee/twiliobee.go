// Package twiliobee is a Bee that is able to send SMS messages.
package twiliobee

import (
	"github.com/muesli/beehive/bees"
	twilio "github.com/carlosdp/twiliogo"
)

// TwilioBee is a Bee that is able to send SMS messages.
type TwilioBee struct {
	bees.Bee

	account_sid string
	auth_token string
	from_number string
  to_number string
}

// Action triggers the action passed to it.
func (mod *TwilioBee) Action(action bees.Action) []bees.Placeholder {
	outs := []bees.Placeholder{}

	switch action.Name {
	case "send":
		body := ""

		action.Options.Bind("body", &body)

    client := twilio.NewClient(mod.account_sid, mod.auth_token)
    twilio.NewMessage(client, mod.from_number, mod.to_number, twilio.Body(body))

	default:
		panic("Unknown action triggered in " + mod.Name() + ": " + action.Name)
	}

	return outs
}

// ReloadOptions parses the config options and initializes the Bee.
func (mod *TwilioBee) ReloadOptions(options bees.BeeOptions) {
	mod.SetOptions(options)

	options.Bind("account_sid", &mod.account_sid)
	options.Bind("auth_token", &mod.auth_token)
	options.Bind("from_number", &mod.from_number)
	options.Bind("to_number", &mod.to_number)
}

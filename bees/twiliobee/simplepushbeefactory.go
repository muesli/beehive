package twiliobee

import (
	"github.com/muesli/beehive/bees"
)

// TwilioBeeFactory is a factory for TwilioBees.
type TwilioBeeFactory struct {
	bees.BeeFactory
}

// New returns a new Bee instance configured with the supplied options.
func (factory *TwilioBeeFactory) New(name, description string, options bees.BeeOptions) bees.BeeInterface {
	bee := TwilioBee{
		Bee: bees.NewBee(name, factory.ID(), description, options),
	}
	bee.ReloadOptions(options)

	return &bee
}

// ID returns the ID of this Bee.
func (factory *TwilioBeeFactory) ID() string {
	return "twiliobee"
}

// Name returns the name of this Bee.
func (factory *TwilioBeeFactory) Name() string {
	return "Twilio"
}

// Description returns the description of this Bee.
func (factory *TwilioBeeFactory) Description() string {
	return "Lets you send SMS messages to your phone"
}

// Image returns the filename of an image for this Bee.
func (factory *TwilioBeeFactory) Image() string {
	return factory.ID() + ".png"
}

// LogoColor returns the preferred logo background color (used by the admin interface).
func (factory *TwilioBeeFactory) LogoColor() string {
	return "#ee3248"
}

// Options returns the options available to configure this Bee.
func (factory *TwilioBeeFactory) Options() []bees.BeeOptionDescriptor {
	opts := []bees.BeeOptionDescriptor{
		{
			Name:        "key",
			Description: "Twilio key which you get after installing the app",
			Type:        "string",
			Mandatory:   true,
		},
		{
			Name:        "password",
			Description: "Password for end-to-end encryption (optional)",
			Type:        "string",
			Mandatory:   false,
		},
		{
			Name:        "salt",
			Description: "Salt for end-to-end encryption (optional)",
			Type:        "url",
			Mandatory:   false,
		},
	}
	return opts
}

// Events describes the available events provided by this Bee.
func (factory *TwilioBeeFactory) Events() []bees.EventDescriptor {
	events := []bees.EventDescriptor{}
	return events
}

// Actions describes the available actions provided by this Bee.
func (factory *TwilioBeeFactory) Actions() []bees.ActionDescriptor {
	actions := []bees.ActionDescriptor{
		{
			Namespace:   factory.Name(),
			Name:        "send",
			Description: "Sends a push notification",
			Options: []bees.PlaceholderDescriptor{
				{
					Name:        "title",
					Description: "Title of push notification (optional)",
					Type:        "string",
					Mandatory:   false,
				},
				{
					Name:        "message",
					Description: "Content of push notification",
					Type:        "string",
					Mandatory:   true,
				},
				{
					Name:        "event",
					Description: "Event id for customizing vibration and ringtone (optional)",
					Type:        "string",
					Mandatory:   false,
				},
			},
		},
	}
	return actions
}

func init() {
	f := TwilioBeeFactory{}
	bees.RegisterFactory(&f)
}

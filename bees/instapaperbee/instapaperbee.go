package instapaperbee

import (
	"github.com/muesli/beehive/bees"
	"net/url"
	"net/http"
)

type InstapaperBee struct {
	bees.Bee

	username string
	password string
}

func (mod *InstapaperBee) Run(cin chan bees.Event) {
	select {
	case <-mod.SigChan:
		return
	}
}

func (mod *InstapaperBee) Action(action bees.Action) []bees.Placeholder {
	switch action.Name {
	case "save":
		var title, page_url string
		action.Options.Bind("title", &title)
		action.Options.Bind("url", &page_url)

		msg := url.Values{}
		msg.Set("username", mod.username)
		msg.Set("password", mod.password)
		msg.Set("url", page_url)

		if title != "" {
			msg.Set("title", title)
		}
		mod.LogDebugf("Message: %s", msg)
		resp, err := http.PostForm("https://www.instapaper.com/api/add", msg)
		if err != nil {
			panic(err)
		}
		defer resp.Body.Close()
		if resp.StatusCode == 200 {
			mod.LogDebugf("Added article to instapaper.")
		}
	default:
		panic("Unknown action triggered in " + mod.Name() + ": " + action.Name)
	}
	return []bees.Placeholder{}
}

func (mod *InstapaperBee) ReloadOptions(options bees.BeeOptions) {
	mod.SetOptions(options)
	options.Bind("username", &mod.username)
	options.Bind("password", &mod.password)
}

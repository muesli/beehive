# Mastodonbee

You'll need to register and authorize the 'beehive' application in mastodon in order to use the mastodonbee.


You can use following code snippet:

```
package main

import (
	"context"
	"fmt"
	"log"

	"github.com/mattn/go-mastodon"
)

func main() {
	app, err := mastodon.RegisterApp(context.Background(), &mastodon.AppConfig{
		Server:     "https://server.address", // Enter _your_ server address here
		ClientName: "", // Name for the client
		Scopes:     "read write follow", // Permission scopes
		Website:    "", // Optional
	})
	if err != nil {
		log.Fatal(err)
	}
    // Prints out the ClientID and ClientSecret which you'll need to configure
    // beehive propberly
	fmt.Printf("client-id    : %s\n", app.ClientID)
	fmt.Printf("client-secret: %s\n", app.ClientSecret)
}
```
# GotifyBee

This bee can push notification to the specified gotify application as an action to some event.

You can get your **TOKEN** for your specified gotify application by going to the **APPS** section of the gotify server instance
and unhiding the token for your desired application.

## APP

* [Android APP](https://play.google.com/store/apps/details?id=com.github.gotify&hl=en_US&gl=US)
* [Chrome extension](https://chrome.google.com/webstore/detail/gotify-push/cbegkpikakpajcaoblfkeindhhikpfmd?hl=en)

## Configuration

The **message** field is required. If the message's title is empty, it would be replaced by Gotify.

The priority of the message is option.

### Options
```json
"Bees": [
  {
    "Name": "gotify example",
    "Class": "gotify",
    "Description": "This is example of gotify",
    "Options": [
      {
        "Name": "token",
        "Value": "TOKEN"
      }
    ]
  },
]
```

### Actions

```json
"Actions": [
    {
      "Bee": "gotify example",
      "Name": "send",
      "Options": [
        {
          "Name": "title",
          "Type": "string",
          "Value": ""
        },
        {
          "Name": "message",
          "Type": "string",
          "Value": ""
        },
        {
          "Name": "priority",
          "Type": "string",
          "Value": ""
        }
      ]
    },
]
```
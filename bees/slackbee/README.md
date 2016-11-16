# Slack bee

The [Slack](https://slack.com) bee can send and listen to messages in a Slack channel.

## Configuration

### Options

```json
"Bees":[
   {
      "Name":"slackmsg",
      "Class":"slackbee",
      "Description":"slackbee",
      "Options":[
         {
            "Name":"apiKey",
            "Value":"env://SLACK_KEY"
         },
         {
            "Name":"channels",
            "Value":["rubiojr-test"]
         }
      ]
   }
]```

**apiKey**: Slack API Key. You can get one from https://api.slack.com/docs/oauth-test-tokens.

The API key can be added to the recipe/config as-is, via environment variable (`env://MY_API_KEY`) or read from a file (`file:///Users/rubiojr/.slack_key`).

**channels**: The slack channels to listen on.

### Actions

**send**: send a message to a Slack channel. Needs the name of the channel (not the channel ID), and the text to send. You can use interpolation to send something from the event received:

```json
"Elements":[
   {
      "Action":{
         "Bee":"slackmsg",
         "Name":"send",
         "Options":[
            {
               "Name":"channel",
               "Value":"rubiojr-test2"
            },
            {
               "Name":"text",
               "Value":"{{.something}}"
            }
         ]
      }
   }
]
```

## Credits

Slack logo: https://remoteworkspain.slack.com/brand-guidelines

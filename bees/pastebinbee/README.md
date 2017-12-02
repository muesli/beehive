# Pastebin bee

This bee can post some text on pastebin as an action to some event.

You can get your **api developer key** from: <http://pastebin.com/api>

## Configuration

Due to admin interface limitations here are available options for exposure and expiration of the paste:

- exposure:

  - 0 - Public
  - 1 - Unlisted
  - 2 - Private

- expire:

  - N - Never
  - 10M - 10 minutes
  - 1H - 1 hour
  - 1D - 1 day
  - 1W - 1 week
  - 2W - 2 weeks
  - 1M - 1 month

### Options

```json
"Bees": [
  {
    "Name": "Pastebin example",
    "Class": "pastebinbee",
    "Description": "This is example of pastebinbee",
    "Options": [
      {
        "Name": "api_dev_key",
        "Value": "API_DEVELOPER_KEY"
      }
    ]
  },
]
```

### Actions

```json
"Actions": [
{
  "Bee": "Paste",
  "Name": "post",
  "Options": [
    {
      "Name": "title",
      "Type": "string",
      "Value": "beehive"
    },
    {
      "Name": "content",
      "Type": "string",
      "Value": "testing beehive passed!"
    },
    {
      "Name": "expire",
      "Type": "string",
      "Value": "1H"
    },
    {
      "Name": "exposure",
      "Type": "string",
      "Value": "2"
    }
  ]
},
]
```

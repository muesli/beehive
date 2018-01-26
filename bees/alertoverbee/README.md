# AlertOver Bee

This bee can push notification to your iOS or Android device or extension in
Chrome as an action to some event.

You can get your **SOURCE-KEY** and **RECEIVER-KEY** from:
<https://www.alertover.com/pages/api>

## APP

- [iOS APP](https://itunes.apple.com/cn/app/alertover-gao-xiao-mian-fei/id1069760182?l=en&mt=8)
- [Android APP](http://www.wandoujia.com/apps/com.alertover.app)
- [Chrome extension](https://chrome.google.com/webstore/detail/alertover/cgcgodonijnlgljfdbiicdccnldpdgia?hl=zh-CN)

## Configuration

The SOURCE-KEY, RECEIVER-KEY and message content are mandatory. If the
message's title is empty, it would be replaced by AlertOver.

The priority of the message is option, 0 for normal and 1 for emergency, 0 and
1 is string.

### Options

```json
"Bees": [
  {
    "Name": "alertover example",
    "Class": "alertoverbee",
    "Description": "This is example of alertoverbee",
    "Options": [
      {
        "Name": "source",
        "Value": "SOURCE-KEY"
      }
    ]
  },
]
```

### Actions

```json
"Actions": [
    {
      "Bee": "alertover example",
      "Name": "send",
      "Options": [
        {
          "Name": "receiver",
          "Type": "string",
          "Value": "RECEIVER-KEY"
        },
        {
          "Name": "title",
          "Type": "string",
          "Value": ""
        },
        {
          "Name": "content",
          "Type": "string",
          "Value": ""
        },
        {
          "Name": "url",
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

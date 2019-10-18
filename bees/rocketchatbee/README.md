# Rocket.chat bee

The [Rocket.Chat](https://rocket.chat/) bee can send messages into a Rocket.Chat channel.

## Configuration

### Options

```json
"Bees":[
   {
      "Name":"rocketchatmsg",
      "Class":"rocketchatbee",
      "Description":"my rocket.chat bee",
      "Options":[
         {
            "Name":"url",
            "Value":"http://localhost:3000"
         },
         {
            "Name":"user_id",
            "Value":"YOUR_USER_ID"
         },
         {
            "Name":"auth_token",
            "Value":"YOUR_AUTH_TOKEN"
         }
      ]
   }
]
```

See https://rocket.chat/docs/developer-guides/rest-api/personal-access-tokens/ for reference on how to create a `user_id` and `auth_token`.

### Actions

**send**: send a message to a Rocket.Chat channel. This needs the name of the `channel`, and the `text` to send. If you specify an `alias`, messages posted appear under that username.

```json
"Elements":[
   {
      "Action":{
         "Bee":"rocketchatmsg",
         "Name":"send",
         "Options":[
            {
               "Name":"channel",
               "Value":"info"
            },
            {
               "Name":"text",
               "Value":"This is the latest info!"
            },
            {
               "Name":"alias",
               "Value":"Informer"
            }
         ]
         ]
      }
   }
]
```

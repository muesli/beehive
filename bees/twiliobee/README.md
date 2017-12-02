# Twilio bee

The [Twilio](https://twilio.com) bee can send SMS messages to a phone.

## Configuration

### Options

```json
"Bees": [
  {
    "Name":"Twilio Example Bee",
    "Class":"twiliobee",
    "Description":"A bee to demonstrate Twilio with Beehive",
    "Options":[
      {
        "Name": "account_sid",
        "Value": "YOUR_TWILIO_ACCOUNT_SID"
      },
      {
        "Name": "auth_token",
        "Value": "YOUR_TWILIO_AUTH_TOKEN"
      },
      {
        "Name": "from_number",
        "Value": "+15551234567"
      },
      {
        "Name": "to_number",
        "Value": "+15559876543"
      }
    ]
  }
]
```

**account_sid** and **auth_token**: Twilio Account SID and Authentication Token. You can sign up and get them from <https://www.twilio.com/try-twilio>.

These can be added to the recipe/config as-is (`XXXXXXXX`), via environment variable (`env://MY_ACCOUNT_SID`) or read from a file (`file:///home/james//.twilio_config`).

**from_number**: Your Twilio phone number. Must be in the format `+15558675309`.

**to_number**: The phone number to send an SMS message to.

### Actions

**send**: send an SMS message. Needs the body of the message to send. You can use interpolation to send somthing from the event received:

```json
"Actions": [
  {
    "Bee":"Twilio Example Bee",
    "Name":"send",
    "Options":[
      {
        "Name":"body",
        "Value":"Example body with interpolation: {{.something_from_event}}"
      }
    ]
  }
]
```

## Credits

Twilio logo: <https://www.twilio.com/press>

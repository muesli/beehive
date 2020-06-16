# Telegram Bees

[Telegram](https://telegram.org) bot that sends and receives messages to/from Telegram chats and groups.


## Configuration

Bee options:

```json
"Bees": [{
    "Name": "telegram",
    "Class": "telegrambee",
    "Description": "Telegram bot bee",
    "Options": [
        {
            "Name": "api_key",
            "Value": "file:///Users/lalotone/.telegramapi"
        },
        {
            "Name": "formatting_enabled",
            "Value": false
        }
    ]
}]
```

**api_key**: [Telegram bot](https://core.telegram.org/bots) API Key. Can be added
to the recipe, via environment variable (`env://MY_API_KEY`) or read from a file (`file:///Users/lalotone/.telegram_key`)

**formatting_enabled**: Enable [HTML message formatting](https://core.telegram.org/bots/api#html-style).

## Credits

Telegram image by Telegram Messenger LLP - User:Javitomad, Public Domain, https://commons.wikimedia.org/w/index.php?curid=36861817

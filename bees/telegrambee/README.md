# Telegram Bee

[Telegram](https://telegram.org) bot that sends and receives messages to/from
Telegram chats and groups.

## Configuration

Bee options:

```json
"Bees": [{
    "Name": "telegram",
    "Class": "telegrambee",
    "Description": "Telegram bot bee",
    "Options": [{
        "Name": "apiKey",
        "Value": "file:///Users/lalotone/.telegramapi"
    }]
}]
```

**apiKey**: [Telegram bot](https://core.telegram.org/bots) API Key. Can be
added to the recipe, via environment variable (`env://MY_API_KEY`) or read
from a file (`file:///Users/lalotone/.telegram_key`)

## Credits

Telegram image by Telegram Messenger LLP - User:Javitomad, Public Domain,
<https://commons.wikimedia.org/w/index.php?curid=36861817>

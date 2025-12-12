# Receiving Messages

There are two ways to receive messages from users:

1. **Polling** - Periodically fetch new messages
2. **Webhooks** - Receive messages in real-time (recommended)

## Polling with getUpdates

**Endpoint:** `GET /bot/getUpdates`

**Headers:** `Authorization: Bearer YOUR_API_TOKEN`

```bash
curl https://api.solafon.com/api/v1/bot/getUpdates \
  -H "Authorization: Bearer YOUR_API_TOKEN"
```

**Response:**
```json
{
  "ok": true,
  "result": [
    {
      "update_id": 1,
      "message": {
        "message_id": 1,
        "from": {
          "id": 123,
          "email": "user@example.com",
          "name": "John",
          "language": "en"
        },
        "chat": {"id": 123, "type": "private"},
        "date": 1704067200,
        "text": "/start"
      }
    }
  ]
}
```

Messages are marked as read after retrieval.

## Polling Example

### Python

```python
import requests
import time

API_TOKEN = "your_api_token"
BASE_URL = "https://api.solafon.com/api/v1/bot"

def get_updates():
    response = requests.get(
        f"{BASE_URL}/getUpdates",
        headers={"Authorization": f"Bearer {API_TOKEN}"}
    )
    return response.json()

def send_message(chat_id, text):
    requests.post(
        f"{BASE_URL}/sendMessage",
        headers={"Authorization": f"Bearer {API_TOKEN}"},
        json={"chat_id": chat_id, "text": text}
    )

def handle_message(message):
    chat_id = message["chat"]["id"]
    text = message.get("text", "")

    if text == "/start":
        send_message(chat_id, "Welcome to my app!")
    elif text == "/help":
        send_message(chat_id, "Available commands: /start, /help")
    else:
        send_message(chat_id, f"You said: {text}")

# Main loop
while True:
    updates = get_updates()
    if updates.get("ok"):
        for update in updates.get("result", []):
            if "message" in update:
                handle_message(update["message"])
    time.sleep(1)  # Poll every second
```

### Node.js

```javascript
const axios = require('axios');

const API_TOKEN = 'your_api_token';
const BASE_URL = 'https://api.solafon.com/api/v1/bot';

async function getUpdates() {
  const response = await axios.get(`${BASE_URL}/getUpdates`, {
    headers: { Authorization: `Bearer ${API_TOKEN}` }
  });
  return response.data;
}

async function sendMessage(chatId, text) {
  await axios.post(`${BASE_URL}/sendMessage`,
    { chat_id: chatId, text },
    { headers: { Authorization: `Bearer ${API_TOKEN}` } }
  );
}

async function handleMessage(message) {
  const chatId = message.chat.id;
  const text = message.text || '';

  if (text === '/start') {
    await sendMessage(chatId, 'Welcome to my app!');
  } else if (text === '/help') {
    await sendMessage(chatId, 'Available commands: /start, /help');
  } else {
    await sendMessage(chatId, `You said: ${text}`);
  }
}

// Main loop
async function poll() {
  while (true) {
    try {
      const updates = await getUpdates();
      if (updates.ok) {
        for (const update of updates.result || []) {
          if (update.message) {
            await handleMessage(update.message);
          }
        }
      }
    } catch (error) {
      console.error('Error:', error.message);
    }
    await new Promise(r => setTimeout(r, 1000));
  }
}

poll();
```

## Message Object

| Field | Type | Description |
|-------|------|-------------|
| message_id | integer | Unique message identifier |
| from | object | User who sent the message |
| from.id | integer | User ID (use as chat_id for replies) |
| from.email | string | User's email |
| from.name | string | User's display name |
| from.language | string | User's preferred language |
| chat | object | Chat information |
| chat.id | integer | Chat ID (same as user ID for private chats) |
| chat.type | string | Always "private" |
| date | integer | Unix timestamp |
| text | string | Message text |

## Polling vs Webhooks

| Feature | Polling | Webhooks |
|---------|---------|----------|
| Setup | Simple | Requires public URL |
| Latency | 1+ seconds | Instant |
| Server Load | Higher | Lower |
| Reliability | Guaranteed | Depends on your server |

For production apps, [webhooks](webhooks.md) are recommended.

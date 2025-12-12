# Developer API Overview

The Developer API allows you to build interactive apps on Solafon platform.

## Base URL

```
https://api.solafon.com/api/v1/bot
```

## Authentication

All Developer API requests require your app's API token:

```bash
Authorization: Bearer YOUR_API_TOKEN
```

Or simply:
```bash
Authorization: YOUR_API_TOKEN
```

## Getting Your API Token

1. Create an app via [Dev Studio](../dev-studio/overview.md) or the Apps API
2. Your token is provided upon creation
3. Store it securely - you can regenerate if compromised

## Available Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | /bot/getMe | Get app info |
| POST | /bot/sendMessage | Send message to user |
| GET | /bot/getUpdates | Get pending messages (polling) |
| POST | /bot/setWebhook | Set webhook URL |
| POST | /bot/deleteWebhook | Remove webhook |
| GET | /bot/getWebhookInfo | Get webhook status |
| POST | /bot/setMyCommands | Set app commands |
| GET | /bot/getMyCommands | Get app commands |

## Message Flow

### Option 1: Polling

```
User → Solafon → Your Server (polls getUpdates)
Your Server → sendMessage → Solafon → User
```

### Option 2: Webhooks (Recommended)

```
User → Solafon → Webhook POST to your server
Your Server → sendMessage → Solafon → User
```

## Quick Example

```python
import requests

API_TOKEN = "your_api_token"
BASE_URL = "https://api.solafon.com/api/v1/bot"

# Get app info
response = requests.get(
    f"{BASE_URL}/getMe",
    headers={"Authorization": f"Bearer {API_TOKEN}"}
)
print(response.json())

# Send a message
response = requests.post(
    f"{BASE_URL}/sendMessage",
    headers={"Authorization": f"Bearer {API_TOKEN}"},
    json={
        "chat_id": 123,
        "text": "Hello from my app!"
    }
)
print(response.json())
```

## Response Format

All responses follow this structure:

**Success:**
```json
{
  "ok": true,
  "result": {...}
}
```

**Error:**
```json
{
  "ok": false,
  "error_code": 400,
  "description": "Bad Request: chat not found"
}
```

## Next Steps

- [Send Messages](send-message.md)
- [Receive Messages](receive-messages.md)
- [Set Up Webhooks](webhooks.md)
- [Configure Commands](commands.md)

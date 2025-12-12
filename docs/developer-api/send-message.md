# Sending Messages

Send messages to users from your app.

## Send Message

**Endpoint:** `POST /bot/sendMessage`

**Headers:** `Authorization: Bearer YOUR_API_TOKEN`

### Basic Text Message

```bash
curl -X POST https://api.solafon.com/api/v1/bot/sendMessage \
  -H "Authorization: Bearer YOUR_API_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "chat_id": 123,
    "text": "Hello, user!"
  }'
```

**Response:**
```json
{
  "ok": true,
  "result": {
    "message_id": 456,
    "chat": {"id": 123, "type": "private"},
    "date": 1704067200,
    "text": "Hello, user!"
  }
}
```

### Request Parameters

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| chat_id | integer | Yes | User ID to send message to |
| text | string | Yes | Message text |
| message_type | string | No | Message type (default: "text") |
| metadata | string | No | JSON metadata for rich content |

### Message with Buttons

```bash
curl -X POST https://api.solafon.com/api/v1/bot/sendMessage \
  -H "Authorization: Bearer YOUR_API_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "chat_id": 123,
    "text": "Choose an option:",
    "message_type": "button",
    "metadata": "{\"buttons\": [{\"text\": \"Option 1\"}, {\"text\": \"Option 2\"}]}"
  }'
```

## Code Examples

### Python

```python
import requests

def send_message(chat_id, text):
    response = requests.post(
        "https://api.solafon.com/api/v1/bot/sendMessage",
        headers={"Authorization": f"Bearer {API_TOKEN}"},
        json={"chat_id": chat_id, "text": text}
    )
    return response.json()

# Usage
result = send_message(123, "Hello!")
print(result)
```

### Node.js

```javascript
const axios = require('axios');

async function sendMessage(chatId, text) {
  const response = await axios.post(
    'https://api.solafon.com/api/v1/bot/sendMessage',
    { chat_id: chatId, text: text },
    { headers: { Authorization: `Bearer ${API_TOKEN}` } }
  );
  return response.data;
}

// Usage
sendMessage(123, 'Hello!').then(console.log);
```

### Go

```go
package main

import (
    "bytes"
    "encoding/json"
    "net/http"
)

func sendMessage(chatID int, text string) error {
    body, _ := json.Marshal(map[string]interface{}{
        "chat_id": chatID,
        "text":    text,
    })

    req, _ := http.NewRequest("POST",
        "https://api.solafon.com/api/v1/bot/sendMessage",
        bytes.NewBuffer(body))
    req.Header.Set("Authorization", "Bearer "+apiToken)
    req.Header.Set("Content-Type", "application/json")

    client := &http.Client{}
    _, err := client.Do(req)
    return err
}
```

## Error Handling

| Error Code | Description |
|------------|-------------|
| 400 | Bad request - missing required fields |
| 401 | Unauthorized - invalid API token |
| 404 | Chat not found - user hasn't interacted with app |
| 429 | Too many requests - rate limited |

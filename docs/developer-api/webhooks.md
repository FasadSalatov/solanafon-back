# Webhooks

Webhooks deliver messages to your server in real-time.

## Set Webhook

**Endpoint:** `POST /bot/setWebhook`

```bash
curl -X POST https://api.solafon.com/api/v1/bot/setWebhook \
  -H "Authorization: Bearer YOUR_API_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"url": "https://your-server.com/webhook"}'
```

**Response:**
```json
{
  "ok": true,
  "result": true,
  "description": "Webhook was set"
}
```

> **Important:** URL must use HTTPS.

---

## Delete Webhook

**Endpoint:** `POST /bot/deleteWebhook`

```bash
curl -X POST https://api.solafon.com/api/v1/bot/deleteWebhook \
  -H "Authorization: Bearer YOUR_API_TOKEN"
```

---

## Get Webhook Info

**Endpoint:** `GET /bot/getWebhookInfo`

```bash
curl https://api.solafon.com/api/v1/bot/getWebhookInfo \
  -H "Authorization: Bearer YOUR_API_TOKEN"
```

**Response:**
```json
{
  "ok": true,
  "result": {
    "url": "https://your-server.com/webhook",
    "has_custom_certificate": false,
    "pending_update_count": 5
  }
}
```

---

## Webhook Payload

When a user sends a message, this payload is POST'd to your webhook URL:

```json
{
  "update_id": 123,
  "message_id": 123,
  "from": {
    "id": 456,
    "email": "user@example.com",
    "name": "John",
    "language": "en"
  },
  "chat": {
    "id": 456,
    "type": "private"
  },
  "date": 1704067200,
  "text": "/start"
}
```

Your server should respond with HTTP 200 OK.

---

## Webhook Server Examples

### Python (Flask)

```python
from flask import Flask, request
import requests

app = Flask(__name__)
API_TOKEN = "your_api_token"

@app.route('/webhook', methods=['POST'])
def webhook():
    data = request.json
    chat_id = data['from']['id']
    text = data.get('text', '')

    # Process message
    if text == '/start':
        reply = "Welcome to my app!"
    else:
        reply = f"You said: {text}"

    # Send reply
    requests.post(
        'https://api.solafon.com/api/v1/bot/sendMessage',
        headers={'Authorization': f'Bearer {API_TOKEN}'},
        json={'chat_id': chat_id, 'text': reply}
    )

    return 'OK', 200

if __name__ == '__main__':
    app.run(port=8000)
```

### Node.js (Express)

```javascript
const express = require('express');
const axios = require('axios');

const app = express();
app.use(express.json());

const API_TOKEN = 'your_api_token';

app.post('/webhook', async (req, res) => {
  const { from, text } = req.body;
  const chatId = from.id;

  // Process message
  let reply;
  if (text === '/start') {
    reply = 'Welcome to my app!';
  } else {
    reply = `You said: ${text}`;
  }

  // Send reply
  await axios.post(
    'https://api.solafon.com/api/v1/bot/sendMessage',
    { chat_id: chatId, text: reply },
    { headers: { Authorization: `Bearer ${API_TOKEN}` } }
  );

  res.status(200).send('OK');
});

app.listen(8000);
```

### Go

```go
package main

import (
    "bytes"
    "encoding/json"
    "fmt"
    "net/http"
)

var apiToken = "your_api_token"

type WebhookPayload struct {
    From struct {
        ID int `json:"id"`
    } `json:"from"`
    Text string `json:"text"`
}

func sendMessage(chatID int, text string) {
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
    client.Do(req)
}

func webhookHandler(w http.ResponseWriter, r *http.Request) {
    var payload WebhookPayload
    json.NewDecoder(r.Body).Decode(&payload)

    var reply string
    if payload.Text == "/start" {
        reply = "Welcome to my app!"
    } else {
        reply = fmt.Sprintf("You said: %s", payload.Text)
    }

    sendMessage(payload.From.ID, reply)
    w.WriteHeader(http.StatusOK)
}

func main() {
    http.HandleFunc("/webhook", webhookHandler)
    http.ListenAndServe(":8000", nil)
}
```

---

## Troubleshooting

| Issue | Solution |
|-------|----------|
| Webhook not receiving | Check URL is HTTPS and publicly accessible |
| 5xx errors | Check server logs, ensure 200 response |
| Timeouts | Respond within 10 seconds |
| Duplicate messages | Implement idempotency with update_id |

## Webhook Logs

View webhook delivery logs via App Settings:

```bash
curl https://api.solafon.com/api/v1/apps/YOUR_APP_ID/settings \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

The `webhookLogs` array shows recent delivery attempts with status codes and response times.

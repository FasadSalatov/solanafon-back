# App Messages API

Endpoints for chat interaction with apps.

## Get App Messages

Get chat history with an app.

**Endpoint:** `GET /apps/:id/messages`

**Headers:** `Authorization: Bearer <token>` (required)

**Response:**
```json
{
  "messages": [
    {
      "id": 1,
      "appId": 1,
      "userId": 1,
      "content": "/start",
      "isFromBot": false,
      "isRead": true,
      "createdAt": "2024-01-01T00:00:00Z"
    },
    {
      "id": 2,
      "appId": 1,
      "userId": 1,
      "content": "Welcome to the app!",
      "isFromBot": true,
      "isRead": true,
      "createdAt": "2024-01-01T00:00:01Z"
    }
  ],
  "total": 2
}
```

---

## Send Message to App

Send a message to an app.

**Endpoint:** `POST /apps/:id/messages`

**Headers:** `Authorization: Bearer <token>` (required)

**Request:**
```json
{
  "content": "/start"
}
```

**Response:**
```json
{
  "userMessage": {
    "id": 1,
    "appId": 1,
    "userId": 1,
    "content": "/start",
    "isFromBot": false,
    "isRead": true,
    "createdAt": "2024-01-01T00:00:00Z"
  },
  "botMessage": {
    "id": 2,
    "appId": 1,
    "userId": 1,
    "content": "Welcome to the app!",
    "isFromBot": true,
    "isRead": false,
    "createdAt": "2024-01-01T00:00:01Z"
  }
}
```

## Message Flow

1. User sends a message via `POST /apps/:id/messages`
2. If the message matches a defined command, auto-response is returned
3. If webhook is configured, message is forwarded to the webhook URL
4. App developer can respond via Developer API

## Message Types

| Type | Description |
|------|-------------|
| text | Plain text message |
| image | Image with URL |
| button | Message with interactive buttons |

Messages can include metadata for rich content:

```json
{
  "content": "Choose an option:",
  "messageType": "button",
  "metadata": {
    "buttons": [
      {"text": "Option 1", "callback": "opt1"},
      {"text": "Option 2", "callback": "opt2"}
    ]
  }
}
```

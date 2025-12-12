# App Commands

Define commands that your app responds to automatically.

## Set Commands

**Endpoint:** `POST /bot/setMyCommands`

```bash
curl -X POST https://api.solafon.com/api/v1/bot/setMyCommands \
  -H "Authorization: Bearer YOUR_API_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "commands": [
      {"command": "start", "description": "Start the app"},
      {"command": "help", "description": "Get help"},
      {"command": "settings", "description": "App settings"}
    ]
  }'
```

**Response:**
```json
{
  "ok": true,
  "result": true
}
```

> **Note:** Don't include the `/` prefix in the command name.

---

## Get Commands

**Endpoint:** `GET /bot/getMyCommands`

```bash
curl https://api.solafon.com/api/v1/bot/getMyCommands \
  -H "Authorization: Bearer YOUR_API_TOKEN"
```

**Response:**
```json
{
  "ok": true,
  "result": [
    {"command": "start", "description": "Start the app"},
    {"command": "help", "description": "Get help"}
  ]
}
```

---

## Commands via REST API

You can also manage commands via the Apps API (with JWT token):

### Add Command

**Endpoint:** `POST /apps/:id/commands`

```bash
curl -X POST https://api.solafon.com/api/v1/apps/123/commands \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "command": "/help",
    "description": "Show help information",
    "response": "Available commands:\n/start - Begin\n/help - This message"
  }'
```

### Update Command

**Endpoint:** `PUT /apps/:id/commands/:cmdId`

```bash
curl -X PUT https://api.solafon.com/api/v1/apps/123/commands/1 \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "description": "Updated description",
    "response": "Updated response text"
  }'
```

### Delete Command

**Endpoint:** `DELETE /apps/:id/commands/:cmdId`

```bash
curl -X DELETE https://api.solafon.com/api/v1/apps/123/commands/1 \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

### Get All Commands

**Endpoint:** `GET /apps/:id/commands`

```bash
curl https://api.solafon.com/api/v1/apps/123/commands \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

**Response:**
```json
{
  "commands": [
    {
      "id": 1,
      "command": "/start",
      "description": "Start the app",
      "response": "Welcome!",
      "isEnabled": true
    },
    {
      "id": 2,
      "command": "/help",
      "description": "Get help",
      "response": "Help text here...",
      "isEnabled": true
    }
  ]
}
```

---

## Auto-Response

When you define a command with a `response`, the platform automatically replies when users send that command.

**Example:**
```json
{
  "command": "/start",
  "description": "Start the app",
  "response": "Welcome to My App!\n\nType /help for available commands."
}
```

When a user sends `/start`, they immediately receive:
```
Welcome to My App!

Type /help for available commands.
```

No server-side code required for simple commands!

---

## Best Practices

1. **Always define /start** - First command users will send
2. **Include /help** - Show all available commands
3. **Keep descriptions short** - 50 characters max
4. **Use clear command names** - Lowercase, no spaces
5. **Handle unknown commands** - Provide helpful error message

## Common Commands

| Command | Description |
|---------|-------------|
| /start | Initial greeting and app introduction |
| /help | List all available commands |
| /settings | User preferences or settings |
| /about | Information about the app |
| /cancel | Cancel current operation |

# Creating Your First App

This guide walks you through creating your first app on Solafon.

## Using Dev Studio (Recommended)

Dev Studio is the official app for creating and managing your apps. It provides an interactive chat-based interface.

### Step 1: Start Dev Studio

```bash
curl -X POST https://api.solafon.com/api/v1/apps/devstudio/message \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"content": "/start"}'
```

### Step 2: Create New App

```bash
curl -X POST https://api.solafon.com/api/v1/apps/devstudio/message \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"content": "/newapp"}'
```

### Step 3: Follow the Prompts

Dev Studio will guide you through:

1. **App Name** - Give your app a descriptive name (3-50 characters)
2. **Description** - Explain what your app does (min 10 characters)
3. **Icon** - Choose an emoji to represent your app
4. **Category** - Select from: AI, Games, Trading, DePIN, DeFi, NFT, Staking, Services
5. **Username** - Create a unique identifier ending with 'app' (e.g., `mytradingapp`)
6. **Welcome Message** - The first message users see (or skip with `/skip`)

### Step 4: Receive Your API Token

After creation, you'll receive:
- Confirmation of your app details
- Your API token (save this securely!)
- Moderation status (apps are reviewed within 24 hours)

---

## Using the API Directly

You can also create apps via the REST API:

**Endpoint:** `POST /apps`

```bash
curl -X POST https://api.solafon.com/api/v1/apps \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "icon": "🎮",
    "title": "My Game App",
    "description": "An awesome gaming experience",
    "categoryId": 2,
    "botUsername": "mygameapp",
    "welcomeMessage": "Welcome to My Game! Type /play to start.",
    "webhookUrl": "https://myserver.com/webhook"
  }'
```

**Response:**
```json
{
  "message": "App created successfully. It will be reviewed by moderators within 24 hours.",
  "app": {
    "id": 123,
    "title": "My Game App",
    "icon": "🎮",
    "botUsername": "mygameapp",
    "moderationStatus": "pending"
  },
  "apiToken": "abc123def456..."
}
```

## After Creation

### Set Up Commands

Define automatic responses for common commands:

```bash
# Via Dev Studio
curl -X POST https://api.solafon.com/api/v1/apps/devstudio/message \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{"content": "/commands"}'

# Via API
curl -X POST https://api.solafon.com/api/v1/apps/123/commands \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "command": "/help",
    "description": "Show help information",
    "response": "Available commands:\n/play - Start playing\n/stats - View your stats"
  }'
```

### Configure Webhook

Receive messages in real-time:

```bash
curl -X POST https://api.solafon.com/api/v1/bot/setWebhook \
  -H "Authorization: Bearer YOUR_API_TOKEN" \
  -d '{"url": "https://your-server.com/webhook"}'
```

### Wait for Moderation

- Apps are reviewed within 24 hours
- You'll be notified when approved/rejected
- Rejected apps include feedback on what to fix

## Next Steps

- [Set up webhooks](../developer-api/webhooks.md)
- [Learn about the Developer API](../developer-api/overview.md)
- [Manage your apps](../dev-studio/managing-apps.md)

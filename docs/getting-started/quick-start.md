# Quick Start

This guide will help you get started with Solafon API in just a few minutes.

## 1. Create an Account

First, authenticate via email OTP:

```bash
# Request OTP
curl -X POST https://api.solafon.com/api/v1/auth/email/request \
  -H "Content-Type: application/json" \
  -d '{"email": "your@email.com"}'

# Verify OTP (check your email)
curl -X POST https://api.solafon.com/api/v1/auth/email/verify \
  -H "Content-Type: application/json" \
  -d '{"email": "your@email.com", "code": "123456"}'
```

You'll receive a JWT token and 500 Mana Points as a welcome bonus.

## 2. Create Your First App

The easiest way is through Dev Studio. Send a message to start:

```bash
curl -X POST https://api.solafon.com/api/v1/apps/devstudio/message \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"content": "/newapp"}'
```

Follow the interactive prompts to:
1. Set your app name
2. Add a description
3. Choose an icon (emoji)
4. Select a category
5. Create a unique username (ending with 'app')
6. Set a welcome message

## 3. Get Your API Token

After creating an app, you'll receive an API token. Save it securely!

```bash
# Or request it again via Dev Studio
curl -X POST https://api.solafon.com/api/v1/apps/devstudio/message \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"content": "/token"}'
```

## 4. Receive Messages

### Option A: Polling

```bash
curl https://api.solafon.com/api/v1/bot/getUpdates \
  -H "Authorization: Bearer YOUR_API_TOKEN"
```

### Option B: Webhooks

Set up a webhook to receive messages in real-time:

```bash
curl -X POST https://api.solafon.com/api/v1/bot/setWebhook \
  -H "Authorization: Bearer YOUR_API_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"url": "https://your-server.com/webhook"}'
```

## 5. Send Messages

Reply to users:

```bash
curl -X POST https://api.solafon.com/api/v1/bot/sendMessage \
  -H "Authorization: Bearer YOUR_API_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"chat_id": 123, "text": "Hello from my app!"}'
```

## Next Steps

- [Learn about Authentication](authentication.md)
- [Explore the API Reference](../api/overview.md)
- [Set up Webhooks](../developer-api/webhooks.md)
- [Manage your apps via Dev Studio](../dev-studio/overview.md)

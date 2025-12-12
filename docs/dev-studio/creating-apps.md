# Creating Apps with Dev Studio

This guide walks you through creating an app using Dev Studio's interactive interface.

## Start the Wizard

Send `/newapp` to begin:

```bash
curl -X POST https://api.solafon.com/api/v1/apps/devstudio/message \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{"content": "/newapp"}'
```

## Step 1: App Name

Enter a name for your app (3-50 characters):

```
You: My Trading App
```

## Step 2: Description

Describe what your app does (minimum 10 characters):

```
You: A trading signals and market analysis platform
```

## Step 3: Icon

Choose an emoji to represent your app:

```
You: 📊
```

## Step 4: Category

Select from available categories (enter the number):

```
Dev Studio:
1. 🤖 AI
2. 🎮 Games
3. 📊 Trading
4. 🌐 DePIN
5. 💎 DeFi
6. 🖼️ NFT
7. 🔒 Staking
8. 🛠️ Services

You: 3
```

## Step 5: Username

Create a unique username ending with 'app':

```
You: mytradingapp
```

Requirements:
- Minimum 5 characters
- Only a-z, 0-9, and underscore
- Must end with 'app'
- Must be unique

## Step 6: Welcome Message

Set the first message users see (or `/skip` to skip):

```
You: Welcome to My Trading App! 📊

Get real-time trading signals and market analysis.

Commands:
/signals - Latest signals
/portfolio - Your portfolio
/settings - Preferences
```

## Completion

After completing all steps, you'll receive:

```
🎉 Congratulations! App created!

📊 My Trading App
@mytradingapp

🔑 API Token (save it!):
abc123def456...

📋 Status: ⏳ Under moderation
App will be reviewed within 24 hours.

Next steps:
• /commands - add commands
• /webhook - set up webhook
• /token - show token again

API Documentation: /help
```

## Important Notes

1. **Save your API token** - You'll need it for the Developer API
2. **Moderation takes up to 24 hours** - Your app will be reviewed
3. **Username is permanent** - Choose carefully, it can't be changed
4. **Welcome message is optional** - But recommended for user experience

## What's Next?

After creating your app:

1. **Set up commands** - Define auto-responses with `/commands`
2. **Configure webhook** - Receive messages in real-time with `/webhook`
3. **Build your backend** - Use the [Developer API](../developer-api/overview.md)
4. **Wait for approval** - Moderation typically completes within 24 hours

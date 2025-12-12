# Dev Studio Overview

Dev Studio is the official app for creating and managing your apps on Solafon.

## What is Dev Studio?

Dev Studio provides an interactive chat-based interface for developers. Instead of calling APIs directly, you can create and manage apps through a conversational interface.

## Accessing Dev Studio

```bash
curl -X POST https://api.solafon.com/api/v1/apps/devstudio/message \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"content": "/start"}'
```

## Features

### Create Apps
Use `/newapp` to create a new app with a guided wizard:
1. App name
2. Description
3. Icon (emoji)
4. Category
5. Unique username
6. Welcome message

### Manage Apps
- `/myapps` - View all your apps
- `/edit` - Edit app settings
- `/delete` - Remove an app

### Configure
- `/commands` - Set up auto-responses
- `/webhook` - Configure webhook URL
- `/token` - Get or view API token

### Help
- `/help` - List all commands
- `/cancel` - Cancel current action

## Example Session

```
You: /start

Dev Studio: Welcome to Dev Studio!
Here you can:
• Create new apps
• Configure commands and auto-responses
• Get API tokens for integrations
• Set up webhooks

Start with /newapp to create your first app!

You: /newapp

Dev Studio: Great! Let's create a new app.
What will your app be called?
(Example: "Crypto Trading" or "NFT Gallery")

You: My Trading App

Dev Studio: Excellent name: My Trading App!
Now describe what your app does (brief description):

You: A trading signals and market analysis tool

Dev Studio: Now choose an icon (emoji) for your app:
(Send one emoji, e.g.: 📊 💰 🔥)

...
```

## Next Steps

- [Creating Apps](creating-apps.md)
- [Managing Apps](managing-apps.md)
- [Commands Reference](commands.md)

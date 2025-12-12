# Dev Studio Commands Reference

Complete list of all Dev Studio commands.

## Basic Commands

| Command | Description |
|---------|-------------|
| `/start` | Welcome message and introduction |
| `/help` | Show all available commands |
| `/cancel` | Cancel current action |

## App Management

| Command | Description |
|---------|-------------|
| `/newapp` | Create a new app |
| `/myapps` | List all your apps |
| `/edit` | Edit an existing app |
| `/delete` | Delete an app |

## Configuration

| Command | Description |
|---------|-------------|
| `/token` | Get or view API token |
| `/commands` | Manage app commands |
| `/webhook` | Configure webhook URL |

---

## Command Details

### /start

Shows welcome message and capabilities overview.

```
Dev Studio: Welcome to Dev Studio!

Here you can:
• Create new apps
• Configure commands and auto-responses
• Get API tokens for integrations
• Set up webhooks

Start with /newapp to create your first app!

/help - all commands
```

---

### /help

Lists all available commands with descriptions.

```
Dev Studio:
📋 Available commands:

🆕 Create and manage
/newapp - Create new app
/myapps - List my apps
/edit - Edit app
/delete - Delete app

⚙️ App settings
/token - Get/view API token
/commands - Configure commands
/webhook - Set up webhook

❌ /cancel - Cancel current action
```

---

### /cancel

Cancels any ongoing operation and returns to idle state.

```
Dev Studio: Action cancelled. How can I help?

Use /help for list of commands.
```

---

### /newapp

Starts the app creation wizard. See [Creating Apps](creating-apps.md) for full walkthrough.

Steps:
1. App name (3-50 chars)
2. Description (min 10 chars)
3. Icon (emoji)
4. Category (1-8)
5. Username (ends with 'app')
6. Welcome message (or /skip)

---

### /myapps

Shows all apps created by you.

```
Dev Studio:
📱 Your apps:

1. 📊 My Trading App
   Username: @mytradingapp
   Status: ✅ Active
   Users: 1,500

2. 🎮 Game Zone
   Username: @gamezoneapp
   Status: ⏳ Under moderation
   Users: 0

Use /edit to modify or /token to get token
```

---

### /edit

Opens the edit menu for selected app.

Options:
1. 📝 Name
2. 📄 Description
3. 📋 Commands
4. 🔗 Webhook
5. 🔑 Show token

---

### /delete

Permanently deletes an app after confirmation.

Requires typing `YES` to confirm.

---

### /token

Shows API token for selected app.

If you have multiple apps, you'll be asked to select one.

---

### /commands

Opens command management for selected app.

Actions:
- Enter `/command` to add new command
- Enter `delete /command` to remove existing
- Enter `/cancel` to exit

---

### /webhook

Configure webhook URL for selected app.

- Enter HTTPS URL to set webhook
- Enter `clear` to remove webhook
- Enter `/cancel` to exit

---

## Status Icons

| Icon | Meaning |
|------|---------|
| ✅ | Active / Approved |
| ⏳ | Under moderation |
| ❌ | Rejected |

## Error Messages

| Message | Meaning |
|---------|---------|
| "You don't have any apps" | Create one with /newapp |
| "Name too short" | Minimum 3 characters |
| "Username taken" | Choose a different username |
| "Must end with 'app'" | Username format requirement |
| "URL must start with https://" | Webhook URL requirement |

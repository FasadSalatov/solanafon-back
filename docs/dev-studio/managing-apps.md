# Managing Apps with Dev Studio

After creating apps, use Dev Studio to manage them.

## View Your Apps

List all apps you've created:

```
You: /myapps

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

Use /edit to modify or /token to get API token
```

## Edit an App

```
You: /edit

Dev Studio:
Select app to edit (enter number):

1. 📊 My Trading App (@mytradingapp)
2. 🎮 Game Zone (@gamezoneapp)

You: 1

Dev Studio:
⚙️ Editing: 📊 My Trading App

Choose what to change:

1. 📝 Name
2. 📄 Description
3. 📋 Commands
4. 🔗 Webhook
5. 🔑 Show token

Enter number or /cancel to exit:
```

### Edit Name

```
You: 1

Dev Studio:
Current name: My Trading App

Enter new name:

You: Pro Trading App

Dev Studio:
✅ Name changed to: Pro Trading App

App sent for re-moderation.
```

### Edit Description

```
You: 2

Dev Studio:
Current description: A trading signals platform

Enter new description:

You: Professional trading signals with AI-powered analysis

Dev Studio:
✅ Description updated!

App sent for re-moderation.
```

## Manage Commands

```
You: /commands

Dev Studio:
📋 Commands for My Trading App:

/start - Start the app
/signals - Get latest signals
/help - Show help

Enter new command (e.g., /menu)
Or 'delete /command' to remove

/cancel - exit
```

### Add a Command

```
You: /portfolio

Dev Studio:
Command: /portfolio

Enter brief description:

You: View your portfolio

Dev Studio:
Now enter the app's response to this command:

You: 📊 Your Portfolio

Loading your positions...

Use /buy or /sell to trade.

Dev Studio:
✅ Command /portfolio added!
```

### Delete a Command

```
You: delete /signals

Dev Studio:
✅ Command /signals deleted!
```

## Configure Webhook

```
You: /webhook

Dev Studio:
🔗 Webhook for My Trading App

Current URL: not set

Enter new webhook URL:

You: https://myserver.com/webhook

Dev Studio:
✅ Webhook set:
https://myserver.com/webhook

Now user messages will be sent to this URL.
```

### Clear Webhook

```
You: /webhook
You: clear

Dev Studio:
✅ Webhook removed
```

## Get API Token

```
You: /token

Dev Studio:
🔑 API Token for 📊 My Trading App

abc123def456...

⚠️ Keep this token secure!

To regenerate, use /edit → Token
```

## Delete an App

```
You: /delete

Dev Studio:
Select app to delete (enter number):

1. 📊 My Trading App (@mytradingapp)

You: 1

Dev Studio:
⚠️ Are you sure you want to delete "My Trading App"?

Enter YES to confirm or /cancel to abort.

You: YES

Dev Studio:
✅ App "My Trading App" deleted.
```

## Tips

1. **Changes require re-moderation** - Name and description changes reset approval
2. **Webhook URL must be HTTPS** - HTTP URLs won't be accepted
3. **Test commands locally first** - Make sure responses are formatted correctly
4. **Keep tokens secure** - Regenerate if compromised via /edit → Token

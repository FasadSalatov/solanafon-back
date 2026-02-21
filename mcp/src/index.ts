#!/usr/bin/env node

/**
 * Solafon MCP Server
 *
 * Model Context Protocol server for building mini-apps and bots
 * on the Solafon platform. Provides AI assistants with tools to:
 * - Manage apps via Bot API
 * - Read platform documentation
 * - Scaffold code for integrations
 * - Debug webhooks and message flows
 */

import { McpServer } from "@modelcontextprotocol/sdk/server/mcp.js";
import { StdioServerTransport } from "@modelcontextprotocol/sdk/server/stdio.js";
import { z } from "zod";

// ---------------------------------------------------------------------------
// Configuration
// ---------------------------------------------------------------------------

const API_BASE_URL = process.env.SOLAFON_API_URL || "https://api.solafon.com/api/v1";
const BOT_TOKEN = process.env.SOLAFON_BOT_TOKEN || "";

// ---------------------------------------------------------------------------
// HTTP Client
// ---------------------------------------------------------------------------

interface ApiResponse {
  ok?: boolean;
  error?: string;
  error_code?: number;
  description?: string;
  [key: string]: unknown;
}

async function apiRequest(
  method: string,
  path: string,
  body?: Record<string, unknown>,
  token?: string
): Promise<ApiResponse> {
  const url = `${API_BASE_URL}${path}`;
  const headers: Record<string, string> = {
    "Content-Type": "application/json",
  };

  const authToken = token || BOT_TOKEN;
  if (authToken) {
    headers["Authorization"] = `Bearer ${authToken}`;
  }

  const options: RequestInit = { method, headers };
  if (body && (method === "POST" || method === "PUT" || method === "PATCH")) {
    options.body = JSON.stringify(body);
  }

  const response = await fetch(url, options);
  const data = await response.json();
  return data as ApiResponse;
}

// ---------------------------------------------------------------------------
// Documentation Content
// ---------------------------------------------------------------------------

const DOCS: Record<string, { title: string; content: string }> = {
  introduction: {
    title: "Introduction to Solafon",
    content: `# Solafon Platform

Solafon is a platform for creating and deploying mini-applications (mini-apps) that users interact with through a chat-based interface.

## For Users
- Browse and discover mini-apps by category (AI, Games, Trading, DePIN, DeFi, NFT, Staking, Services)
- Interact with apps via chat interface
- Mana Points (MP) economy for premium features
- Secret Login with virtual phone numbers for privacy

## For Developers
- Create mini-apps via Dev Studio or REST API
- Telegram Bot API-compatible interface for message handling
- Webhooks and polling for receiving user messages
- Bot commands with auto-response capability
- Moderation system for app quality

## Base URL
\`https://api.solafon.com/api/v1\`

## Authentication
- **User Auth**: JWT via email OTP (for user-facing features)
- **Bot API Auth**: API Token (for bot/app integrations — this is what developers use)

## Quick Start
1. Create an app (via Dev Studio or POST /apps)
2. Get your API token
3. Set up webhooks or polling to receive messages
4. Send responses via POST /bot/sendMessage`,
  },

  "quick-start": {
    title: "Quick Start Guide",
    content: `# Quick Start — Build Your First Solafon Bot in 5 Minutes

## Step 1: Create Your App
POST to create a new mini-app:
\`\`\`bash
curl -X POST ${API_BASE_URL}/apps \\
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \\
  -H "Content-Type: application/json" \\
  -d '{
    "title": "My Bot",
    "description": "A helpful bot",
    "icon": "🤖",
    "category_id": 1,
    "bot_username": "mybotapp",
    "welcome_message": "Hello! How can I help?"
  }'
\`\`\`
Response includes your \`apiToken\` — save it!

## Step 2: Get Messages (Polling)
\`\`\`bash
curl -H "Authorization: Bearer YOUR_API_TOKEN" \\
  ${API_BASE_URL}/bot/getUpdates
\`\`\`

## Step 3: Send Response
\`\`\`bash
curl -X POST ${API_BASE_URL}/bot/sendMessage \\
  -H "Authorization: Bearer YOUR_API_TOKEN" \\
  -H "Content-Type: application/json" \\
  -d '{"chat_id": 1, "text": "Hello from my bot!"}'
\`\`\`

## Step 4: Set Up Webhook (Optional)
\`\`\`bash
curl -X POST ${API_BASE_URL}/bot/setWebhook \\
  -H "Authorization: Bearer YOUR_API_TOKEN" \\
  -H "Content-Type: application/json" \\
  -d '{"url": "https://your-server.com/webhook"}'
\`\`\`

## Step 5: Define Commands
\`\`\`bash
curl -X POST ${API_BASE_URL}/bot/setMyCommands \\
  -H "Authorization: Bearer YOUR_API_TOKEN" \\
  -H "Content-Type: application/json" \\
  -d '{"commands": [
    {"command": "start", "description": "Start the bot", "response": "Welcome!"},
    {"command": "help", "description": "Show help", "response": "Available commands: /start, /help"}
  ]}'
\`\`\`

Note: New users get 500 Mana Points welcome bonus. Apps go through moderation (pending → approved/rejected).`,
  },

  authentication: {
    title: "Authentication",
    content: `# Authentication

## Two Auth Systems

### 1. User Authentication (JWT)
For user-facing features (creating apps, profiles, etc.)

**Request OTP:**
\`\`\`
POST /auth/email/request
Body: {"email": "user@example.com"}
\`\`\`

**Verify OTP:**
\`\`\`
POST /auth/email/verify
Body: {"email": "user@example.com", "code": "123456"}
Response: {"token": "eyJ...", "user": {...}, "isNewUser": true}
\`\`\`

**Use JWT:**
\`\`\`
Authorization: Bearer eyJhbGciOiJIUzI1NiIs...
\`\`\`

### 2. Bot API Token
For developer integrations (sending messages, managing bots)

**Get your token:**
- Via Dev Studio: send \`/token\` command
- Via API: \`GET /apps/:id/settings\` (requires JWT)
- Generated when creating app: \`POST /apps\`

**Use API Token:**
\`\`\`
Authorization: Bearer your-api-token-here
# or just:
Authorization: your-api-token-here
\`\`\`

**Regenerate if compromised:**
\`\`\`
POST /apps/:id/regenerate-token
\`\`\`

**Verify your token works:**
\`\`\`
GET /bot/getMe
\`\`\`

## Security Best Practices
- Store tokens in environment variables, never in code
- Use HTTPS for all webhook URLs
- Regenerate tokens if compromised
- Each app has its own unique API token`,
  },

  "bot-api": {
    title: "Bot API Reference",
    content: `# Bot API Reference

Base URL: \`/api/v1/bot\`
Auth: API Token in Authorization header

## Endpoints

### GET /bot/getMe
Get bot info. Returns: id, title, username, is_verified, users_count, webhook_url.

### POST /bot/sendMessage
Send a message to a user.
\`\`\`json
{
  "chat_id": 123,          // Required: user ID
  "text": "Hello!",         // Required: message text
  "message_type": "text",   // Optional: "text" | "image" | "button"
  "metadata": "{}"          // Optional: JSON string for rich content
}
\`\`\`

**Button metadata example:**
\`\`\`json
{
  "message_type": "button",
  "metadata": "{\\"buttons\\": [{\\"text\\": \\"Buy\\", \\"action\\": \\"buy_item\\"}]}"
}
\`\`\`

### GET /bot/getUpdates
Long-polling for new messages. Returns unread messages, marks them as read.

**Message object:**
\`\`\`json
{
  "message_id": 1,
  "from": {"id": 1, "email": "user@mail.com", "name": "User", "language": "en"},
  "chat": {"id": 1, "type": "private"},
  "date": 1704067200,
  "text": "/start"
}
\`\`\`

### POST /bot/setWebhook
\`\`\`json
{"url": "https://your-server.com/webhook"}
\`\`\`
URL must be HTTPS. Solafon sends POST with Telegram-style update payload.

**Webhook payload:**
\`\`\`json
{
  "update_id": 1,
  "message": {
    "message_id": 1,
    "from": {"id": 1, "email": "...", "name": "..."},
    "chat": {"id": 1, "type": "private"},
    "date": 1704067200,
    "text": "Hello"
  }
}
\`\`\`

### POST /bot/deleteWebhook
Remove webhook configuration.

### GET /bot/getWebhookInfo
Returns: url, has_custom_certificate, pending_update_count.

### POST /bot/setMyCommands
\`\`\`json
{
  "commands": [
    {"command": "start", "description": "Start bot", "response": "Welcome!"},
    {"command": "help", "description": "Show help"}
  ]
}
\`\`\`
Note: command names WITHOUT the "/" prefix.
Commands with \`response\` field auto-reply without hitting your server.

### GET /bot/getMyCommands
Returns array of bot commands.

## Response Format
\`\`\`json
{"ok": true, "result": {...}}
{"ok": false, "error_code": 400, "description": "Bad request"}
\`\`\`

## Error Codes
- 400: Bad request (missing/invalid params)
- 401: Unauthorized (invalid/missing token)
- 403: Forbidden (app not approved yet)
- 404: Not found (user/app not found)
- 429: Too many requests (rate limited)`,
  },

  "app-management": {
    title: "App Management API",
    content: `# App Management API

All endpoints require JWT authentication.

## Create App
\`\`\`
POST /apps
{
  "title": "My App",           // Required, max 50 chars
  "description": "App desc",   // Optional
  "icon": "🤖",                // Optional, emoji
  "category_id": 1,            // Optional, 1-8
  "bot_username": "myapp",     // Optional, unique, a-z/0-9/_
  "welcome_message": "Hi!",    // Optional
  "is_secret": false           // Optional, secret apps only visible to Secret Login users
}
\`\`\`
Returns the created app WITH \`apiToken\`.

## List Apps
\`\`\`
GET /apps                    // All approved apps
GET /apps?category=ai        // Filter by category slug
GET /apps/search?q=trading   // Search by title/description
GET /apps/my                 // Your own apps
\`\`\`

## Get App Details
\`\`\`
GET /apps/:id
\`\`\`

## Update App
\`\`\`
PUT /apps/:id
{"title": "New Title", "description": "New desc"}
\`\`\`
Note: Changing title/description resets moderation status to "pending".

## Delete App
\`\`\`
DELETE /apps/:id
\`\`\`

## App Settings (Owner Only)
\`\`\`
GET /apps/:id/settings
\`\`\`
Returns: app details + apiToken + commands + recent webhookLogs.

## Regenerate API Token
\`\`\`
POST /apps/:id/regenerate-token
\`\`\`

## Bot Commands (REST)
\`\`\`
GET    /apps/:id/commands           // List commands
POST   /apps/:id/commands           // Add command
PUT    /apps/:id/commands/:cmdId    // Update command
DELETE /apps/:id/commands/:cmdId    // Delete command
\`\`\`

Command body:
\`\`\`json
{
  "command": "help",
  "description": "Show help",
  "response": "Here's how to use...",  // Auto-reply text
  "is_enabled": true
}
\`\`\`

## Categories
\`\`\`
GET /categories
GET /categories/:slug/apps
\`\`\`

Available categories:
| # | Name | Slug | Icon |
|---|------|------|------|
| 1 | AI | ai | 🤖 |
| 2 | Games | games | 🎮 |
| 3 | Trading | trading | 📊 |
| 4 | DePIN | depin | 🌐 |
| 5 | DeFi | defi | 💎 |
| 6 | NFT | nft | 🖼️ |
| 7 | Staking | staking | 🔒 |
| 8 | Services | services | 🛠️ |

## Messages
\`\`\`
GET  /apps/:id/messages              // Chat history
POST /apps/:id/messages              // Send message
     {"content": "Hello", "message_type": "text"}
\`\`\`

## Moderation
Apps go through moderation: pending → approved/rejected.
Only approved apps are visible to users and can use Bot API.`,
  },

  "dev-studio": {
    title: "Dev Studio Guide",
    content: `# Dev Studio

Dev Studio is a chat-based interface for creating and managing mini-apps.

## Access
\`\`\`
POST /apps/devstudio/message
{"content": "/start"}
\`\`\`

## Commands

### Basic
- \`/start\` — Welcome message and overview
- \`/help\` — List all available commands
- \`/cancel\` — Cancel current operation

### App Management
- \`/newapp\` — Start creating a new app (6-step wizard)
- \`/myapps\` — List your apps with status
- \`/edit\` — Edit an existing app
- \`/delete\` — Delete an app

### Configuration
- \`/token\` — View your API token
- \`/commands\` — Manage bot commands
- \`/webhook\` — Configure webhook URL

## Creating an App (/newapp)
Step-by-step wizard:
1. **Name** (3-50 chars)
2. **Description** (min 10 chars)
3. **Icon** (emoji)
4. **Category** (1-8)
5. **Username** (min 5 chars, a-z/0-9/_, must end with "app", unique)
6. **Welcome message** (or /skip)

## Editing an App (/edit)
Select app → choose what to edit:
1. Name
2. Description
3. Commands (add/delete)
4. Webhook URL (set HTTPS URL or "clear")
5. API Token (view)

## Status Icons
- ⏳ pending — awaiting moderation
- ✅ approved — live and accessible
- ❌ rejected — needs changes`,
  },

  "webhooks-guide": {
    title: "Webhooks Integration Guide",
    content: `# Webhooks Guide

## How Webhooks Work
1. User sends message to your app
2. If no matching auto-response command → Solafon POSTs to your webhook URL
3. Your server processes the message
4. You respond via POST /bot/sendMessage

## Setup
\`\`\`bash
curl -X POST ${API_BASE_URL}/bot/setWebhook \\
  -H "Authorization: Bearer YOUR_TOKEN" \\
  -d '{"url": "https://your-server.com/webhook"}'
\`\`\`

## Webhook Payload
\`\`\`json
{
  "update_id": 12345,
  "message": {
    "message_id": 67,
    "from": {
      "id": 1,
      "email": "user@mail.com",
      "name": "John",
      "language": "en"
    },
    "chat": {
      "id": 1,
      "type": "private"
    },
    "date": 1704067200,
    "text": "Hello bot!"
  }
}
\`\`\`

## Server Examples

### Node.js (Express)
\`\`\`javascript
const express = require('express');
const app = express();
app.use(express.json());

const BOT_TOKEN = process.env.SOLAFON_BOT_TOKEN;
const API_URL = 'https://api.solafon.com/api/v1';

app.post('/webhook', async (req, res) => {
  const { message } = req.body;
  const chatId = message.from.id;
  const text = message.text;

  // Process message and send response
  await fetch(\`\${API_URL}/bot/sendMessage\`, {
    method: 'POST',
    headers: {
      'Authorization': \`Bearer \${BOT_TOKEN}\`,
      'Content-Type': 'application/json'
    },
    body: JSON.stringify({
      chat_id: chatId,
      text: \`You said: \${text}\`
    })
  });

  res.json({ ok: true });
});

app.listen(3000);
\`\`\`

### Python (Flask)
\`\`\`python
from flask import Flask, request
import requests

app = Flask(__name__)
BOT_TOKEN = "your-token"
API_URL = "https://api.solafon.com/api/v1"

@app.route("/webhook", methods=["POST"])
def webhook():
    data = request.json
    message = data["message"]
    chat_id = message["from"]["id"]
    text = message["text"]

    requests.post(f"{API_URL}/bot/sendMessage",
        headers={"Authorization": f"Bearer {BOT_TOKEN}"},
        json={"chat_id": chat_id, "text": f"You said: {text}"})

    return {"ok": True}

app.run(port=3000)
\`\`\`

### Go
\`\`\`go
package main

import (
    "bytes"
    "encoding/json"
    "fmt"
    "net/http"
)

const botToken = "your-token"
const apiURL = "https://api.solafon.com/api/v1"

type Update struct {
    Message struct {
        From struct { ID int \`json:"id"\` } \`json:"from"\`
        Text string \`json:"text"\`
    } \`json:"message"\`
}

func webhook(w http.ResponseWriter, r *http.Request) {
    var update Update
    json.NewDecoder(r.Body).Decode(&update)

    body, _ := json.Marshal(map[string]interface{}{
        "chat_id": update.Message.From.ID,
        "text":    fmt.Sprintf("You said: %s", update.Message.Text),
    })

    req, _ := http.NewRequest("POST", apiURL+"/bot/sendMessage", bytes.NewBuffer(body))
    req.Header.Set("Authorization", "Bearer "+botToken)
    req.Header.Set("Content-Type", "application/json")
    http.DefaultClient.Do(req)

    json.NewEncoder(w).Encode(map[string]bool{"ok": true})
}

func main() {
    http.HandleFunc("/webhook", webhook)
    http.ListenAndServe(":3000", nil)
}
\`\`\`

## Polling vs Webhooks
| Feature | Polling | Webhooks |
|---------|---------|----------|
| Latency | Higher (poll interval) | Near-instant |
| Server Required | No | Yes (HTTPS) |
| Complexity | Simple | Moderate |
| Reliability | Manual retry | Auto-retry |
| Best For | Development, simple bots | Production, real-time apps |

## Troubleshooting
- **Not receiving updates**: Check webhook URL is HTTPS and accessible
- **401 errors**: Verify your API token
- **403 errors**: App must be approved (moderation status = "approved")
- **Webhook logs**: Check via GET /apps/:id/settings for recent webhook delivery attempts`,
  },

  "message-types": {
    title: "Message Types",
    content: `# Message Types

Solafon supports three message types: text, image, and button.

## Text Messages
\`\`\`json
{
  "chat_id": 1,
  "text": "Hello! This is a plain text message."
}
\`\`\`

## Image Messages
\`\`\`json
{
  "chat_id": 1,
  "text": "Check out this image!",
  "message_type": "image",
  "metadata": "{\\"image_url\\": \\"https://example.com/image.png\\"}"
}
\`\`\`

## Button Messages
\`\`\`json
{
  "chat_id": 1,
  "text": "Choose an option:",
  "message_type": "button",
  "metadata": "{\\"buttons\\": [{\\"text\\": \\"Option A\\", \\"action\\": \\"select_a\\"}, {\\"text\\": \\"Option B\\", \\"action\\": \\"select_b\\"}]}"
}
\`\`\`

Note: metadata is a JSON string, not a JSON object.`,
  },

  "mana-points": {
    title: "Mana Points System",
    content: `# Mana Points (MP)

Internal currency for the Solafon platform.

## Getting MP
- **Welcome bonus**: 500 MP for new users
- **Top-up**: Purchase via POST /mana/topup

## Pricing
| Package | Price |
|---------|-------|
| 100 MP | $2.30 |
| 500 MP | $11.50 |
| 1000 MP | $23.00 |

1 MP = $0.023

## Spending MP
- **Secret Login**: 50-150 MP for virtual phone numbers
- **Premium features**: Various costs

## API
\`\`\`
GET  /mana           // Balance + transaction history
POST /mana/topup     // Add MP ({"amount": 100})
\`\`\`

## Transaction Types
- \`topup\` — purchased MP
- \`purchase\` — spent on features
- \`reward\` — earned (welcome bonus, etc.)
- \`refund\` — returned MP`,
  },

  "api-overview": {
    title: "API Overview",
    content: `# API Overview

## Base URL
\`https://api.solafon.com/api/v1\`

## Response Format
**Success:**
\`\`\`json
{"data": {...}}
// or for lists:
{"items": [...], "total": 42}
\`\`\`

**Error:**
\`\`\`json
{"error": "Error message"}
\`\`\`

**Bot API:**
\`\`\`json
{"ok": true, "result": {...}}
{"ok": false, "error_code": 400, "description": "Bad request"}
\`\`\`

## HTTP Status Codes
- 200: Success
- 201: Created
- 400: Bad request
- 401: Unauthorized
- 403: Forbidden
- 404: Not found
- 429: Rate limited
- 500: Server error

## Rate Limits
- Authenticated: 100 requests/minute
- Unauthenticated: 20 requests/minute
- Headers: X-RateLimit-Limit, X-RateLimit-Remaining, X-RateLimit-Reset

## CORS
All origins allowed. Methods: GET, POST, PUT, PATCH, DELETE, OPTIONS.

## All Endpoints Summary

### Auth (Public)
POST /auth/email/request
POST /auth/email/verify

### Auth (Protected)
GET  /auth/me
POST /auth/logout

### Categories
GET /categories
GET /categories/:slug/apps

### Apps
GET    /apps
GET    /apps/search?q=
POST   /apps
GET    /apps/my
GET    /apps/:id
PUT    /apps/:id
DELETE /apps/:id

### Messages
GET  /apps/:id/messages
POST /apps/:id/messages

### App Settings
GET  /apps/:id/settings
POST /apps/:id/regenerate-token

### Commands (REST)
GET    /apps/:id/commands
POST   /apps/:id/commands
PUT    /apps/:id/commands/:cmdId
DELETE /apps/:id/commands/:cmdId

### Dev Studio
POST /apps/devstudio/message

### Profile
GET   /profile
PATCH /profile
PUT   /profile/settings

### Mana Points
GET  /mana
POST /mana/topup

### Secret Login
GET    /secret/numbers
POST   /secret/activate
GET    /secret/status
DELETE /secret/deactivate

### Bot API (Token Auth)
GET  /bot/getMe
POST /bot/sendMessage
GET  /bot/getUpdates
POST /bot/setWebhook
POST /bot/deleteWebhook
GET  /bot/getWebhookInfo
POST /bot/setMyCommands
GET  /bot/getMyCommands`,
  },

  "code-examples": {
    title: "Code Examples",
    content: `# Code Examples

## Echo Bot (Node.js — Polling)
\`\`\`javascript
const TOKEN = process.env.SOLAFON_BOT_TOKEN;
const API = 'https://api.solafon.com/api/v1';

async function poll() {
  while (true) {
    try {
      const res = await fetch(\`\${API}/bot/getUpdates\`, {
        headers: { 'Authorization': \`Bearer \${TOKEN}\` }
      });
      const data = await res.json();

      if (data.ok && data.result) {
        for (const msg of data.result) {
          await fetch(\`\${API}/bot/sendMessage\`, {
            method: 'POST',
            headers: {
              'Authorization': \`Bearer \${TOKEN}\`,
              'Content-Type': 'application/json'
            },
            body: JSON.stringify({
              chat_id: msg.from.id,
              text: \`Echo: \${msg.text}\`
            })
          });
        }
      }
    } catch (e) {
      console.error('Poll error:', e);
    }
    await new Promise(r => setTimeout(r, 2000));
  }
}

poll();
\`\`\`

## Echo Bot (Python — Polling)
\`\`\`python
import requests
import time
import os

TOKEN = os.environ["SOLAFON_BOT_TOKEN"]
API = "https://api.solafon.com/api/v1"
HEADERS = {"Authorization": f"Bearer {TOKEN}"}

while True:
    try:
        res = requests.get(f"{API}/bot/getUpdates", headers=HEADERS)
        data = res.json()

        if data.get("ok") and data.get("result"):
            for msg in data["result"]:
                requests.post(f"{API}/bot/sendMessage",
                    headers=HEADERS,
                    json={"chat_id": msg["from"]["id"], "text": f"Echo: {msg['text']}"})
    except Exception as e:
        print(f"Error: {e}")

    time.sleep(2)
\`\`\`

## AI Bot with OpenAI (Node.js)
\`\`\`javascript
import OpenAI from 'openai';

const TOKEN = process.env.SOLAFON_BOT_TOKEN;
const API = 'https://api.solafon.com/api/v1';
const openai = new OpenAI();

const conversations = new Map();

async function handleMessage(msg) {
  const chatId = msg.from.id;

  if (!conversations.has(chatId)) {
    conversations.set(chatId, [
      { role: 'system', content: 'You are a helpful Solana assistant.' }
    ]);
  }

  const history = conversations.get(chatId);
  history.push({ role: 'user', content: msg.text });

  const completion = await openai.chat.completions.create({
    model: 'gpt-4o-mini',
    messages: history
  });

  const reply = completion.choices[0].message.content;
  history.push({ role: 'assistant', content: reply });

  await fetch(\`\${API}/bot/sendMessage\`, {
    method: 'POST',
    headers: {
      'Authorization': \`Bearer \${TOKEN}\`,
      'Content-Type': 'application/json'
    },
    body: JSON.stringify({ chat_id: chatId, text: reply })
  });
}

// Polling loop
async function poll() {
  while (true) {
    const res = await fetch(\`\${API}/bot/getUpdates\`, {
      headers: { 'Authorization': \`Bearer \${TOKEN}\` }
    });
    const data = await res.json();
    if (data.ok && data.result) {
      for (const msg of data.result) await handleMessage(msg);
    }
    await new Promise(r => setTimeout(r, 2000));
  }
}

poll();
\`\`\`

## Webhook Server Template (Express)
\`\`\`javascript
import express from 'express';

const app = express();
app.use(express.json());

const TOKEN = process.env.SOLAFON_BOT_TOKEN;
const API = 'https://api.solafon.com/api/v1';

app.post('/webhook', async (req, res) => {
  const { message } = req.body;

  if (message.text.startsWith('/')) {
    // Handle commands
    const command = message.text.split(' ')[0].slice(1);
    switch (command) {
      case 'start':
        await sendMessage(message.from.id, 'Welcome to my bot!');
        break;
      case 'help':
        await sendMessage(message.from.id, 'Commands: /start, /help, /info');
        break;
      default:
        await sendMessage(message.from.id, 'Unknown command. Try /help');
    }
  } else {
    // Handle regular messages
    await sendMessage(message.from.id, \`You said: \${message.text}\`);
  }

  res.json({ ok: true });
});

async function sendMessage(chatId, text) {
  await fetch(\`\${API}/bot/sendMessage\`, {
    method: 'POST',
    headers: {
      'Authorization': \`Bearer \${TOKEN}\`,
      'Content-Type': 'application/json'
    },
    body: JSON.stringify({ chat_id: chatId, text })
  });
}

app.listen(3000, () => console.log('Webhook server running on :3000'));
\`\`\``,
  },

  "project-architecture": {
    title: "Project Architecture",
    content: `# Solafon Backend Architecture

## Tech Stack
- **Language**: Go 1.21+
- **Web Framework**: Fiber v2
- **ORM**: GORM
- **Database**: PostgreSQL 15
- **Auth**: JWT (HS256) + API Tokens

## Project Structure
\`\`\`
solanafon-back/
├── cmd/
│   ├── server/main.go         # Entry point
│   └── seed/main.go           # Database seeder
├── internal/
│   ├── config/config.go       # Environment configuration
│   ├── database/database.go   # DB connection + migrations
│   ├── handlers/
│   │   ├── auth.go            # OTP + JWT auth
│   │   ├── miniapp.go         # App CRUD + messages
│   │   ├── devstudio.go       # Chat-based app management
│   │   ├── bot.go             # Bot API (Telegram-style)
│   │   ├── profile.go         # User profile + Mana Points
│   │   └── secret.go          # Secret Login
│   ├── middleware/auth.go     # JWT middleware
│   ├── models/
│   │   ├── user.go            # User, OTP, ManaTransaction
│   │   ├── miniapp.go         # App, Category, Messages, Commands
│   │   └── secret.go          # SecretNumber, SecretAccess
│   ├── routes/routes.go       # All route definitions
│   └── utils/
│       ├── jwt.go             # JWT generation/validation
│       ├── otp.go             # OTP generation
│       └── email.go           # SMTP email sending
├── docs/                      # Documentation (GitBook)
├── mcp/                       # MCP server for AI tools
├── Makefile                   # Dev commands
├── Dockerfile                 # Multi-stage build
├── docker-compose.yml         # PostgreSQL service
└── go.mod                     # Dependencies
\`\`\`

## Key Patterns
1. **Handler Pattern**: Each feature has its own handler struct with DB + config
2. **State Machine**: Dev Studio uses conversation states for multi-step workflows
3. **Dual Auth**: JWT for users, API tokens for bots
4. **Telegram-compatible Bot API**: Same message format and endpoint naming
5. **Webhook + Polling**: Both push and pull message delivery
6. **Soft Deletes**: GORM DeletedAt for data preservation
7. **Auto-migration**: Schema managed by GORM on startup

## Database Tables (13)
User, OTP, ManaTransaction, Category, MiniApp, AppUser, AppMessage,
BotCommand, WebhookLog, ConversationState, SecretNumber, SecretAccess

## Environment Variables
PORT, DATABASE_URL, JWT_SECRET, SMTP_HOST, SMTP_PORT,
SMTP_USER, SMTP_PASSWORD, OTP_EXPIRY_MINUTES,
RATE_LIMIT_REQUESTS, RATE_LIMIT_WINDOW`,
  },

  "secret-login": {
    title: "Secret Login",
    content: `# Secret Login

Privacy feature that provides virtual phone numbers for anonymous access.

## API Endpoints

### List Available Numbers
\`\`\`
GET /secret/numbers
\`\`\`
Returns available virtual numbers with pricing.
Format: +999 (XXX) XXX-XX-XX
Standard: 50 MP, Premium: 100-150 MP.

### Activate Secret Access
\`\`\`
POST /secret/activate
{"number_id": 1}
\`\`\`
Deducts MP from balance, activates access to secret apps.

### Check Status
\`\`\`
GET /secret/status
\`\`\`

### Deactivate
\`\`\`
DELETE /secret/deactivate
\`\`\`

## Benefits
- Anonymous identity within the platform
- Access to secret/hidden apps (is_secret: true)
- Virtual number for privacy
- No personal phone number exposure`,
  },

  "mcp-setup": {
    title: "MCP Server Setup",
    content: `# Solafon MCP Server Setup

The Solafon MCP Server gives your AI assistant (Claude, Cursor, etc.) direct access to Solafon developer tools.

## Installation

\`\`\`bash
cd mcp
npm install
npm run build
\`\`\`

## Claude Desktop Configuration
Add to ~/Library/Application Support/Claude/claude_desktop_config.json (macOS)
or %APPDATA%/Claude/claude_desktop_config.json (Windows):

\`\`\`json
{
  "mcpServers": {
    "solafon": {
      "command": "node",
      "args": ["/path/to/solanafon-back/mcp/dist/index.js"],
      "env": {
        "SOLAFON_BOT_TOKEN": "your-api-token-here",
        "SOLAFON_API_URL": "https://api.solafon.com/api/v1"
      }
    }
  }
}
\`\`\`

## Cursor Configuration
Add to .cursor/mcp.json in your project:

\`\`\`json
{
  "mcpServers": {
    "solafon": {
      "command": "node",
      "args": ["./mcp/dist/index.js"],
      "env": {
        "SOLAFON_BOT_TOKEN": "your-api-token-here"
      }
    }
  }
}
\`\`\`

## VS Code (Claude Code) Configuration
Add to .vscode/mcp.json:

\`\`\`json
{
  "servers": {
    "solafon": {
      "command": "node",
      "args": ["./mcp/dist/index.js"],
      "env": {
        "SOLAFON_BOT_TOKEN": "your-api-token-here"
      }
    }
  }
}
\`\`\`

## Available Tools
After connecting, your AI assistant gets these tools:

### Bot API Tools
- \`solafon_send_message\` — Send message to user
- \`solafon_get_updates\` — Get pending messages
- \`solafon_get_bot_info\` — Get your bot info
- \`solafon_set_webhook\` — Configure webhook
- \`solafon_delete_webhook\` — Remove webhook
- \`solafon_get_webhook_info\` — Check webhook status
- \`solafon_set_commands\` — Define bot commands
- \`solafon_get_commands\` — List bot commands

### Development Tools
- \`solafon_api_request\` — Make any API request
- \`solafon_health_check\` — Check API status
- \`solafon_scaffold_bot\` — Generate bot code template

### Documentation
- \`solafon_read_docs\` — Read platform documentation
- \`solafon_search_docs\` — Search documentation

## Environment Variables
| Variable | Required | Description |
|----------|----------|-------------|
| SOLAFON_BOT_TOKEN | Yes* | Your bot API token |
| SOLAFON_API_URL | No | API URL (default: https://api.solafon.com/api/v1) |

*Required for Bot API tools. Documentation tools work without a token.`,
  },
};

// All available doc topics for listing
const DOC_TOPICS = Object.entries(DOCS).map(([key, doc]) => ({
  key,
  title: doc.title,
}));

// ---------------------------------------------------------------------------
// Code Scaffold Templates
// ---------------------------------------------------------------------------

const SCAFFOLDS: Record<string, { title: string; description: string; code: string; language: string }> = {
  "node-polling": {
    title: "Node.js Polling Bot",
    description: "Simple polling bot using Node.js with message handling",
    language: "javascript",
    code: `// Solafon Polling Bot — Node.js
// Install: npm init -y
// Run: SOLAFON_BOT_TOKEN=your-token node bot.js

const TOKEN = process.env.SOLAFON_BOT_TOKEN;
const API = 'https://api.solafon.com/api/v1';

async function sendMessage(chatId, text, options = {}) {
  const body = { chat_id: chatId, text, ...options };
  const res = await fetch(\`\${API}/bot/sendMessage\`, {
    method: 'POST',
    headers: {
      'Authorization': \`Bearer \${TOKEN}\`,
      'Content-Type': 'application/json'
    },
    body: JSON.stringify(body)
  });
  return res.json();
}

async function handleMessage(msg) {
  const chatId = msg.from.id;
  const text = msg.text || '';

  if (text === '/start') {
    await sendMessage(chatId, 'Welcome! I am your bot. Send me anything!');
  } else if (text === '/help') {
    await sendMessage(chatId, 'Commands:\\n/start - Start\\n/help - Show help');
  } else {
    await sendMessage(chatId, \`You said: \${text}\`);
  }
}

async function poll() {
  console.log('Bot started polling...');
  while (true) {
    try {
      const res = await fetch(\`\${API}/bot/getUpdates\`, {
        headers: { 'Authorization': \`Bearer \${TOKEN}\` }
      });
      const data = await res.json();
      if (data.ok && data.result?.length > 0) {
        for (const msg of data.result) {
          await handleMessage(msg);
        }
      }
    } catch (error) {
      console.error('Polling error:', error.message);
    }
    await new Promise(resolve => setTimeout(resolve, 2000));
  }
}

poll();
`,
  },

  "node-webhook": {
    title: "Node.js Webhook Bot",
    description: "Express.js webhook bot with command handling",
    language: "javascript",
    code: `// Solafon Webhook Bot — Node.js + Express
// Install: npm install express
// Run: SOLAFON_BOT_TOKEN=your-token node bot.js

import express from 'express';

const app = express();
app.use(express.json());

const TOKEN = process.env.SOLAFON_BOT_TOKEN;
const API = 'https://api.solafon.com/api/v1';

async function sendMessage(chatId, text, options = {}) {
  await fetch(\`\${API}/bot/sendMessage\`, {
    method: 'POST',
    headers: {
      'Authorization': \`Bearer \${TOKEN}\`,
      'Content-Type': 'application/json'
    },
    body: JSON.stringify({ chat_id: chatId, text, ...options })
  });
}

// Webhook endpoint
app.post('/webhook', async (req, res) => {
  try {
    const { message } = req.body;
    const chatId = message.from.id;
    const text = message.text || '';

    // Command handling
    if (text.startsWith('/')) {
      const command = text.split(' ')[0].slice(1);
      const args = text.split(' ').slice(1).join(' ');

      switch (command) {
        case 'start':
          await sendMessage(chatId, 'Welcome to my bot!');
          break;
        case 'help':
          await sendMessage(chatId, 'Available commands: /start, /help');
          break;
        default:
          await sendMessage(chatId, \`Unknown command: /\${command}\`);
      }
    } else {
      // Regular message handling
      await sendMessage(chatId, \`Received: \${text}\`);
    }

    res.json({ ok: true });
  } catch (error) {
    console.error('Webhook error:', error);
    res.status(500).json({ ok: false });
  }
});

const PORT = process.env.PORT || 3000;
app.listen(PORT, () => {
  console.log(\`Webhook server running on port \${PORT}\`);
  console.log('Set webhook URL: POST /bot/setWebhook {"url": "https://your-domain.com/webhook"}');
});
`,
  },

  "python-polling": {
    title: "Python Polling Bot",
    description: "Simple Python polling bot with requests",
    language: "python",
    code: `# Solafon Polling Bot — Python
# Install: pip install requests
# Run: SOLAFON_BOT_TOKEN=your-token python bot.py

import os
import time
import requests

TOKEN = os.environ["SOLAFON_BOT_TOKEN"]
API = "https://api.solafon.com/api/v1"
HEADERS = {"Authorization": f"Bearer {TOKEN}", "Content-Type": "application/json"}


def send_message(chat_id: int, text: str, **kwargs):
    payload = {"chat_id": chat_id, "text": text, **kwargs}
    return requests.post(f"{API}/bot/sendMessage", headers=HEADERS, json=payload).json()


def handle_message(msg: dict):
    chat_id = msg["from"]["id"]
    text = msg.get("text", "")

    if text == "/start":
        send_message(chat_id, "Welcome! I am your Solafon bot.")
    elif text == "/help":
        send_message(chat_id, "Commands:\\n/start - Start\\n/help - Help")
    else:
        send_message(chat_id, f"You said: {text}")


def poll():
    print("Bot started polling...")
    while True:
        try:
            res = requests.get(f"{API}/bot/getUpdates", headers=HEADERS)
            data = res.json()
            if data.get("ok") and data.get("result"):
                for msg in data["result"]:
                    handle_message(msg)
        except Exception as e:
            print(f"Error: {e}")
        time.sleep(2)


if __name__ == "__main__":
    poll()
`,
  },

  "python-webhook": {
    title: "Python Webhook Bot",
    description: "Flask webhook bot with command handling",
    language: "python",
    code: `# Solafon Webhook Bot — Python + Flask
# Install: pip install flask requests
# Run: SOLAFON_BOT_TOKEN=your-token python bot.py

import os
import requests
from flask import Flask, request, jsonify

app = Flask(__name__)

TOKEN = os.environ["SOLAFON_BOT_TOKEN"]
API = "https://api.solafon.com/api/v1"
HEADERS = {"Authorization": f"Bearer {TOKEN}", "Content-Type": "application/json"}


def send_message(chat_id: int, text: str, **kwargs):
    payload = {"chat_id": chat_id, "text": text, **kwargs}
    requests.post(f"{API}/bot/sendMessage", headers=HEADERS, json=payload)


@app.route("/webhook", methods=["POST"])
def webhook():
    data = request.json
    message = data.get("message", {})
    chat_id = message["from"]["id"]
    text = message.get("text", "")

    if text.startswith("/"):
        command = text.split()[0][1:]
        if command == "start":
            send_message(chat_id, "Welcome to my bot!")
        elif command == "help":
            send_message(chat_id, "Commands: /start, /help")
        else:
            send_message(chat_id, f"Unknown command: /{command}")
    else:
        send_message(chat_id, f"Received: {text}")

    return jsonify({"ok": True})


if __name__ == "__main__":
    app.run(host="0.0.0.0", port=int(os.environ.get("PORT", 3000)))
`,
  },

  "go-webhook": {
    title: "Go Webhook Bot",
    description: "Go webhook bot using net/http",
    language: "go",
    code: `// Solafon Webhook Bot — Go
// Run: SOLAFON_BOT_TOKEN=your-token go run main.go

package main

import (
\t"bytes"
\t"encoding/json"
\t"fmt"
\t"log"
\t"net/http"
\t"os"
\t"strings"
)

var (
\tbotToken = os.Getenv("SOLAFON_BOT_TOKEN")
\tapiURL   = "https://api.solafon.com/api/v1"
)

type Update struct {
\tUpdateID int \`json:"update_id"\`
\tMessage  struct {
\t\tMessageID int \`json:"message_id"\`
\t\tFrom      struct {
\t\t\tID       int    \`json:"id"\`
\t\t\tEmail    string \`json:"email"\`
\t\t\tName     string \`json:"name"\`
\t\t\tLanguage string \`json:"language"\`
\t\t} \`json:"from"\`
\t\tChat struct {
\t\t\tID   int    \`json:"id"\`
\t\t\tType string \`json:"type"\`
\t\t} \`json:"chat"\`
\t\tDate int    \`json:"date"\`
\t\tText string \`json:"text"\`
\t} \`json:"message"\`
}

func sendMessage(chatID int, text string) error {
\tbody, _ := json.Marshal(map[string]interface{}{
\t\t"chat_id": chatID,
\t\t"text":    text,
\t})
\treq, _ := http.NewRequest("POST", apiURL+"/bot/sendMessage", bytes.NewBuffer(body))
\treq.Header.Set("Authorization", "Bearer "+botToken)
\treq.Header.Set("Content-Type", "application/json")
\t_, err := http.DefaultClient.Do(req)
\treturn err
}

func webhookHandler(w http.ResponseWriter, r *http.Request) {
\tvar update Update
\tif err := json.NewDecoder(r.Body).Decode(&update); err != nil {
\t\thttp.Error(w, "Bad request", 400)
\t\treturn
\t}

\tchatID := update.Message.From.ID
\ttext := update.Message.Text

\tif strings.HasPrefix(text, "/") {
\t\tcmd := strings.Split(text, " ")[0][1:]
\t\tswitch cmd {
\t\tcase "start":
\t\t\tsendMessage(chatID, "Welcome to my Solafon bot!")
\t\tcase "help":
\t\t\tsendMessage(chatID, "Commands: /start, /help")
\t\tdefault:
\t\t\tsendMessage(chatID, fmt.Sprintf("Unknown command: /%s", cmd))
\t\t}
\t} else {
\t\tsendMessage(chatID, fmt.Sprintf("You said: %s", text))
\t}

\tjson.NewEncoder(w).Encode(map[string]bool{"ok": true})
}

func main() {
\thttp.HandleFunc("/webhook", webhookHandler)
\tport := os.Getenv("PORT")
\tif port == "" {
\t\tport = "3000"
\t}
\tlog.Printf("Webhook server running on :%s", port)
\tlog.Fatal(http.ListenAndServe(":"+port, nil))
}
`,
  },

  "node-ai-bot": {
    title: "AI Bot with OpenAI (Node.js)",
    description: "AI-powered bot using OpenAI GPT with conversation memory",
    language: "javascript",
    code: `// Solafon AI Bot — Node.js + OpenAI
// Install: npm install openai
// Run: SOLAFON_BOT_TOKEN=x OPENAI_API_KEY=x node bot.js

import OpenAI from 'openai';

const TOKEN = process.env.SOLAFON_BOT_TOKEN;
const API = 'https://api.solafon.com/api/v1';
const openai = new OpenAI();

// Conversation memory per user
const conversations = new Map();

async function sendMessage(chatId, text) {
  await fetch(\`\${API}/bot/sendMessage\`, {
    method: 'POST',
    headers: {
      'Authorization': \`Bearer \${TOKEN}\`,
      'Content-Type': 'application/json'
    },
    body: JSON.stringify({ chat_id: chatId, text })
  });
}

async function handleMessage(msg) {
  const chatId = msg.from.id;
  const text = msg.text || '';

  // Reset conversation on /start
  if (text === '/start') {
    conversations.delete(chatId);
    await sendMessage(chatId, 'Hello! I am an AI assistant. Ask me anything!');
    return;
  }

  // Get or create conversation history
  if (!conversations.has(chatId)) {
    conversations.set(chatId, [{
      role: 'system',
      content: 'You are a helpful assistant on the Solafon platform. Be concise and friendly.'
    }]);
  }

  const history = conversations.get(chatId);
  history.push({ role: 'user', content: text });

  // Keep last 20 messages to avoid token limits
  if (history.length > 21) {
    history.splice(1, history.length - 21);
  }

  try {
    const completion = await openai.chat.completions.create({
      model: 'gpt-4o-mini',
      messages: history,
      max_tokens: 500
    });

    const reply = completion.choices[0].message.content;
    history.push({ role: 'assistant', content: reply });

    await sendMessage(chatId, reply);
  } catch (error) {
    console.error('OpenAI error:', error);
    await sendMessage(chatId, 'Sorry, I encountered an error. Please try again.');
  }
}

async function poll() {
  console.log('AI Bot started...');
  while (true) {
    try {
      const res = await fetch(\`\${API}/bot/getUpdates\`, {
        headers: { 'Authorization': \`Bearer \${TOKEN}\` }
      });
      const data = await res.json();
      if (data.ok && data.result?.length > 0) {
        for (const msg of data.result) {
          await handleMessage(msg);
        }
      }
    } catch (error) {
      console.error('Poll error:', error.message);
    }
    await new Promise(r => setTimeout(r, 2000));
  }
}

poll();
`,
  },
};

const SCAFFOLD_LIST = Object.entries(SCAFFOLDS).map(([key, s]) => ({
  key,
  title: s.title,
  description: s.description,
  language: s.language,
}));

// ---------------------------------------------------------------------------
// Server Setup
// ---------------------------------------------------------------------------

const server = new McpServer({
  name: "solafon",
  version: "1.0.0",
});

// ---------------------------------------------------------------------------
// Tools: Documentation
// ---------------------------------------------------------------------------

server.tool(
  "solafon_read_docs",
  "Read Solafon platform documentation by topic. Returns detailed documentation about a specific topic.",
  {
    topic: z.enum(Object.keys(DOCS) as [string, ...string[]]).describe(
      `Documentation topic. Available: ${DOC_TOPICS.map((t) => `${t.key} (${t.title})`).join(", ")}`
    ),
  },
  async ({ topic }) => {
    const doc = DOCS[topic];
    if (!doc) {
      return {
        content: [
          {
            type: "text" as const,
            text: `Topic "${topic}" not found. Available topics: ${DOC_TOPICS.map((t) => t.key).join(", ")}`,
          },
        ],
      };
    }
    return {
      content: [{ type: "text" as const, text: doc.content }],
    };
  }
);

server.tool(
  "solafon_search_docs",
  "Search across all Solafon documentation for a keyword or phrase. Returns matching sections.",
  {
    query: z.string().describe("Search query (keyword or phrase)"),
  },
  async ({ query }) => {
    const q = query.toLowerCase();
    const results: string[] = [];

    for (const [key, doc] of Object.entries(DOCS)) {
      if (
        doc.content.toLowerCase().includes(q) ||
        doc.title.toLowerCase().includes(q)
      ) {
        // Find the matching lines with context
        const lines = doc.content.split("\n");
        const matchingLines: string[] = [];
        for (let i = 0; i < lines.length; i++) {
          if (lines[i].toLowerCase().includes(q)) {
            const start = Math.max(0, i - 1);
            const end = Math.min(lines.length, i + 3);
            matchingLines.push(lines.slice(start, end).join("\n"));
          }
        }
        results.push(
          `## ${doc.title} (topic: "${key}")\n${matchingLines.slice(0, 3).join("\n---\n")}`
        );
      }
    }

    if (results.length === 0) {
      return {
        content: [
          {
            type: "text" as const,
            text: `No results found for "${query}". Try different keywords or use solafon_read_docs to browse topics: ${DOC_TOPICS.map((t) => t.key).join(", ")}`,
          },
        ],
      };
    }

    return {
      content: [
        {
          type: "text" as const,
          text: `Found ${results.length} matching doc(s) for "${query}":\n\n${results.join("\n\n---\n\n")}`,
        },
      ],
    };
  }
);

server.tool(
  "solafon_list_docs",
  "List all available Solafon documentation topics.",
  {},
  async () => {
    const list = DOC_TOPICS.map((t) => `- **${t.key}**: ${t.title}`).join("\n");
    return {
      content: [
        {
          type: "text" as const,
          text: `# Available Documentation Topics\n\n${list}\n\nUse \`solafon_read_docs\` with a topic key to read the full documentation.`,
        },
      ],
    };
  }
);

// ---------------------------------------------------------------------------
// Tools: Bot API
// ---------------------------------------------------------------------------

server.tool(
  "solafon_get_bot_info",
  "Get information about your Solafon bot (id, title, username, webhook status). Requires SOLAFON_BOT_TOKEN.",
  {},
  async () => {
    const data = await apiRequest("GET", "/bot/getMe");
    return {
      content: [{ type: "text" as const, text: JSON.stringify(data, null, 2) }],
    };
  }
);

server.tool(
  "solafon_send_message",
  "Send a message from your bot to a user on Solafon. Supports text, image, and button message types.",
  {
    chat_id: z.number().describe("User ID to send the message to"),
    text: z.string().describe("Message text content"),
    message_type: z
      .enum(["text", "image", "button"])
      .optional()
      .describe("Message type: text (default), image, or button"),
    metadata: z
      .string()
      .optional()
      .describe(
        'JSON string with extra data. For images: {"image_url": "..."}, for buttons: {"buttons": [{"text": "...", "action": "..."}]}'
      ),
  },
  async ({ chat_id, text, message_type, metadata }) => {
    const body: Record<string, unknown> = { chat_id, text };
    if (message_type) body.message_type = message_type;
    if (metadata) body.metadata = metadata;

    const data = await apiRequest("POST", "/bot/sendMessage", body);
    return {
      content: [{ type: "text" as const, text: JSON.stringify(data, null, 2) }],
    };
  }
);

server.tool(
  "solafon_get_updates",
  "Get pending (unread) messages from users via long-polling. Messages are marked as read after retrieval.",
  {},
  async () => {
    const data = await apiRequest("GET", "/bot/getUpdates");
    return {
      content: [{ type: "text" as const, text: JSON.stringify(data, null, 2) }],
    };
  }
);

server.tool(
  "solafon_set_webhook",
  "Set a webhook URL to receive messages in real-time. URL must be HTTPS.",
  {
    url: z.string().url().describe("HTTPS webhook URL to receive message updates"),
  },
  async ({ url }) => {
    const data = await apiRequest("POST", "/bot/setWebhook", { url });
    return {
      content: [{ type: "text" as const, text: JSON.stringify(data, null, 2) }],
    };
  }
);

server.tool(
  "solafon_delete_webhook",
  "Remove the currently configured webhook. Switch back to polling mode.",
  {},
  async () => {
    const data = await apiRequest("POST", "/bot/deleteWebhook");
    return {
      content: [{ type: "text" as const, text: JSON.stringify(data, null, 2) }],
    };
  }
);

server.tool(
  "solafon_get_webhook_info",
  "Get current webhook configuration including URL and pending update count.",
  {},
  async () => {
    const data = await apiRequest("GET", "/bot/getWebhookInfo");
    return {
      content: [{ type: "text" as const, text: JSON.stringify(data, null, 2) }],
    };
  }
);

server.tool(
  "solafon_set_commands",
  "Define bot commands. Commands with a 'response' field will auto-reply without hitting your server.",
  {
    commands: z
      .array(
        z.object({
          command: z.string().describe("Command name without / prefix (e.g. 'help')"),
          description: z.string().describe("Short description of the command"),
          response: z
            .string()
            .optional()
            .describe("Auto-reply text. If set, Solafon responds automatically without webhook/polling."),
        })
      )
      .describe("Array of command objects"),
  },
  async ({ commands }) => {
    const data = await apiRequest("POST", "/bot/setMyCommands", { commands });
    return {
      content: [{ type: "text" as const, text: JSON.stringify(data, null, 2) }],
    };
  }
);

server.tool(
  "solafon_get_commands",
  "Get all defined bot commands.",
  {},
  async () => {
    const data = await apiRequest("GET", "/bot/getMyCommands");
    return {
      content: [{ type: "text" as const, text: JSON.stringify(data, null, 2) }],
    };
  }
);

// ---------------------------------------------------------------------------
// Tools: Development Helpers
// ---------------------------------------------------------------------------

server.tool(
  "solafon_api_request",
  "Make a custom API request to any Solafon endpoint. Use for endpoints not covered by specific tools.",
  {
    method: z.enum(["GET", "POST", "PUT", "PATCH", "DELETE"]).describe("HTTP method"),
    path: z.string().describe("API path starting with / (e.g. /apps, /bot/getMe)"),
    body: z
      .record(z.string(), z.unknown())
      .optional()
      .describe("Request body as JSON object (for POST/PUT/PATCH)"),
    token: z
      .string()
      .optional()
      .describe("Override token. Defaults to SOLAFON_BOT_TOKEN env var."),
  },
  async ({ method, path, body, token }) => {
    const data = await apiRequest(method, path, body as Record<string, unknown>, token);
    return {
      content: [{ type: "text" as const, text: JSON.stringify(data, null, 2) }],
    };
  }
);

server.tool(
  "solafon_health_check",
  "Check if the Solafon API is online and responding.",
  {},
  async () => {
    try {
      const baseUrl = API_BASE_URL.replace("/api/v1", "");
      const response = await fetch(`${baseUrl}/health`);
      const data = await response.json();
      return {
        content: [
          {
            type: "text" as const,
            text: `API Status: ${response.ok ? "Online" : "Error"}\n${JSON.stringify(data, null, 2)}`,
          },
        ],
      };
    } catch (error) {
      return {
        content: [
          {
            type: "text" as const,
            text: `API Status: Offline\nError: ${error instanceof Error ? error.message : String(error)}`,
          },
        ],
      };
    }
  }
);

server.tool(
  "solafon_scaffold_bot",
  "Generate a bot code template for the Solafon platform. Choose language and architecture.",
  {
    template: z
      .enum(Object.keys(SCAFFOLDS) as [string, ...string[]])
      .describe(
        `Template name. Available: ${SCAFFOLD_LIST.map((s) => `${s.key} (${s.title})`).join(", ")}`
      ),
  },
  async ({ template }) => {
    const scaffold = SCAFFOLDS[template];
    if (!scaffold) {
      return {
        content: [
          {
            type: "text" as const,
            text: `Template "${template}" not found. Available: ${SCAFFOLD_LIST.map((s) => s.key).join(", ")}`,
          },
        ],
      };
    }
    return {
      content: [
        {
          type: "text" as const,
          text: `# ${scaffold.title}\n\n${scaffold.description}\n\n\`\`\`${scaffold.language}\n${scaffold.code}\`\`\``,
        },
      ],
    };
  }
);

server.tool(
  "solafon_list_templates",
  "List all available bot code templates for scaffolding.",
  {},
  async () => {
    const list = SCAFFOLD_LIST.map(
      (s) => `- **${s.key}** (${s.language}): ${s.title} — ${s.description}`
    ).join("\n");
    return {
      content: [
        {
          type: "text" as const,
          text: `# Available Bot Templates\n\n${list}\n\nUse \`solafon_scaffold_bot\` with a template key to generate the code.`,
        },
      ],
    };
  }
);

// ---------------------------------------------------------------------------
// Resources: Documentation Pages
// ---------------------------------------------------------------------------

for (const [key, doc] of Object.entries(DOCS)) {
  server.resource(
    `docs-${key}`,
    `solafon://docs/${key}`,
    { mimeType: "text/markdown", description: doc.title },
    async () => ({
      contents: [
        {
          uri: `solafon://docs/${key}`,
          mimeType: "text/markdown" as const,
          text: doc.content,
        },
      ],
    })
  );
}

// Resource: Full API reference
server.resource(
  "docs-full",
  "solafon://docs/full",
  { mimeType: "text/markdown", description: "Complete Solafon documentation" },
  async () => {
    const fullDocs = Object.values(DOCS)
      .map((d) => d.content)
      .join("\n\n---\n\n");
    return {
      contents: [
        {
          uri: "solafon://docs/full",
          mimeType: "text/markdown" as const,
          text: fullDocs,
        },
      ],
    };
  }
);

// Resource: Bot templates
for (const [key, scaffold] of Object.entries(SCAFFOLDS)) {
  server.resource(
    `template-${key}`,
    `solafon://templates/${key}`,
    { mimeType: "text/plain", description: scaffold.title },
    async () => ({
      contents: [
        {
          uri: `solafon://templates/${key}`,
          mimeType: "text/plain" as const,
          text: scaffold.code,
        },
      ],
    })
  );
}

// ---------------------------------------------------------------------------
// Prompts
// ---------------------------------------------------------------------------

server.prompt(
  "create-solafon-bot",
  "Step-by-step guide to create a new bot on the Solafon platform",
  {
    language: z
      .enum(["javascript", "python", "go"])
      .optional()
      .describe("Programming language for the bot"),
    architecture: z
      .enum(["polling", "webhook"])
      .optional()
      .describe("Message delivery: polling or webhook"),
  },
  async ({ language, architecture }) => ({
    messages: [
      {
        role: "user" as const,
        content: {
          type: "text" as const,
          text: `Help me create a new Solafon bot using ${language || "any language"} with ${architecture || "your recommended"} architecture.

Here's what I need:
1. First, use solafon_read_docs with topic "quick-start" to understand the setup process
2. Use solafon_scaffold_bot to generate a starter template
3. Set up bot commands using solafon_set_commands
4. Configure ${architecture === "webhook" ? "a webhook using solafon_set_webhook" : "polling to receive messages"}
5. Test the bot by sending a test message

Guide me through each step with explanations.`,
        },
      },
    ],
  })
);

server.prompt(
  "debug-solafon-bot",
  "Debug issues with your Solafon bot",
  {
    issue: z.string().optional().describe("Description of the issue"),
  },
  async ({ issue }) => ({
    messages: [
      {
        role: "user" as const,
        content: {
          type: "text" as const,
          text: `Help me debug my Solafon bot. ${issue ? `Issue: ${issue}` : ""}

Please:
1. Check if the API is online using solafon_health_check
2. Verify my bot token with solafon_get_bot_info
3. Check webhook configuration with solafon_get_webhook_info
4. Check for pending messages with solafon_get_updates
5. Review my bot commands with solafon_get_commands

Based on the results, diagnose the issue and suggest fixes.`,
        },
      },
    ],
  })
);

server.prompt(
  "solafon-api-explorer",
  "Explore the Solafon API interactively",
  {},
  async () => ({
    messages: [
      {
        role: "user" as const,
        content: {
          type: "text" as const,
          text: `I want to explore the Solafon API. Please:
1. Show me the available documentation topics using solafon_list_docs
2. Let me know what tools are available for interacting with the API
3. Show me available bot templates using solafon_list_templates

Then ask me what I'd like to do — create a bot, explore endpoints, or learn about a specific feature.`,
        },
      },
    ],
  })
);

// ---------------------------------------------------------------------------
// Start Server
// ---------------------------------------------------------------------------

async function main() {
  const transport = new StdioServerTransport();
  await server.connect(transport);
}

main().catch((error) => {
  console.error("Fatal error:", error);
  process.exit(1);
});

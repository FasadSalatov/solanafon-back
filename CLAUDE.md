# Solafon Backend — AI Development Guide

## Project Overview
Solafon is a platform for creating and deploying mini-applications (mini-apps) that users interact with through a chat-based interface. This is the Go backend.

## Tech Stack
- **Language**: Go 1.21+
- **Web Framework**: Fiber v2 (Express-like)
- **ORM**: GORM (PostgreSQL)
- **Database**: PostgreSQL 15
- **Auth**: JWT (HS256) for users, API Tokens for bots

## Project Structure
```
cmd/server/main.go         — Entry point, Fiber setup, middleware
cmd/seed/main.go           — DB seeder (categories, dev studio, secret numbers)
internal/
  config/config.go         — Env vars loading
  database/database.go     — PostgreSQL connection, auto-migrate
  handlers/
    auth.go                — Email OTP, JWT, user creation (500 MP bonus)
    miniapp.go             — App CRUD, messages, webhooks, commands
    devstudio.go           — Chat-based state machine for app management
    bot.go                 — Telegram-style Bot API (token auth)
    profile.go             — User profile, Mana Points, settings
    secret.go              — Secret Login, virtual numbers
  middleware/auth.go       — JWT validation middleware
  models/
    user.go                — User, OTP, ManaTransaction
    miniapp.go             — Category, MiniApp, AppUser, AppMessage, BotCommand, WebhookLog, ConversationState
    secret.go              — SecretNumber, SecretAccess
  routes/routes.go         — All API routes
  utils/
    jwt.go                 — JWT gen/validate
    otp.go                 — Crypto-secure OTP generation
    email.go               — SMTP email sending
docs/                      — GitBook documentation
mcp/                       — MCP server for AI dev tools
```

## Key Patterns
1. **Handler structs** with `db *gorm.DB` and `cfg *config.Config`
2. **Two auth systems**: JWT (user endpoints) and API Token (bot endpoints)
3. **State machine** in devstudio.go for multi-step conversations
4. **Telegram-compatible Bot API** at `/api/v1/bot/*`
5. **Soft deletes** via GORM `DeletedAt`
6. **Auto-migration** on server start

## API Structure
- `/api/v1/auth/*` — Email OTP authentication
- `/api/v1/apps/*` — Mini-app CRUD + messages
- `/api/v1/apps/devstudio/message` — Dev Studio chat interface
- `/api/v1/categories/*` — App categories
- `/api/v1/profile/*` — User profile
- `/api/v1/mana/*` — Mana Points economy
- `/api/v1/secret/*` — Secret Login
- `/api/v1/bot/*` — Bot API (token auth, not JWT)

## Development Commands
```bash
make run          # Start server
make dev          # Hot reload (Air)
make seed         # Seed database
make build        # Build binary
make docker-up    # Start PostgreSQL
make test         # Run tests
```

## Database
13 tables auto-migrated by GORM. Key models: User, MiniApp, AppMessage, BotCommand, Category, ManaTransaction, ConversationState.

## Environment Variables
PORT, DATABASE_URL, JWT_SECRET, SMTP_HOST, SMTP_PORT, SMTP_USER, SMTP_PASSWORD, OTP_EXPIRY_MINUTES, RATE_LIMIT_REQUESTS, RATE_LIMIT_WINDOW

## Conventions
- Handlers return Fiber JSON responses: `c.JSON(fiber.Map{...})`
- Error format: `{"error": "message"}`
- Bot API format: `{"ok": true/false, "result": ...}`
- App usernames must end with "app"
- API tokens are 64-char hex strings
- New users get 500 MP welcome bonus
- Apps require moderation: pending → approved/rejected

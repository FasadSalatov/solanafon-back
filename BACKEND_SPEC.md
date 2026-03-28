# Solafon — Full Backend API Specification

> Complete backend specification for the Solafon Android app.
> Base URL: `https://api.solafon.com`
> Generated: 2026-03-02

---

## Table of Contents

1. [Architecture Overview](#1-architecture-overview)
2. [Authentication](#2-authentication)
3. [User Management](#3-user-management)
4. [Apps Marketplace](#4-apps-marketplace)
5. [Developer Platform](#5-developer-platform)
6. [Chat & Conversations](#6-chat--conversations)
7. [WebSocket Protocol](#7-websocket-protocol)
8. [Wallet & Blockchain](#8-wallet--blockchain)
9. [Mana Points](#9-mana-points)
10. [Referral System](#10-referral-system)
11. [WebRTC Calls](#11-webrtc-calls)
12. [Notifications](#12-notifications)
13. [News Feed](#13-news-feed)
14. [Support & Legal](#14-support--legal)
15. [Internationalization](#15-internationalization)
16. [Crash Reporting](#16-crash-reporting)
17. [Bot API & Webhooks](#17-bot-api--webhooks)
18. [Mini-App JS Bridge](#18-mini-app-js-bridge)
19. [Data Models Reference](#19-data-models-reference)
20. [Error Handling](#20-error-handling)
21. [Currently Mocked — Needs Backend](#21-currently-mocked--needs-backend)

---

## 1. Architecture Overview

### Tech Stack (Client)
| Layer | Technology |
|-------|-----------|
| HTTP | Retrofit 2 + OkHttp 3 |
| Serialization | Gson (lenient) |
| Real-time | OkHttp WebSocket (chat, WebRTC signaling) |
| Auth Storage | EncryptedSharedPreferences (AES256_GCM) |
| Wallet | Local Ed25519 (BIP39 / BIP44 / SLIP-0010) |
| Push | Firebase Cloud Messaging |

### HTTP Configuration
| Parameter | Value |
|-----------|-------|
| Connect Timeout | 30s |
| Read Timeout | 120s |
| Write Timeout | 60s |

### Authentication Model
- **Type:** Bearer JWT in `Authorization` header
- **Refresh:** Automatic on 401 via `POST /api/auth/refresh`
- **Storage:** EncryptedSharedPreferences (`solafon_secure_prefs`, AES256_SIV key / AES256_GCM value)

### Public Endpoints (No Auth Required)
```
GET  /health
POST /api/auth/send-code
POST /api/auth/verify-code
POST /api/auth/refresh
GET  /api/apps              (listing only, excluding /launch)
GET  /api/apps/categories
GET  /api/i18n/languages
GET  /api/support/faq
GET  /api/legal/*
GET  /api/referral/validate/{code}
POST /api/crash/mobile
```

---

## 2. Authentication

### 2.1 Health Check
```
GET /health
```
**Response:**
```json
{
  "status": "ok",
  "version": "1.0.0",
  "timestamp": "2025-01-01T00:00:00Z"
}
```

### 2.2 Send Verification Code
```
POST /api/auth/send-code
```
**Request:**
```json
{ "email": "user@example.com" }
```
**Response:**
```json
{
  "success": true,
  "message": "Code sent to email",
  "expiresIn": 300
}
```

### 2.3 Verify Code & Login
```
POST /api/auth/verify-code
```
**Request:**
```json
{
  "email": "user@example.com",
  "code": "123456",
  "referralCode": "ABCD1234"
}
```
`referralCode` is optional. Validated on login screen in real-time via `GET /api/referral/validate/{code}`.

**Response:**
```json
{
  "success": true,
  "token": "eyJhbGci...",
  "refreshToken": "eyJhbGci...",
  "user": {
    "id": "user_123",
    "email": "user@example.com",
    "displayName": "John",
    "avatarUrl": null
  }
}
```

### 2.4 Refresh Token
```
POST /api/auth/refresh
```
**Request:**
```json
{ "refreshToken": "eyJhbGci..." }
```
**Response:**
```json
{
  "token": "new_access_token",
  "refreshToken": "new_refresh_token"
}
```
Called automatically by client's `TokenAuthenticator` on any 401. If refresh fails → clear auth → redirect to login.

### 2.5 Logout
```
POST /api/logout
Authorization: Bearer {token}
```
**Response:**
```json
{ "success": true }
```
Client clears local auth data regardless of response.

### Token Flow Summary
1. Email → `send-code` → OTP sent
2. OTP → `verify-code` → JWT + refresh token returned
3. Both stored in EncryptedSharedPreferences (`access_token`, `refresh_token`, `user` as JSON, `is_authenticated` boolean)
4. All authenticated requests: `Authorization: Bearer {accessToken}`
5. On 401 → auto-call `refresh` → get new pair → retry original request
6. If refresh fails → clear all → login screen

---

## 3. User Management

### 3.1 Get Profile
```
GET /api/users/me
Authorization: Bearer {token}
```
**Response:**
```json
{
  "id": "user_123",
  "email": "user@example.com",
  "displayName": "John",
  "avatarUrl": "https://...",
  "createdAt": "2025-01-01T00:00:00Z",
  "stats": {
    "appsUsed": 5,
    "transactions": 12,
    "nftsOwned": 0,
    "appsCreated": 2
  },
  "settings": {
    "language": "en",
    "hapticFeedback": true,
    "pushNotifications": true,
    "emailNotifications": true,
    "marketingEmails": false,
    "biometricEnabled": false,
    "twoFactorEnabled": false
  }
}
```

### 3.2 Update Profile
```
PATCH /api/users/me
Authorization: Bearer {token}
```
**Request (all optional):**
```json
{
  "displayName": "New Name",
  "avatarUrl": "https://..."
}
```
**Response:**
```json
{
  "success": true,
  "user": { "id": "...", "email": "...", "displayName": "...", "avatarUrl": "..." }
}
```

### 3.3 Update Settings
```
PATCH /api/users/me/settings
Authorization: Bearer {token}
```
**Request (all optional, send only changed fields):**
```json
{
  "language": "ru",
  "hapticFeedback": true,
  "pushNotifications": false,
  "emailNotifications": true,
  "marketingEmails": false,
  "biometricEnabled": true,
  "twoFactorEnabled": false
}
```
**Response:**
```json
{ "success": true, "settings": { ... } }
```

### 3.4 Delete Account
```
DELETE /api/users/me
Authorization: Bearer {token}
```
**Response:** `{ "success": true }`

### 3.5 Get Sessions
```
GET /api/users/me/sessions
Authorization: Bearer {token}
```
**Response:**
```json
{
  "sessions": [
    {
      "id": "sess_abc",
      "deviceType": "android",
      "deviceName": "Pixel 7",
      "ipAddress": "1.2.3.4",
      "location": "Moscow, RU",
      "lastActive": "2025-01-01T12:00:00Z",
      "isCurrent": true
    }
  ]
}
```

### 3.6 Revoke Session
```
DELETE /api/users/me/sessions/{sessionId}
Authorization: Bearer {token}
```
**Response:** `{ "success": true }`

### 3.7 Revoke All Sessions
```
DELETE /api/users/me/sessions
Authorization: Bearer {token}
```
Client clears local auth data on success.

---

## 4. Apps Marketplace

### 4.1 List Apps
```
GET /api/apps?category={cat}&search={q}&page={p}&limit={l}&sortBy={sort}
```
| Param | Type | Default | Values |
|-------|------|---------|--------|
| category | string? | null | null/"all" = no filter, or category ID |
| search | string? | null | Free text search |
| page | int | 1 | Page number |
| limit | int | 20 | Items per page |
| sortBy | string | "popular" | "popular", "new", "trending" |

**Response:**
```json
{
  "apps": [
    {
      "id": "app_123",
      "name": "My App",
      "description": "Short description",
      "icon": "🤖",
      "iconUrl": "https://...",
      "category": "AI",
      "url": "https://myapp.com",
      "users": "1.2K",
      "usersCount": 1200,
      "isVerified": true,
      "isTrending": false,
      "rating": 4.5,
      "createdAt": "2025-01-01T00:00:00Z",
      "developer": {
        "id": "dev_456",
        "name": "Dev Studio",
        "isVerified": true
      }
    }
  ],
  "pagination": {
    "currentPage": 1,
    "totalPages": 5,
    "totalItems": 100,
    "hasMore": true
  }
}
```
**Auth:** Not required for listing.

### 4.2 Get Categories
```
GET /api/apps/categories
```
**Response:**
```json
{
  "categories": [
    { "id": "ai", "name": "AI", "icon": "🤖", "appsCount": 15 },
    { "id": "games", "name": "Games", "icon": "🎮", "appsCount": 8 },
    { "id": "trading", "name": "Trading", "icon": "📈", "appsCount": 12 },
    { "id": "depin", "name": "DePIN", "icon": "📡", "appsCount": 3 },
    { "id": "defi", "name": "DeFi", "icon": "💰", "appsCount": 7 },
    { "id": "nft", "name": "NFT", "icon": "🖼️", "appsCount": 5 },
    { "id": "stacking", "name": "Stacking", "icon": "📊", "appsCount": 4 },
    { "id": "tools", "name": "Tools", "icon": "🛠️", "appsCount": 10 }
  ]
}
```
**Note:** Currently 8 categories are hardcoded on the client. This endpoint makes them dynamic.

### 4.3 Get App Detail
```
GET /api/apps/{appId}
```
**Response:**
```json
{
  "id": "app_123",
  "name": "My App",
  "description": "Short description",
  "longDescription": "Full markdown description...",
  "icon": "🤖",
  "iconUrl": "https://...",
  "category": "AI",
  "url": "https://myapp.com",
  "users": "1.2K",
  "usersCount": 1200,
  "isVerified": true,
  "isTrending": false,
  "rating": 4.5,
  "screenshots": ["https://...", "https://..."],
  "tags": ["ai", "chatbot"],
  "reviewsCount": 42,
  "createdAt": "2025-01-01T00:00:00Z",
  "updatedAt": "2025-01-15T00:00:00Z",
  "permissions": ["wallet_read", "notifications"],
  "version": "1.0.0",
  "developer": { "id": "dev_456", "name": "Dev Studio", "isVerified": true }
}
```

### 4.4 Launch App
```
POST /api/apps/{appId}/launch
Authorization: Bearer {token}
```
**Response:**
```json
{ "success": true, "sessionId": "session_abc" }
```
Records app usage for analytics. Client opens app URL in WebView.

---

## 5. Developer Platform

### 5.1 List My Apps
```
GET /api/developer/apps
Authorization: Bearer {token}
```
**Response:**
```json
{
  "apps": [
    {
      "id": "app_123",
      "name": "My Bot",
      "description": "Bot description",
      "icon": "🤖",
      "category": "AI",
      "url": "https://mybot.com",
      "status": "approved",
      "isVerified": false,
      "users": "100",
      "usersCount": 100,
      "rating": 4.2,
      "createdAt": "2025-01-01T00:00:00Z",
      "updatedAt": "2025-01-15T00:00:00Z",
      "moderationNote": null
    }
  ],
  "stats": {
    "totalApps": 3,
    "approvedApps": 2,
    "pendingApps": 1,
    "rejectedApps": 0,
    "totalUsers": 350
  }
}
```

### 5.2 Create App
```
POST /api/developer/apps
Authorization: Bearer {token}
```
**Request:**
```json
{
  "name": "My Bot",
  "description": "A helpful bot",
  "icon": "🤖",
  "category": "AI",
  "url": "https://mybot.com"
}
```
**Response:**
```json
{
  "success": true,
  "app": { ... },
  "message": "App created successfully",
  "apiKey": "sk_live_abc123...",
  "keyId": "key_xyz",
  "apiKeyHint": "sk_live_abc..."
}
```
**IMPORTANT:** `apiKey` is returned **ONLY ONCE** at creation. Client shows it in a modal for the user to copy. Cannot be retrieved later.

App statuses: `"pending"` → `"approved"` or `"rejected"`

### 5.3 Get App Detail (Developer)
```
GET /api/developer/apps/{appId}
Authorization: Bearer {token}
```
Returns full `DeveloperApp` with `moderationNote`.

### 5.4 Update App
```
PATCH /api/developer/apps/{appId}
Authorization: Bearer {token}
```
**Request (all optional):**
```json
{
  "name": "Updated Name",
  "description": "Updated description",
  "icon": "🎮",
  "category": "Games",
  "url": "https://new-url.com"
}
```

### 5.5 Delete App
```
DELETE /api/developer/apps/{appId}
Authorization: Bearer {token}
```

### 5.6 Generate API Key
```
POST /api/developer/apps/{appId}/api-key
Authorization: Bearer {token}
```
**Request:**
```json
{
  "appId": "app_123",
  "permissions": ["messages.send", "messages.read"]
}
```
**Response:**
```json
{
  "success": true,
  "apiKey": "sk_live_newkey...",
  "apiSecret": "secret_...",
  "keyId": "key_new",
  "message": "Key generated"
}
```

### 5.7 List API Credentials
```
GET /api/developer/apps/{appId}/api-key
Authorization: Bearer {token}
```
**Response:**
```json
{
  "credentials": [
    {
      "id": "key_xyz",
      "appId": "app_123",
      "apiKeyPrefix": "sk_live_abc...",
      "webhookUrl": "https://mybot.com/webhook",
      "webhookSecret": "whsec_...",
      "isActive": true,
      "createdAt": "2025-01-01T00:00:00Z",
      "lastUsedAt": "2025-01-15T12:00:00Z",
      "permissions": ["messages.send", "messages.read"]
    }
  ]
}
```

### 5.8 Revoke API Key
```
DELETE /api/developer/apps/{appId}/credentials/{keyId}
Authorization: Bearer {token}
```
**Legacy (revoke all):** `DELETE /api/developer/apps/{appId}/api-key`

### 5.9 Update Webhook
```
PUT /api/developer/apps/{appId}/webhook
Authorization: Bearer {token}
```
**Request:**
```json
{
  "webhookUrl": "https://mybot.com/webhook",
  "events": ["message", "callback", "conversation_start"]
}
```
**Response:**
```json
{
  "success": true,
  "webhookSecret": "whsec_abc123..."
}
```

### 5.10 Get Welcome Message
```
GET /api/developer/apps/{appId}/welcome-message
Authorization: Bearer {token}
```
**Response:**
```json
{
  "welcomeMessage": {
    "id": "wm_123",
    "appId": "app_123",
    "content": { "type": "text", "text": "Hello! How can I help?" },
    "isActive": true,
    "createdAt": "2025-01-01T00:00:00Z"
  }
}
```

### 5.11 Update Welcome Message
```
PUT /api/developer/apps/{appId}/welcome-message
Authorization: Bearer {token}
```
**Request:**
```json
{
  "content": { "type": "text", "text": "Welcome! Ask me anything." },
  "isActive": true
}
```

### 5.12 File Upload (NEEDED)
```
POST /api/developer/upload
Authorization: Bearer {token}
Content-Type: multipart/form-data
```
**Request:**
- `file`: binary
- `type`: `"icon"` | `"banner"` (optional)

**Response:**
```json
{
  "url": "https://cdn.solafon.com/uploads/abc123.png",
  "filename": "abc123.png",
  "size": 245000,
  "file_type": "image/png"
}
```
**Constraints:**
- Max size: 5MB
- Formats: PNG, JPG, WEBP, GIF
- `type=icon`: auto-resize to 256x256
- `type=banner`: validate min 640x360

**Note:** Currently only `POST /api/admin/upload` exists (admin only). Need developer-accessible version.

### Needed Model Changes
- Add `iconUrl: String?` to `CreateAppRequest` and `UpdateAppRequest`
- Add `welcomeBannerUrl: String?` to `CreateAppRequest`, `UpdateAppRequest`, and `App` model
- Add `welcome_banner_url` column to `apps` table

---

## 6. Chat & Conversations

### 6.1 List Conversations
```
GET /api/conversations?page={p}&limit={l}
Authorization: Bearer {token}
```
**Response:**
```json
{
  "conversations": [
    {
      "id": "conv_123",
      "appId": "app_456",
      "userId": "user_789",
      "appName": "My Bot",
      "appIcon": "🤖",
      "appIconUrl": "https://...",
      "appUrl": "https://mybot.com",
      "lastMessage": {
        "id": "msg_abc",
        "content": { "type": "text", "text": "Hello!" },
        "timestamp": 1704067200000,
        "senderType": "bot"
      },
      "unreadCount": 2,
      "createdAt": "2025-01-01T00:00:00Z",
      "updatedAt": "2025-01-01T12:00:00Z",
      "isActive": true
    }
  ],
  "pagination": { ... }
}
```

### 6.2 Start Conversation
```
POST /api/apps/{appId}/conversations
Authorization: Bearer {token}
```
**Request:**
```json
{
  "appId": "app_456",
  "initialMessage": "Hello!"
}
```
**Response:**
```json
{
  "success": true,
  "conversation": { ... },
  "welcomeMessage": {
    "id": "msg_welcome",
    "content": { "type": "text", "text": "Welcome!" },
    "senderType": "bot",
    "timestamp": 1704067200000
  }
}
```
WebSocket auto-subscribes after creation.

### 6.3 Get Messages
```
GET /api/conversations/{conversationId}/messages?limit={l}&before={cursor}
Authorization: Bearer {token}
```
| Param | Type | Default | Description |
|-------|------|---------|-------------|
| limit | int | 50 | Messages per page |
| before | string? | null | Message ID cursor for pagination |

**Response:**
```json
{
  "messages": [
    {
      "id": "msg_123",
      "appId": "app_456",
      "conversationId": "conv_789",
      "senderId": "user_123",
      "senderType": "user",
      "content": {
        "type": "text",
        "text": "Hello!"
      },
      "timestamp": 1704067200000,
      "status": "read",
      "replyToId": null,
      "metadata": null
    },
    {
      "id": "msg_124",
      "senderType": "bot",
      "content": {
        "type": "button",
        "text": "Choose an option:",
        "buttons": [
          { "id": "btn_1", "text": "Option A", "action": "callback", "payload": "opt_a" },
          { "id": "btn_2", "text": "Visit Site", "action": "url", "url": "https://..." }
        ]
      }
    }
  ],
  "pagination": { ... },
  "hasMore": true
}
```

### 6.4 Send Message
```
POST /api/conversations/{conversationId}/messages
Authorization: Bearer {token}
```
**Request:**
```json
{
  "content": { "type": "text", "text": "Hello bot!" },
  "replyToId": null,
  "metadata": null
}
```
**Response:**
```json
{ "success": true, "message": { ... } }
```
Triggers `message.received` webhook on the bot's server.

### 6.5 Button Callback
```
POST /api/conversations/{conversationId}/callback
Authorization: Bearer {token}
```
**Request:**
```json
{
  "messageId": "msg_124",
  "buttonId": "btn_1",
  "payload": "opt_a"
}
```
**Response:**
```json
{
  "success": true,
  "message": { ... },
  "action": "openUrl"
}
```
`action`: `"openUrl"` | `"openWebApp"` | `"showAlert"` | `null`

Triggers `callback.received` webhook on the bot's server.

### 6.6 Mark as Read
```
POST /api/conversations/{conversationId}/read
Authorization: Bearer {token}
```
Resets `unreadCount` to 0.

### 6.7 Delete Conversation
```
DELETE /api/conversations/{conversationId}
Authorization: Bearer {token}
```
Triggers `conversation.ended` webhook on the bot's server.

---

## 7. WebSocket Protocol

### 7.1 Chat WebSocket

**URL:** `wss://api.solafon.com/ws?token={accessToken}`

| Config | Value |
|--------|-------|
| Ping interval | 30s |
| Max reconnect attempts | 5 |
| Reconnect base delay | 5s (exponential backoff) |
| Fallback | HTTP polling (3s active / 10s idle) |

**Connection States:** `DISCONNECTED` → `CONNECTING` → `CONNECTED` | `ERROR`

#### Server → Client Events

| Event | Payload | Description |
|-------|---------|-------------|
| `new_message` | `ChatMessage` | New message from bot |
| `typing` | `{ conversationId, isTyping, timestamp }` | Bot typing indicator |
| `message_status` | `{ messageId, status }` | Delivery status update |
| `conversation_update` | `{ conversationId, unreadCount }` | Conversation metadata change |
| `pong` | `{}` | Keepalive response |
| `error` | `{ message }` | Error notification |

#### Client → Server Events

| Event | Payload | Description |
|-------|---------|-------------|
| `ping` | `{}` | Keepalive |
| `subscribe` | `{ type: "subscribe", conversationId }` | Subscribe to updates |
| `unsubscribe` | `{ type: "unsubscribe", conversationId }` | Unsubscribe |

#### Polling Fallback
When WebSocket is unavailable:
- **Active** (receiving messages): poll every **3 seconds**
- **Idle** (no new messages for 5+ polls): every **10 seconds**
- Endpoint: `GET /api/conversations/{id}/messages?limit=20`

### 7.2 WebRTC Signaling WebSocket

**URL:** `wss://api.solafon.com/ws/webrtc?roomId={roomId}`
**Auth:** `Authorization: Bearer {token}` header

#### Server → Client

| Event | Payload |
|-------|---------|
| `offer` | `{ type: "offer", data: { type, sdp } }` |
| `answer` | `{ type: "answer", data: { type, sdp } }` |
| `ice-candidate` | `{ type: "ice-candidate", data: { sdpMid, sdpMLineIndex, candidate } }` |
| `user-joined` | `{ type: "user-joined", data: { userId } }` |
| `user-left` | `{ type: "user-left", data: { userId } }` |
| `call-ended` | `{ type: "call-ended" }` |

#### Client → Server

| Event | Payload |
|-------|---------|
| `join` | `{ type: "join", roomId }` |
| `offer` | `{ type: "offer", roomId, data: { type, sdp } }` |
| `answer` | `{ type: "answer", roomId, data: { type, sdp } }` |
| `ice-candidate` | `{ type: "ice-candidate", roomId, data: { sdpMid, sdpMLineIndex, candidate } }` |

---

## 8. Wallet & Blockchain

### Architecture
| Property | Value |
|----------|-------|
| Type | Non-custodial (keys never leave device) |
| Chain | Solana (mainnet-beta) |
| Key type | Ed25519 via BIP39/BIP44/SLIP-0010 |
| Derivation path | `m/44'/501'/0'/0'` |
| Client storage | EncryptedSharedPreferences (`solafon_wallet_secure`) |
| RPC | `https://api.mainnet-beta.solana.com` |
| System Program | `11111111111111111111111111111111` |
| Lamports | 1 SOL = 1,000,000,000 lamports |

### Key Management (Client-Side Only)
```
Create wallet:
  1. Generate 12-word BIP39 mnemonic (128-bit entropy)
  2. Derive seed via PBKDF2
  3. Apply SLIP-0010 Ed25519 derivation (m/44'/501'/0'/0')
  4. Extract 32-byte private key → derive Ed25519 public key
  5. Encode as Base58 → wallet address
  6. Store mnemonic + address in EncryptedSharedPreferences

Import wallet:
  - Validate mnemonic (12 or 24 words, BIP39 wordlist)
  - Same derivation flow

Transaction signing:
  1. Build Solana transaction message
  2. Sign locally with Ed25519 private key (64-byte signature)
  3. Format: [num_signatures(1)][signature(64)][message]
  4. Send base64-encoded signed tx to backend for broadcast
```

### 8.1 Get Balance
```
GET /api/wallet/balance?address={solanaAddress}
Authorization: Bearer {token}
```
**Response:**
```json
{
  "address": "DRpbCBMxVnDK...",
  "sol": "1.5",
  "tokens": [
    {
      "mint": "EPjFWdd5AufqSSqeM2qN1xzybapC8G4wEGGkZwyTDt1v",
      "symbol": "USDC",
      "name": "USD Coin",
      "balance": "100.00",
      "decimals": 6,
      "icon": "https://...",
      "priceUsd": 1.0,
      "valueUsd": 100.0
    }
  ],
  "totalUsd": "325.50"
}
```

### 8.2 Send Transaction
```
POST /api/wallet/send
Authorization: Bearer {token}
```
**Request:**
```json
{ "signedTransaction": "base64_encoded_signed_tx..." }
```
**Response:**
```json
{ "signature": "5KtP3...", "status": "confirmed" }
```
**IMPORTANT:** Transaction is signed locally on device. Server only broadcasts the pre-signed transaction to Solana network. Private keys NEVER sent to server.

### 8.3 Transaction History
```
GET /api/wallet/transactions?address={addr}&limit={l}&before={cursor}
Authorization: Bearer {token}
```
**Response:**
```json
{
  "transactions": [
    {
      "signature": "5KtP3...",
      "type": "transfer",
      "direction": "out",
      "token": "SOL",
      "tokenMint": null,
      "amount": "0.5",
      "from": "DRpbCBMx...",
      "to": "7nYBs...",
      "fee": "0.000005",
      "timestamp": 1704067200000,
      "status": "finalized",
      "slot": 123456789
    }
  ],
  "hasMore": true,
  "nextCursor": "before_sig_abc..."
}
```

### 8.4 Supported Tokens
```
GET /api/wallet/tokens
Authorization: Bearer {token}
```
**Response:**
```json
{
  "tokens": [
    { "mint": "So11111...", "symbol": "SOL", "name": "Solana", "decimals": 9, "icon": "https://..." },
    { "mint": "EPjFWdd5...", "symbol": "USDC", "name": "USD Coin", "decimals": 6, "icon": "https://..." }
  ]
}
```

### 8.5 Token Prices
```
GET /api/wallet/prices?mints={comma_separated_or_all}
Authorization: Bearer {token}
```
**Response:**
```json
{
  "prices": {
    "So11111...": 148.35,
    "EPjFWdd5...": 1.0
  }
}
```

### 8.6 Transaction Status
```
GET /api/wallet/status?signature={txSignature}
Authorization: Bearer {token}
```
**Response:**
```json
{ "signature": "5KtP3...", "status": "finalized" }
```
Statuses: `"pending"` → `"confirmed"` → `"finalized"` | `"failed"`

---

## 9. Mana Points

> Currently fully mocked on client. All endpoints below need backend implementation.

### 9.1 Get Mana Points Balance (NEEDED)
```
GET /api/wallet/mana-points
Authorization: Bearer {token}
```
**Response:**
```json
{
  "success": true,
  "balance": 500,
  "balanceUsdt": 11.50,
  "income": 500,
  "rewards": 12,
  "history": [
    {
      "id": "mp_tx_123",
      "type": "earned",
      "amount": 50,
      "source": "daily_reward",
      "description": "Daily login reward",
      "createdAt": "2026-03-01T10:00:00Z"
    }
  ]
}
```

### 9.2 Get Tariffs (NEEDED)
```
GET /api/wallet/mana-points/tariffs
Authorization: Bearer {token}
```
**Response:**
```json
{
  "success": true,
  "tariffs": [
    { "id": "mp_250", "mp": 250, "priceUsdt": 5.79, "currency": "USD" },
    { "id": "mp_500", "mp": 500, "priceUsdt": 11.90, "isPopular": true, "currency": "USD" },
    { "id": "mp_1000", "mp": 1000, "priceUsdt": 22.99, "currency": "USD" },
    { "id": "mp_2500", "mp": 2500, "priceUsdt": 57.98, "currency": "USD" }
  ]
}
```

### 9.3 Purchase Mana Points (NEEDED)
```
POST /api/wallet/mana-points/purchase
Authorization: Bearer {token}
```
**Request:**
```json
{ "tariffId": "mp_500", "paymentMethod": "solana" }
```
**Response:**
```json
{
  "success": true,
  "transactionId": "tx_abc123",
  "newBalance": 1000,
  "amountPurchased": 500,
  "amountCharged": 11.90
}
```

### 9.4 Gift Mana Points (NEEDED)
```
POST /api/wallet/mana-points/gift
Authorization: Bearer {token}
```
**Request:**
```json
{
  "recipientAddress": "DRpbCBMxVnDK...",
  "amount": 100
}
```
**Response:**
```json
{
  "success": true,
  "transactionId": "tx_gift_123",
  "newBalance": 400,
  "amountSent": 100
}
```

### 9.5 Get Wallet Networks (NEEDED)
```
GET /api/wallet/networks
Authorization: Bearer {token}
```
**Response:**
```json
{
  "success": true,
  "networks": [
    {
      "id": "solana",
      "name": "Solana",
      "address": "DRpbCBMxVnDK...",
      "iconColor": "#9945FF",
      "isDefault": true
    },
    {
      "id": "bitcoin",
      "name": "Bitcoin",
      "address": "bc1qw508d6...",
      "iconColor": "#F7931A",
      "isDefault": false
    }
  ]
}
```

### Database Schema (Suggested)
```sql
CREATE TABLE mana_point_transactions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id),
    type VARCHAR(20) NOT NULL, -- 'purchase', 'earned', 'gift_sent', 'gift_received', 'spent'
    amount INT NOT NULL,
    source VARCHAR(50),
    description TEXT,
    related_tx_id UUID,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE mana_point_tariffs (
    id VARCHAR(20) PRIMARY KEY,
    mp_amount INT NOT NULL,
    price_usdt DECIMAL(10,2) NOT NULL,
    currency VARCHAR(3) DEFAULT 'USD',
    is_popular BOOLEAN DEFAULT false,
    is_active BOOLEAN DEFAULT true
);

CREATE INDEX idx_mana_tx_user ON mana_point_transactions(user_id, created_at DESC);
```

---

## 10. Referral System

### 10.1 Get Referral Info
```
GET /api/referral
Authorization: Bearer {token}
```
**Response:**
```json
{
  "data": {
    "code": "ABCD1234",
    "link": "https://solafon.com/ref/ABCD1234",
    "totalCount": 5,
    "referrals": [
      {
        "id": "user_abc",
        "displayName": "Alice",
        "avatarUrl": "https://...",
        "registeredAt": "2025-01-15T00:00:00Z"
      }
    ]
  }
}
```

### 10.2 Validate Referral Code
```
GET /api/referral/validate/{code}
```
**Auth:** Not required (used during registration)

**Response:**
```json
{
  "data": {
    "valid": true,
    "referrerName": "John"
  }
}
```

---

## 11. WebRTC Calls

### 11.1 Create Call Room
```
POST /api/calls/rooms/create
Authorization: Bearer {token}
```
**Request:**
```json
{
  "type": "direct",
  "participantIds": ["user_abc"]
}
```
**Response:**
```json
{
  "data": {
    "room": {
      "id": "room_123",
      "roomCode": "ABC-DEF",
      "type": "direct",
      "status": "waiting",
      "createdBy": "user_me",
      "startedAt": null,
      "endedAt": null,
      "duration": 0,
      "createdAt": "2025-01-01T00:00:00Z",
      "participants": [
        {
          "id": "part_1",
          "roomId": "room_123",
          "userId": "user_me",
          "status": "joined",
          "joinedAt": "2025-01-01T00:00:00Z",
          "isMuted": false,
          "isVideoOn": true,
          "isAudioOn": true,
          "user": { ... }
        }
      ]
    },
    "deepLink": "solafon://call/ABC-DEF",
    "roomCode": "ABC-DEF"
  }
}
```

### 11.2 Get Room by Code
```
GET /api/calls/rooms/code/{roomCode}
Authorization: Bearer {token}
```

### 11.3 Join Room
```
POST /api/calls/rooms/{roomId}/join
Authorization: Bearer {token}
```

### 11.4 End Call
```
POST /api/calls/rooms/{roomId}/end
Authorization: Bearer {token}
```

### 11.5 Update Participant Status
```
PATCH /api/calls/rooms/{roomId}/status
Authorization: Bearer {token}
```
**Request:**
```json
{ "isMuted": true, "isVideoOn": false, "isAudioOn": true }
```

---

## 12. Notifications

### Existing Endpoints

#### 12.1 Register Push Token
```
POST /api/notifications/register
Authorization: Bearer {token}
```
**Request:**
```json
{
  "fcmToken": "firebase_token...",
  "deviceId": "unique_device_id",
  "platform": "android"
}
```

#### 12.2 Unregister Push Token
```
DELETE /api/notifications/unregister
Authorization: Bearer {token}
```

### Needed Endpoints

#### 12.3 List Notifications (NEEDED)
```
GET /api/notifications?page={p}&limit={l}&unreadOnly={bool}
Authorization: Bearer {token}
```
**Response:**
```json
{
  "success": true,
  "notifications": [
    {
      "id": "notif_abc123",
      "title": "Your app was approved",
      "body": "Your app 'MyBot' has been approved and is now live.",
      "type": "app_moderation",
      "isRead": false,
      "createdAt": "2026-02-27T10:30:00Z",
      "actionUrl": "/apps/app_123"
    }
  ],
  "pagination": { "page": 1, "limit": 20, "total": 45, "hasMore": true }
}
```

**Notification types:**
- `app_moderation` — app approved/rejected
- `security` — new device login, password change
- `system` — platform updates, maintenance
- `transaction` — deposit, withdrawal, transfer
- `promotion` — special offers, new features

#### 12.4 Mark as Read (NEEDED)
```
POST /api/notifications/{id}/read
Authorization: Bearer {token}
```

#### 12.5 Mark All as Read (NEEDED)
```
POST /api/notifications/read-all
Authorization: Bearer {token}
```
**Response:**
```json
{ "success": true, "updatedCount": 32 }
```

#### 12.6 Get Unread Count (NEEDED)
```
GET /api/notifications/count
Authorization: Bearer {token}
```
**Response:**
```json
{ "success": true, "unreadCount": 32 }
```
Used for badge on bell icon in tab bar.

### Database Schema (Suggested)
```sql
CREATE TABLE notifications (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id),
    title VARCHAR(255) NOT NULL,
    body TEXT NOT NULL,
    type VARCHAR(50) NOT NULL,
    is_read BOOLEAN DEFAULT false,
    action_url VARCHAR(500),
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_notifications_user_unread
    ON notifications(user_id, is_read)
    WHERE is_read = false;
```

---

## 13. News Feed

> Currently fully mocked on client. All endpoints below need backend implementation.

### 13.1 Get News Feed (NEEDED)
```
GET /api/news?page={p}&limit={l}
Authorization: Bearer {token}
```
**Response:**
```json
{
  "success": true,
  "posts": [
    {
      "id": "post_abc123",
      "appId": "app_456",
      "appName": "Moltbook",
      "appIcon": "https://cdn.solafon.com/icons/moltbook.png",
      "text": "We just launched our new feature...",
      "imageUrl": "https://cdn.solafon.com/posts/img123.jpg",
      "commentsCount": 3307,
      "likesCount": 166,
      "sharesCount": 46,
      "isLiked": false,
      "createdAt": "2026-03-01T03:30:00Z"
    }
  ],
  "pagination": { ... }
}
```

### 13.2 Like/Unlike Post (NEEDED)
```
POST /api/news/{postId}/like
Authorization: Bearer {token}
```
**Response:**
```json
{ "success": true, "liked": true, "likesCount": 167 }
```

### 13.3 Share Post (NEEDED)
```
POST /api/news/{postId}/share
Authorization: Bearer {token}
```
**Response:**
```json
{ "success": true, "sharesCount": 47 }
```

### 13.4 Get Comments (NEEDED)
```
GET /api/news/{postId}/comments?page={p}&limit={l}
Authorization: Bearer {token}
```
**Response:**
```json
{
  "success": true,
  "comments": [
    {
      "id": "comment_789",
      "userId": "user_123",
      "userName": "John",
      "userAvatar": "https://...",
      "text": "Great update!",
      "createdAt": "2026-03-01T04:00:00Z"
    }
  ]
}
```

### 13.5 Post Comment (NEEDED)
```
POST /api/news/{postId}/comments
Authorization: Bearer {token}
```
**Request:**
```json
{ "text": "Nice!" }
```

### 13.6 Developer: Create News Post (NEEDED)
```
POST /api/developer/apps/{appId}/news
Authorization: Bearer {token}
```
**Request:**
```json
{
  "text": "We just launched...",
  "imageUrl": "https://cdn.solafon.com/posts/img.jpg"
}
```

### Database Schema (Suggested)
```sql
CREATE TABLE news_posts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    app_id UUID NOT NULL REFERENCES apps(id),
    text TEXT NOT NULL,
    image_url VARCHAR(500),
    comments_count INT DEFAULT 0,
    likes_count INT DEFAULT 0,
    shares_count INT DEFAULT 0,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE news_likes (
    post_id UUID REFERENCES news_posts(id),
    user_id UUID REFERENCES users(id),
    created_at TIMESTAMPTZ DEFAULT NOW(),
    PRIMARY KEY (post_id, user_id)
);

CREATE TABLE news_comments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    post_id UUID NOT NULL REFERENCES news_posts(id),
    user_id UUID NOT NULL REFERENCES users(id),
    text TEXT NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_news_posts_created ON news_posts(created_at DESC);
CREATE INDEX idx_news_comments_post ON news_comments(post_id, created_at);
```

---

## 14. Support & Legal

### 14.1 Get FAQ
```
GET /api/support/faq?language={lang}
```
**Response:**
```json
{
  "faq": [
    { "id": "faq_1", "question": "How to create a wallet?", "answer": "...", "category": "wallet" }
  ]
}
```

### 14.2 Create Support Ticket
```
POST /api/support/tickets
Authorization: Bearer {token}
```
**Request:**
```json
{ "subject": "Issue with wallet", "message": "Detailed description...", "category": "wallet" }
```
**Response:**
```json
{ "ticketId": "ticket_123", "status": "open" }
```

### 14.3 Terms of Service
```
GET /api/legal/terms
```
**Response:**
```json
{ "content": "Markdown content...", "version": "1.0", "lastUpdated": "2025-01-01" }
```

### 14.4 Privacy Policy
```
GET /api/legal/privacy
```
Same format as Terms.

---

## 15. Internationalization

### 15.1 Get Languages
```
GET /api/i18n/languages
```
**Response:**
```json
{
  "languages": [
    { "code": "en", "name": "English", "flag": "🇺🇸" },
    { "code": "ru", "name": "Русский", "flag": "🇷🇺" },
    { "code": "zh", "name": "中文", "flag": "🇨🇳" },
    { "code": "es", "name": "Español", "flag": "🇪🇸" }
  ],
  "defaultLanguage": "en"
}
```

---

## 16. Crash Reporting

### 16.1 Report Crash
```
POST /api/crash/mobile
```
**Auth:** Not required (device-level)

**Request:**
```json
{
  "userId": "user_123",
  "sessionId": "sess_abc",
  "errorType": "crash",
  "message": "NullPointerException in WalletScreen",
  "stack": "java.lang.NullPointerException...",
  "userAgent": "Android/14 Solafon/1.1.62",
  "source": "android",
  "metadata": { "screen": "WalletScreen", "buildVersion": "1.1.62" }
}
```
**Error types:** `"crash"`, `"error"`, `"warning"`

**Response:**
```json
{
  "data": {
    "id": "crash_123",
    "appId": "solafon-android",
    "errorType": "crash",
    "message": "NullPointerException...",
    "createdAt": "2025-01-01T00:00:00Z"
  }
}
```

---

## 17. Bot API & Webhooks

### Webhook Delivery

When a user interacts with a bot, the server sends webhook requests to the developer's configured `webhookUrl`.

**Headers:**
```
Content-Type: application/json
X-Webhook-Signature: HMAC-SHA256(body, webhookSecret)
```

### Webhook Events

#### message.received
User sends a message to the bot.
```json
{
  "event": "message.received",
  "timestamp": 1704067200000,
  "data": {
    "conversationId": "conv_123",
    "message": {
      "id": "msg_456",
      "senderId": "user_789",
      "senderType": "user",
      "content": { "type": "text", "text": "Hello!" },
      "timestamp": 1704067200000
    }
  }
}
```

#### callback.received
User clicks a callback button.
```json
{
  "event": "callback.received",
  "timestamp": 1704067200000,
  "data": {
    "conversationId": "conv_123",
    "messageId": "msg_124",
    "buttonId": "btn_1",
    "payload": "opt_a",
    "userId": "user_789"
  }
}
```

#### conversation.started
User opens a conversation with the bot.
```json
{
  "event": "conversation.started",
  "timestamp": 1704067200000,
  "data": {
    "conversationId": "conv_123",
    "userId": "user_789",
    "initialMessage": "Hello!"
  }
}
```

#### conversation.ended
User deletes the conversation.
```json
{
  "event": "conversation.ended",
  "timestamp": 1704067200000,
  "data": {
    "conversationId": "conv_123",
    "userId": "user_789"
  }
}
```

### Bot Authentication
Bots authenticate using the API key:
- **Header:** `X-Bot-Token: {apiKey}`
- **Or query:** `?token={apiKey}`

### Bot Message Capabilities
Bots can send rich messages via the same message endpoint:
```
POST /api/conversations/{conversationId}/messages
X-Bot-Token: {apiKey}
```

**Rich content types:**
```json
{
  "content": {
    "type": "card",
    "text": "Check out these options:",
    "cards": [
      {
        "id": "card_1",
        "title": "Premium Plan",
        "subtitle": "$9.99/month",
        "imageUrl": "https://...",
        "buttons": [
          { "id": "buy_1", "text": "Buy Now", "action": "callback", "payload": "buy_premium" },
          { "id": "info_1", "text": "Learn More", "action": "url", "url": "https://..." },
          { "id": "open_1", "text": "Open App", "action": "webApp", "url": "https://..." }
        ]
      }
    ]
  }
}
```

**Button action types:**
| Action | Description |
|--------|-------------|
| `callback` | Triggers `callback.received` webhook |
| `url` | Opens external link in browser |
| `webApp` | Opens mini-app URL in WebView |

---

## 18. Mini-App JS Bridge

When a mini-app runs in the Solafon WebView, it has access to these JavaScript APIs:

### `window.solana` — Wallet Adapter
```javascript
// Connect wallet (returns publicKey)
const resp = await window.solana.connect();
console.log(resp.publicKey); // Base58 wallet address

// Sign and send transaction
const signature = await window.solana.signAndSendTransaction(transaction);

// Sign transaction without sending
const signedTx = await window.solana.signTransaction(transaction);

// Sign multiple transactions
const signedTxs = await window.solana.signAllTransactions(transactions);

// Sign arbitrary message
const signature = await window.solana.signMessage(messageBytes);

// Disconnect
await window.solana.disconnect();

// Events
window.solana.on('connect', (publicKey) => { ... });
window.solana.on('disconnect', () => { ... });
window.solana.on('accountChanged', (publicKey) => { ... });
```

### `window.Solafon` — Native Bridge
```javascript
// User info
const user = Solafon.getUser();
// → { id, email, displayName, avatarUrl }

// Platform info
const platform = Solafon.platform;
// → { os: "android", version: "1.1.62" }

// UI
Solafon.showAlert("Hello!");
Solafon.showConfirm("Are you sure?");
Solafon.showToast("Copied!", "short");  // "short" or "long"
Solafon.haptic("light");                // "light", "medium", "heavy"

// Navigation
Solafon.openLink("https://...");
Solafon.close();                        // close WebView

// Theme
Solafon.getTheme();                     // returns "dark"
```

### WebSocket Bridge
For mini-apps needing WebSocket connections to localhost (`ws://`), the native bridge proxies through OkHttp to solve mixed-content security issues.
- Intercepts: `ws://localhost:*`, `ws://127.0.0.1:*`, `ws://10.0.2.2:*`

### Mobile Wallet Adapter (MWA)
Allows other Solana dApps to connect to the Solafon wallet via MWA protocol:
- `authorize()` — connect wallet
- `signTransaction()` — sign transaction
- `signAndSendTransaction()` — sign and broadcast
- `signMessage()` — sign arbitrary message
- `reauthorize()` — reconnect without user approval

---

## 19. Data Models Reference

### Common Wrappers
```typescript
interface ApiError {
  code: string;
  message: string;
  details?: Record<string, string>;
}

interface ApiResponse<T> {
  data?: T;
  error?: ApiError;
}

interface SuccessResponse {
  success: boolean;
  message?: string;
}

interface Pagination {
  currentPage: number;
  totalPages: number;
  totalItems: number;
  hasMore: boolean;
}
```

### User
```typescript
interface User {
  id: string;
  email: string;
  displayName: string;
  avatarUrl: string | null;
}

interface UserStats {
  appsUsed: number;
  transactions: number;
  nftsOwned: number;
  appsCreated: number;
}

interface UserSettings {
  language: string;
  hapticFeedback: boolean;
  pushNotifications: boolean;
  emailNotifications: boolean;
  marketingEmails: boolean;
  biometricEnabled: boolean;
  twoFactorEnabled: boolean;
}
```

### App
```typescript
interface App {
  id: string;
  name: string;
  description: string;
  icon: string;           // emoji
  iconUrl: string | null; // image URL
  category: string;
  url: string;
  users: string;          // formatted: "1.2K"
  usersCount: number;
  isVerified: boolean;
  isTrending: boolean;
  rating: number | null;
  createdAt: string;
  developer: AppDeveloper | null;
}

interface AppDeveloper {
  id: string;
  name: string;
  isVerified: boolean;
}

interface DeveloperApp extends App {
  status: "pending" | "approved" | "rejected";
  updatedAt: string | null;
  moderationNote: string | null;
}
```

### Chat
```typescript
interface ChatConversation {
  id: string;
  appId: string;
  userId: string | null;
  appName: string | null;
  appIcon: string | null;     // emoji
  appIconUrl: string | null;  // image URL
  appUrl: string | null;      // WebView URL
  lastMessage: ChatMessage | null;
  unreadCount: number;
  createdAt: string | null;
  updatedAt: string | null;
  isActive: boolean;
}

interface ChatMessage {
  id: string;
  appId: string;
  conversationId: string;
  senderId: string;
  senderType: "user" | "bot" | "system";
  content: MessageContent;
  timestamp: number;          // milliseconds
  status: "sending" | "sent" | "delivered" | "read" | "failed";
  replyToId: string | null;
  metadata: Record<string, string> | null;
}

interface MessageContent {
  type: "text" | "image" | "button" | "card" | "carousel";
  text: string | null;
  imageUrl: string | null;
  buttons: MessageButton[] | null;
  cards: MessageCard[] | null;
}

interface MessageButton {
  id: string;
  text: string;
  action: "callback" | "url" | "webApp";
  payload: string | null;
  url: string | null;
}

interface MessageCard {
  id: string;
  title: string;
  subtitle: string | null;
  imageUrl: string | null;
  buttons: MessageButton[] | null;
}
```

### Wallet
```typescript
interface WalletBalance {
  address: string;
  sol: string;
  tokens: TokenBalance[];
  totalUsd: string;
}

interface TokenBalance {
  mint: string;
  symbol: string;
  name: string;
  balance: string;
  decimals: number;
  icon: string | null;
  priceUsd: number | null;
  valueUsd: number | null;
}

interface WalletTransaction {
  signature: string;
  type: "transfer" | "swap" | "unknown";
  direction: "in" | "out";
  token: string;
  tokenMint: string | null;
  amount: string;
  from: string;
  to: string;
  fee: string | null;
  timestamp: number;
  status: "confirmed" | "finalized" | "failed";
  slot: number | null;
}
```

### Calls
```typescript
interface CallRoom {
  id: string;
  roomCode: string;
  type: "direct" | "conference";
  status: string;
  createdBy: string;
  startedAt: string | null;
  endedAt: string | null;
  duration: number;
  createdAt: string;
  participants: CallParticipant[] | null;
}

interface CallParticipant {
  id: string;
  roomId: string;
  userId: string;
  status: string;
  joinedAt: string | null;
  isMuted: boolean;
  isVideoOn: boolean;
  isAudioOn: boolean;
  user: User | null;
}
```

---

## 20. Error Handling

### Standard Error Response
```json
{
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "Email is required",
    "details": { "email": "must not be empty" }
  }
}
```

### HTTP Status Codes
| Code | Meaning | Client Behavior |
|------|---------|-----------------|
| 200 | Success | Process response |
| 400 | Bad Request | Show error message |
| 401 | Unauthorized | Auto-refresh token → retry |
| 403 | Forbidden | Show "access denied" |
| 404 | Not Found | Show not found UI |
| 409 | Conflict | Show conflict message |
| 422 | Validation Error | Show field-level errors |
| 429 | Rate Limited | Show "too many requests" |
| 500 | Server Error | Show generic error |
| 502 | Bad Gateway | Retry with backoff |
| 503 | Service Unavailable | Retry with backoff |

### Client Error Wrapper
```kotlin
sealed class ApiResult<T> {
    data class Success<T>(val data: T) : ApiResult<T>()
    data class Error<T>(val message: String, val code: Int?) : ApiResult<T>()
    data class NetworkError<T>(val throwable: Throwable) : ApiResult<T>()
}
```

---

## 21. Currently Mocked — Needs Backend

Summary of features that have client-side UI but use hardcoded/mock data:

| Feature | Screen | What's Mocked | Backend Needed |
|---------|--------|---------------|----------------|
| **Wallet balances** | WalletScreen | Networks, balances, prices, address | Integrate with `/api/wallet/*` endpoints |
| **Send Token** | SendTokenScreen | Coin list, networks, commission | Token list, dynamic fees, tx execution |
| **QR address** | QRScreen | Address is hardcoded | Use real address from WalletManager |
| **Mana Points** | ManaPointsScreen | Balance (500 MP), income, rewards | New `/api/wallet/mana-points` endpoints |
| **Tariff purchase** | ChooseTarifScreen | 4 hardcoded tariff options | New `/api/wallet/mana-points/tariffs` |
| **Notifications** | NotificationsInboxScreen | 4 mock notifications | New `/api/notifications` endpoints |
| **News feed** | NewsScreen | 2 mock posts | New `/api/news` endpoints |
| **Calls** | CallsScreen | Placeholder "coming soon" | Full WebRTC infra (TURN/STUN) |
| **App banner upload** | CreateAppScreen | Upload UI but no function | `POST /api/developer/upload` endpoint |
| **Categories** | CreateAppScreen | 8 hardcoded categories | Dynamic from `/api/apps/categories` |

### Priority Order
1. **Wallet integration** — highest user impact, APIs exist
2. **Mana Points** — core monetization, needs new endpoints
3. **Notifications inbox** — engagement feature, needs new endpoints
4. **News feed** — social feature, needs new endpoints
5. **File upload** — developer feature, unblocks app icon/banner
6. **Calls** — infrastructure-heavy, lower priority

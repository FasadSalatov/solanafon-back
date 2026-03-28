# Solafon — Frontend API Integration Spec

> Ready-to-use API specification for the Android/iOS app.
> Base URL: `https://api.solafon.com`
> Generated: 2026-03-02

---

## Quick Start

### Base Configuration
```
Base URL:        https://api.solafon.com
Auth Header:     Authorization: Bearer {accessToken}
Content-Type:    application/json
Connect Timeout: 30s
Read Timeout:    120s
Write Timeout:   60s
```

### Token Storage
Store in EncryptedSharedPreferences (`solafon_secure_prefs`):
- `access_token` — JWT access token
- `refresh_token` — refresh token for auto-renewal
- `user` — JSON string of user object
- `is_authenticated` — boolean

### Auto Token Refresh
On any `401` response → call `POST /api/auth/refresh` with saved refreshToken → retry original request. If refresh fails → clear all auth data → redirect to login.

---

## Public Endpoints (No Auth Required)

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/health` | Health check |
| POST | `/api/auth/send-code` | Send email OTP |
| POST | `/api/auth/verify-code` | Verify OTP & login |
| POST | `/api/auth/refresh` | Refresh access token |
| GET | `/api/apps` | List marketplace apps |
| GET | `/api/apps/categories` | Get app categories |
| GET | `/api/i18n/languages` | Get supported languages |
| GET | `/api/support/faq` | Get FAQ |
| GET | `/api/legal/terms` | Terms of Service |
| GET | `/api/legal/privacy` | Privacy Policy |
| GET | `/api/referral/validate/:code` | Validate referral code |
| POST | `/api/crash/mobile` | Report crash |

---

## 1. Authentication

### 1.1 Send Verification Code
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
  "expiresIn": 600
}
```

### 1.2 Verify Code & Login
```
POST /api/auth/verify-code
```
**Request:**
```json
{
  "email": "user@example.com",
  "code": "123456",
  "referralCode": "ABCD1234"   // optional
}
```
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

### 1.3 Refresh Token
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

### 1.4 Logout
```
POST /api/auth/logout
Authorization: Bearer {token}
```
**Response:**
```json
{ "success": true }
```

---

## 2. User Management

### 2.1 Get My Profile
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
    "appsCount": 2,
    "transactionsCount": 12,
    "nftsCount": 0
  },
  "settings": {
    "hapticFeedback": true,
    "pushNotifications": true,
    "emailNotifications": true,
    "marketingEmails": false,
    "biometricEnabled": false
  }
}
```

### 2.2 Update Profile
```
PATCH /api/users/me
Authorization: Bearer {token}
```
**Request (all fields optional):**
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
  "user": { "id": "user_123", "email": "...", "displayName": "...", "avatarUrl": "..." }
}
```

### 2.3 Update Settings
```
PUT /api/users/me/settings
Authorization: Bearer {token}
```
**Request (send only changed fields):**
```json
{
  "hapticFeedback": true,
  "pushNotifications": false,
  "emailNotifications": true,
  "marketingEmails": false,
  "biometricEnabled": true
}
```
**Response:**
```json
{ "success": true, "settings": { ... } }
```

### 2.4 Delete Account
```
DELETE /api/users/me
Authorization: Bearer {token}
```
**Response:** `{ "success": true }`

### 2.5 Get Sessions
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
      "lastActive": "2025-01-01T12:00:00Z",
      "isCurrent": true
    }
  ]
}
```

### 2.6 Revoke Session
```
DELETE /api/users/me/sessions/{sessionId}
Authorization: Bearer {token}
```

### 2.7 Revoke All Sessions
```
DELETE /api/users/me/sessions
Authorization: Bearer {token}
```
Client clears local auth data on success.

---

## 3. Apps Marketplace

### 3.1 List Apps
```
GET /api/apps?category={cat}&search={q}&page={p}&limit={l}&sortBy={sort}
```
| Param | Type | Default | Values |
|-------|------|---------|--------|
| category | string? | null | category slug or null for all |
| search | string? | null | free text |
| page | int | 1 | |
| limit | int | 20 | |
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
**Auth:** Not required.

### 3.2 Get Categories
```
GET /api/apps/categories
```
**Response:**
```json
{
  "categories": [
    { "id": "ai", "name": "AI", "icon": "🤖", "appsCount": 15 }
  ]
}
```

### 3.3 Get App Detail
```
GET /api/apps/{appId}
```
**Response:**
```json
{
  "id": "app_123",
  "name": "My App",
  "description": "Short description",
  "longDescription": "Full markdown...",
  "icon": "🤖",
  "iconUrl": "https://...",
  "category": "AI",
  "url": "https://...",
  "users": "1.2K",
  "usersCount": 1200,
  "isVerified": true,
  "isTrending": false,
  "rating": 4.5,
  "screenshots": ["https://..."],
  "tags": ["ai", "chatbot"],
  "reviewsCount": 42,
  "permissions": ["wallet_read"],
  "version": "1.0.0",
  "developer": { "id": "dev_456", "name": "Dev", "isVerified": true }
}
```

### 3.4 Launch App
```
POST /api/apps/{appId}/launch
Authorization: Bearer {token}
```
**Response:**
```json
{ "success": true, "sessionId": "session_abc" }
```
Records usage analytics. Client opens app URL in WebView.

---

## 4. Developer Platform

### 4.1 List My Apps
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
      "description": "...",
      "icon": "🤖",
      "category": "AI",
      "url": "https://...",
      "status": "approved",
      "isVerified": false,
      "users": "100",
      "usersCount": 100,
      "rating": 4.2,
      "createdAt": "...",
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

### 4.2 Create App
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
**IMPORTANT:** `apiKey` is returned ONLY ONCE at creation. Show in a modal for user to copy.

### 4.3 Get App Detail
```
GET /api/developer/apps/{appId}
Authorization: Bearer {token}
```

### 4.4 Update App
```
PUT /api/developer/apps/{appId}
Authorization: Bearer {token}
```
**Request (all optional):**
```json
{
  "name": "Updated Name",
  "description": "...",
  "icon": "🎮",
  "category": "Games",
  "url": "https://..."
}
```

### 4.5 Delete App
```
DELETE /api/developer/apps/{appId}
Authorization: Bearer {token}
```

### 4.6 Generate API Key
```
POST /api/developer/apps/{appId}/api-keys
Authorization: Bearer {token}
```
**Response:**
```json
{
  "success": true,
  "apiKey": "sk_live_newkey...",
  "keyId": "key_new"
}
```

### 4.7 List API Credentials
```
GET /api/developer/apps/{appId}/api-keys
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
      "isActive": true,
      "createdAt": "...",
      "lastUsedAt": "..."
    }
  ]
}
```

### 4.8 Revoke API Key
```
DELETE /api/developer/apps/{appId}/api-keys/{keyId}
Authorization: Bearer {token}
```

### 4.9 Update Webhook
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

### 4.10 Get Welcome Message
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
    "content": "Hello! How can I help?",
    "bannerUrl": null,
    "isActive": true,
    "createdAt": "..."
  }
}
```

### 4.11 Update Welcome Message
```
PUT /api/developer/apps/{appId}/welcome-message
Authorization: Bearer {token}
```
**Request:**
```json
{
  "content": "Welcome! Ask me anything.",
  "bannerUrl": "https://...",
  "isActive": true
}
```

### 4.12 File Upload
```
POST /api/upload
Authorization: Bearer {token}
Content-Type: multipart/form-data
```
**Form Data:**
- `file` — binary file
- `type` — optional: `"icon"` | `"banner"`

**Response:**
```json
{
  "url": "https://api.solafon.com/uploads/abc123.png",
  "filename": "abc123.png",
  "size": 245000,
  "contentType": "image/png"
}
```
**Constraints:** Max 5MB. Formats: PNG, JPG, WEBP, GIF.

---

## 5. Chat & Conversations

### 5.1 List Conversations
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
      "appName": "My Bot",
      "appIcon": "🤖",
      "lastMessage": {
        "id": "msg_abc",
        "content": "Hello!",
        "senderType": "bot",
        "timestamp": "2025-01-01T12:00:00Z"
      },
      "unreadCount": 2,
      "createdAt": "...",
      "updatedAt": "..."
    }
  ],
  "pagination": { "currentPage": 1, "totalPages": 5, "totalItems": 100, "hasMore": true }
}
```

### 5.2 Start Conversation
```
POST /api/conversations
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
    "content": "Welcome!",
    "senderType": "bot",
    "timestamp": "..."
  }
}
```

### 5.3 Get Messages
```
GET /api/conversations/{conversationId}/messages?limit={l}&before={cursor}
Authorization: Bearer {token}
```
| Param | Type | Default | Description |
|-------|------|---------|-------------|
| limit | int | 50 | Messages per page |
| before | string? | null | Message ID cursor |

**Response:**
```json
{
  "messages": [
    {
      "id": "msg_123",
      "conversationId": "conv_789",
      "senderId": "user_123",
      "senderType": "user",
      "content": "Hello!",
      "timestamp": "2025-01-01T12:00:00Z",
      "status": "read",
      "metadata": null
    }
  ],
  "hasMore": true
}
```

**Content types for bot messages:**
- Text: `{ "type": "text", "text": "..." }`
- Buttons: `{ "type": "button", "text": "Choose:", "buttons": [{ "id": "btn_1", "text": "Option A", "action": "callback", "payload": "opt_a" }] }`

### 5.4 Send Message
```
POST /api/conversations/{conversationId}/messages
Authorization: Bearer {token}
```
**Request:**
```json
{
  "content": "Hello bot!",
  "replyToId": null
}
```
**Response:**
```json
{ "success": true, "message": { ... } }
```

### 5.5 Button Callback
```
POST /api/conversations/{conversationId}/messages/{messageId}/callback
Authorization: Bearer {token}
```
**Request:**
```json
{
  "buttonId": "btn_1",
  "payload": "opt_a"
}
```

### 5.6 Mark as Read
```
POST /api/conversations/{conversationId}/read
Authorization: Bearer {token}
```

### 5.7 Delete Conversation
```
DELETE /api/conversations/{conversationId}
Authorization: Bearer {token}
```

---

## 6. Wallet & Blockchain

### Architecture
- **Non-custodial** — private keys never leave the device
- **Chain:** Solana mainnet-beta
- **Key derivation:** Ed25519 via BIP39/BIP44/SLIP-0010 (`m/44'/501'/0'/0'`)
- Transaction signing happens locally, server only broadcasts

### 6.1 Get Balance
```
GET /api/wallet/balance?address={solanaAddress}
Authorization: Bearer {token}
```
**Response:**
```json
{
  "address": "DRpbCBMxVnDK...",
  "sol": 1.5,
  "lamports": 1500000000,
  "tokens": [
    {
      "mint": "EPjFWdd5...",
      "symbol": "USDC",
      "name": "USD Coin",
      "balance": "100.00",
      "decimals": 6,
      "uiAmount": 100.0
    }
  ]
}
```

### 6.2 Send Transaction
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
**IMPORTANT:** Transaction must be signed locally on device. Server only broadcasts.

### 6.3 Transaction History
```
GET /api/wallet/transactions?address={addr}&limit={l}
Authorization: Bearer {token}
```
**Response:**
```json
{
  "transactions": [
    {
      "signature": "5KtP3...",
      "slot": 123456789,
      "blockTime": 1704067200,
      "confirmationStatus": "finalized"
    }
  ]
}
```

### 6.4 Supported Tokens
```
GET /api/wallet/tokens
Authorization: Bearer {token}
```
**Response:**
```json
{
  "tokens": [
    { "mint": "So11111...", "symbol": "SOL", "name": "Solana", "decimals": 9, "icon": "https://..." }
  ]
}
```

### 6.5 Token Prices
```
GET /api/wallet/prices?mints={comma_separated}
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

### 6.6 Transaction Status
```
GET /api/wallet/transactions/{signature}/status
Authorization: Bearer {token}
```
**Response:**
```json
{ "signature": "5KtP3...", "status": "finalized" }
```
Statuses: `"pending"` → `"confirmed"` → `"finalized"` | `"failed"`

---

## 7. Mana Points

### 7.1 Get Balance & History
```
GET /api/mana
Authorization: Bearer {token}
```
**Response:**
```json
{
  "balance": 500,
  "income": 500,
  "rewards": 12,
  "history": [
    {
      "id": "tx_123",
      "type": "purchase",
      "amount": 500,
      "description": "Purchased 500 MP",
      "createdAt": "..."
    }
  ]
}
```

### 7.2 Get Tariffs
```
GET /api/mana/tariffs
Authorization: Bearer {token}
```
**Response:**
```json
{
  "tariffs": [
    { "id": "tariff_1", "name": "Starter", "mpAmount": 250, "priceUSDT": 5.79, "isPopular": false },
    { "id": "tariff_2", "name": "Popular", "mpAmount": 500, "priceUSDT": 11.90, "isPopular": true },
    { "id": "tariff_3", "name": "Pro", "mpAmount": 1000, "priceUSDT": 22.99, "isPopular": false },
    { "id": "tariff_4", "name": "Max", "mpAmount": 2500, "priceUSDT": 57.98, "isPopular": false }
  ]
}
```

### 7.3 Purchase Mana Points
```
POST /api/mana/purchase
Authorization: Bearer {token}
```
**Request:**
```json
{ "tariffId": "tariff_2" }
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

### 7.4 Gift Mana Points
```
POST /api/mana/gift
Authorization: Bearer {token}
```
**Request:**
```json
{
  "recipientAddress": "user_456",
  "amount": 100
}
```
**Response:**
```json
{
  "success": true,
  "transactionId": "tx_gift_abc",
  "newBalance": 400,
  "amountSent": 100
}
```

### 7.5 Get Wallet Networks
```
GET /api/mana/networks
Authorization: Bearer {token}
```
**Response:**
```json
{
  "networks": [
    { "id": "solana", "name": "Solana", "iconColor": "#9945FF", "isDefault": true },
    { "id": "ethereum", "name": "Ethereum", "iconColor": "#627EEA", "isDefault": false },
    { "id": "polygon", "name": "Polygon", "iconColor": "#8247E5", "isDefault": false }
  ]
}
```

---

## 8. Referral System

### 8.1 Get Referral Info
```
GET /api/referral
Authorization: Bearer {token}
```
**Response:**
```json
{
  "data": {
    "code": "ABC12345",
    "link": "https://solafon.com/ref/ABC12345",
    "totalCount": 5,
    "referrals": [
      {
        "id": "user_456",
        "displayName": "John",
        "avatarUrl": "https://...",
        "registeredAt": "2025-01-15T00:00:00Z"
      }
    ]
  }
}
```

### 8.2 Validate Referral Code
```
GET /api/referral/validate/{code}
```
**Response:**
```json
{ "data": { "valid": true, "referrerName": "John" } }
```
or
```json
{ "data": { "valid": false } }
```

---

## 9. WebRTC Calls

### 9.1 Create Room
```
POST /api/calls/rooms/create
Authorization: Bearer {token}
```
**Request:**
```json
{
  "type": "direct",
  "participantIds": ["user_456"]
}
```
**Response:**
```json
{
  "data": {
    "room": {
      "id": "room_123",
      "roomCode": "ABC123",
      "type": "direct",
      "status": "waiting",
      "createdBy": "user_123",
      "participants": [ ... ],
      "createdAt": "..."
    },
    "deepLink": "solafon://call/ABC123",
    "roomCode": "ABC123"
  }
}
```

### 9.2 Get Room by Code
```
GET /api/calls/rooms/code/{roomCode}
Authorization: Bearer {token}
```

### 9.3 Join Room
```
POST /api/calls/rooms/{roomId}/join
Authorization: Bearer {token}
```

### 9.4 End Call
```
POST /api/calls/rooms/{roomId}/end
Authorization: Bearer {token}
```

### 9.5 Update Participant Status
```
PATCH /api/calls/rooms/{roomId}/status
Authorization: Bearer {token}
```
**Request:**
```json
{
  "isMuted": true,
  "isVideoOn": false,
  "isAudioOn": true
}
```

---

## 10. Notifications

### 10.1 Register Push Token
```
POST /api/notifications/push-token
Authorization: Bearer {token}
```
**Request:**
```json
{
  "token": "fcm_token_string",
  "platform": "android",
  "deviceId": "device_uuid"
}
```

### 10.2 Unregister Push Token
```
DELETE /api/notifications/push-token
Authorization: Bearer {token}
```
**Request:**
```json
{ "token": "fcm_token_string" }
```

### 10.3 List Notifications
```
GET /api/notifications?page={p}&limit={l}
Authorization: Bearer {token}
```
**Response:**
```json
{
  "notifications": [
    {
      "id": "notif_123",
      "type": "message",
      "title": "New message from Bot",
      "body": "Hello!",
      "data": { "conversationId": "conv_456" },
      "isRead": false,
      "createdAt": "..."
    }
  ],
  "pagination": { "currentPage": 1, "totalPages": 3, "totalItems": 50, "hasMore": true }
}
```

### 10.4 Mark as Read
```
POST /api/notifications/{notificationId}/read
Authorization: Bearer {token}
```

### 10.5 Mark All as Read
```
POST /api/notifications/read-all
Authorization: Bearer {token}
```

### 10.6 Get Unread Count
```
GET /api/notifications/unread-count
Authorization: Bearer {token}
```
**Response:**
```json
{ "unreadCount": 5 }
```

---

## 11. News Feed

### 11.1 Get Feed
```
GET /api/news/feed?page={p}&limit={l}
Authorization: Bearer {token}
```
**Response:**
```json
{
  "posts": [
    {
      "id": "post_123",
      "title": "New Feature!",
      "content": "We launched...",
      "imageUrl": "https://...",
      "author": { "id": "user_1", "displayName": "Admin", "avatarUrl": "..." },
      "likesCount": 42,
      "commentsCount": 5,
      "sharesCount": 3,
      "isLiked": false,
      "createdAt": "..."
    }
  ],
  "pagination": { ... }
}
```

### 11.2 Like Post
```
POST /api/news/{postId}/like
Authorization: Bearer {token}
```
**Response:**
```json
{ "success": true, "liked": true, "likesCount": 43 }
```
Toggle — calling again unlikes.

### 11.3 Share Post
```
POST /api/news/{postId}/share
Authorization: Bearer {token}
```
**Response:**
```json
{ "success": true, "sharesCount": 4 }
```

### 11.4 Get Comments
```
GET /api/news/{postId}/comments?page={p}&limit={l}
Authorization: Bearer {token}
```
**Response:**
```json
{
  "comments": [
    {
      "id": "comment_123",
      "content": "Great feature!",
      "author": { "id": "user_456", "displayName": "John", "avatarUrl": "..." },
      "createdAt": "..."
    }
  ],
  "pagination": { ... }
}
```

### 11.5 Post Comment
```
POST /api/news/{postId}/comments
Authorization: Bearer {token}
```
**Request:**
```json
{ "content": "Great feature!" }
```

### 11.6 Create Post (Admin)
```
POST /api/news
Authorization: Bearer {token}
```
**Request:**
```json
{
  "title": "New Feature",
  "content": "We launched...",
  "imageUrl": "https://..."
}
```

---

## 12. Support & Legal

### 12.1 Get FAQ
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

### 12.2 Create Support Ticket
```
POST /api/support/tickets
Authorization: Bearer {token}
```
**Request:**
```json
{
  "subject": "Can't send tokens",
  "message": "I'm trying to...",
  "category": "wallet"
}
```
**Response:**
```json
{ "ticketId": "ticket_123", "status": "open" }
```

### 12.3 Get Terms of Service
```
GET /api/legal/terms
```
**Response:**
```json
{ "content": "# Terms...", "version": "1.0", "lastUpdated": "2026-03-01" }
```

### 12.4 Get Privacy Policy
```
GET /api/legal/privacy
```

---

## 13. Internationalization

### 13.1 Get Supported Languages
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

## 14. Crash Reporting

### 14.1 Report Crash
```
POST /api/crash/mobile
```
**Request:**
```json
{
  "userId": "user_123",
  "sessionId": "sess_abc",
  "errorType": "crash",
  "message": "NullPointerException at...",
  "stack": "at com.solafon...",
  "userAgent": "Android 14; Pixel 7",
  "source": "android",
  "metadata": { "screen": "wallet", "action": "send" }
}
```
**Response:**
```json
{
  "data": {
    "id": "crash_123",
    "appId": "solafon-android",
    "errorType": "crash",
    "message": "NullPointerException at...",
    "createdAt": "..."
  }
}
```

---

## 15. WebSocket Protocol

### 15.1 Chat WebSocket
```
URL: wss://api.solafon.com/ws?token={accessToken}
```

| Config | Value |
|--------|-------|
| Ping interval | 30s |
| Max reconnect attempts | 5 |
| Reconnect base delay | 5s (exponential backoff) |
| Fallback | HTTP polling |

**Server → Client Events:**
| Event | Payload |
|-------|---------|
| `new_message` | `ChatMessage` object |
| `typing` | `{ conversationId, isTyping, timestamp }` |
| `message_status` | `{ messageId, status }` |
| `conversation_update` | `{ conversationId, unreadCount }` |
| `pong` | `{}` |

**Client → Server Events:**
| Event | Payload |
|-------|---------|
| `ping` | `{}` |
| `subscribe` | `{ type: "subscribe", conversationId }` |
| `unsubscribe` | `{ type: "unsubscribe", conversationId }` |

**Polling Fallback** (when WS unavailable):
- Active: poll every 3s via `GET /api/conversations/{id}/messages?limit=20`
- Idle: poll every 10s

### 15.2 WebRTC Signaling WebSocket
```
URL: wss://api.solafon.com/ws/webrtc?roomId={roomId}
Auth: Authorization: Bearer {token}
```

**Signaling events:** `offer`, `answer`, `ice-candidate`, `user-joined`, `user-left`, `call-ended`

---

## 16. Error Format

All new `/api/` endpoints use structured errors:
```json
{
  "error": {
    "code": "ERROR_CODE",
    "message": "Human-readable message"
  }
}
```

**Common error codes:**
| Code | HTTP Status | Description |
|------|-------------|-------------|
| `VALIDATION_ERROR` | 400 | Invalid request body |
| `UNAUTHORIZED` | 401 | Missing or invalid token |
| `NOT_FOUND` | 404 | Resource not found |
| `INSUFFICIENT_BALANCE` | 400 | Not enough Mana Points |
| `INTERNAL_ERROR` | 500 | Server error |

---

## 17. ID Format

All entity IDs in responses are prefixed strings:
| Entity | Format | Example |
|--------|--------|---------|
| User | `user_{id}` | `user_123` |
| App | `app_{id}` | `app_456` |
| Conversation | `conv_{id}` | `conv_789` |
| Message | `msg_{id}` | `msg_abc` |
| Room | `room_{id}` | `room_123` |
| Notification | `notif_{id}` | `notif_456` |
| Transaction | `tx_{id}` | `tx_789` |

---

## 18. Currently Mocked (Backend Ready, Frontend Integration Needed)

These features have backend endpoints ready but may need frontend integration work:

| Feature | Status | Endpoint |
|---------|--------|----------|
| WebSocket chat | Backend: stub (no WS server yet) | `wss://api.solafon.com/ws` |
| WebRTC signaling | Backend: stub | `wss://api.solafon.com/ws/webrtc` |
| FCM push notifications | Backend: stores tokens | `POST /api/notifications/push-token` |
| File upload resize | Backend: saves file, no resize | `POST /api/upload` |
| Mana purchase payment | Backend: mock (instant credit) | `POST /api/mana/purchase` |
| Token prices | Backend: hardcoded | `GET /api/wallet/prices` |

---

## Legacy API (`/api/v1/`)

The original `/api/v1/` endpoints remain available for backward compatibility:
- `/api/v1/auth/*` — original email OTP auth
- `/api/v1/apps/*` — basic app CRUD
- `/api/v1/profile/*` — profile management
- `/api/v1/mana/*` — mana points
- `/api/v1/secret/*` — secret login
- `/api/v1/bot/*` — bot API (token auth)

New frontend should use `/api/*` endpoints exclusively.

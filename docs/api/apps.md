# Apps API

Endpoints for browsing and managing mini-apps.

## Get All Apps

List all approved apps with optional category filter.

**Endpoint:** `GET /apps`

**Query Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| category | string | Filter by category slug (optional) |

**Headers:** `Authorization: Bearer <token>` (required)

**Response:**
```json
{
  "apps": [
    {
      "id": 1,
      "title": "Crypto Quest",
      "subtitle": "Play and earn in adventure RPG",
      "icon": "🎮",
      "categoryId": 2,
      "category": {"id": 2, "name": "Games", "slug": "games"},
      "isSecret": false,
      "isVerified": true,
      "usersCount": 45000,
      "description": "An exciting crypto gaming experience...",
      "botUsername": "cryptoquestapp",
      "moderationStatus": "approved"
    }
  ],
  "total": 10
}
```

---

## Search Apps

Search apps by title or description.

**Endpoint:** `GET /apps/search`

**Query Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| q | string | Search query (required) |

**Headers:** `Authorization: Bearer <token>` (required)

**Response:** Same format as Get All Apps.

---

## Get App by ID

Get detailed information about a specific app.

**Endpoint:** `GET /apps/:id`

**Headers:** `Authorization: Bearer <token>` (required)

**Response:**
```json
{
  "id": 1,
  "title": "Crypto Quest",
  "subtitle": "Play and earn in adventure RPG",
  "icon": "🎮",
  "categoryId": 2,
  "category": {"id": 2, "name": "Games", "slug": "games"},
  "isSecret": false,
  "isVerified": true,
  "usersCount": 45000,
  "description": "An exciting crypto gaming experience...",
  "botUsername": "cryptoquestapp",
  "welcomeMessage": "Welcome to Crypto Quest!",
  "moderationStatus": "approved",
  "createdAt": "2024-01-01T00:00:00Z"
}
```

---

## Create New App

Create a new app (will require moderation).

**Endpoint:** `POST /apps`

**Headers:** `Authorization: Bearer <token>` (required)

**Request:**
```json
{
  "icon": "🎮",
  "title": "My App",
  "description": "App description here",
  "categoryId": 2,
  "botUsername": "myawesomeapp",
  "welcomeMessage": "Welcome to My App!",
  "webhookUrl": "https://myserver.com/webhook"
}
```

**Response:**
```json
{
  "message": "App created successfully. It will be reviewed by moderators within 24 hours.",
  "app": {
    "id": 123,
    "title": "My App",
    "icon": "🎮",
    "botUsername": "myawesomeapp",
    "moderationStatus": "pending"
  },
  "apiToken": "abc123def456..."
}
```

> **Important:** Save the `apiToken` securely - it's used for Developer API authentication.

---

## Get My Apps

List apps created by the authenticated user.

**Endpoint:** `GET /apps/my`

**Headers:** `Authorization: Bearer <token>` (required)

**Response:**
```json
{
  "apps": [
    {
      "id": 123,
      "title": "My App",
      "icon": "🎮",
      "botUsername": "myawesomeapp",
      "moderationStatus": "pending",
      "usersCount": 0,
      "createdAt": "2024-01-01T00:00:00Z"
    }
  ],
  "total": 1
}
```

---

## Update App

Update an app you own.

**Endpoint:** `PUT /apps/:id`

**Headers:** `Authorization: Bearer <token>` (required)

**Request:**
```json
{
  "title": "Updated Title",
  "description": "Updated description",
  "welcomeMessage": "New welcome message"
}
```

**Response:**
```json
{
  "message": "App updated successfully",
  "app": {...}
}
```

> **Note:** Updating title or description will reset moderation status to pending.

---

## Delete App

Delete an app you own.

**Endpoint:** `DELETE /apps/:id`

**Headers:** `Authorization: Bearer <token>` (required)

**Response:**
```json
{
  "message": "App deleted successfully"
}
```

---

## Get App Settings (Owner Only)

Get app settings including API token.

**Endpoint:** `GET /apps/:id/settings`

**Headers:** `Authorization: Bearer <token>` (required)

**Response:**
```json
{
  "app": {
    "id": 123,
    "title": "My App",
    "apiToken": "abc123...",
    "webhookUrl": "https://myserver.com/webhook",
    "botUsername": "myawesomeapp",
    "welcomeMessage": "Welcome!",
    "moderationStatus": "approved",
    "usersCount": 1500
  },
  "commands": [
    {"id": 1, "command": "/start", "description": "Start", "response": "Welcome!"}
  ],
  "webhookLogs": [...]
}
```

---

## Regenerate API Token

Generate a new API token (invalidates the old one).

**Endpoint:** `POST /apps/:id/regenerate-token`

**Headers:** `Authorization: Bearer <token>` (required)

**Response:**
```json
{
  "message": "API token regenerated successfully",
  "apiToken": "new-token-123..."
}
```

# Profile API

Endpoints for user profile management.

## Get Full Profile

Get complete user profile with stats.

**Endpoint:** `GET /profile`

**Headers:** `Authorization: Bearer <token>` (required)

**Response:**
```json
{
  "id": 1,
  "email": "user@example.com",
  "name": "John Doe",
  "avatar": "",
  "manaPoints": {
    "balance": 500,
    "valueUSD": 11.50
  },
  "hasSecretAccess": false,
  "language": "en",
  "notificationsEnabled": true,
  "twoFactorEnabled": false,
  "stats": {
    "appsCount": 2,
    "transactionsCount": 5,
    "nftsCount": 0
  },
  "createdAt": "2024-01-01T00:00:00Z"
}
```

---

## Update Profile

Update user profile information.

**Endpoint:** `PATCH /profile`

**Headers:** `Authorization: Bearer <token>` (required)

**Request:**
```json
{
  "name": "New Name",
  "avatar": "https://example.com/avatar.png",
  "language": "ru"
}
```

**Response:**
```json
{
  "message": "Profile updated successfully",
  "user": {
    "id": 1,
    "name": "New Name",
    "avatar": "https://example.com/avatar.png",
    "language": "ru"
  }
}
```

### Available Fields

| Field | Type | Description |
|-------|------|-------------|
| name | string | Display name |
| avatar | string | Avatar URL |
| language | string | Preferred language (en, ru) |

---

## Update Settings

Update user settings.

**Endpoint:** `PUT /profile/settings`

**Headers:** `Authorization: Bearer <token>` (required)

**Request:**
```json
{
  "notificationsEnabled": true,
  "twoFactorEnabled": false
}
```

**Response:**
```json
{
  "message": "Settings updated successfully",
  "settings": {
    "notificationsEnabled": true,
    "twoFactorEnabled": false
  }
}
```

### Available Settings

| Setting | Type | Description |
|---------|------|-------------|
| notificationsEnabled | boolean | Enable push notifications |
| twoFactorEnabled | boolean | Enable two-factor authentication |

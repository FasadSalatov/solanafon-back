# Auth API

Endpoints for user authentication and session management.

## Request OTP Code

Send OTP code to email.

**Endpoint:** `POST /auth/email/request`

**Request:**
```json
{
  "email": "user@example.com"
}
```

**Response:**
```json
{
  "message": "OTP sent to your email",
  "email": "user@example.com"
}
```

---

## Verify OTP Code

Verify OTP and get JWT token. New users receive 500 MP welcome bonus.

**Endpoint:** `POST /auth/email/verify`

**Request:**
```json
{
  "email": "user@example.com",
  "code": "123456"
}
```

**Response (new user):**
```json
{
  "token": "eyJhbGciOiJIUzI1NiIs...",
  "user": {
    "id": 1,
    "email": "user@example.com",
    "name": "Solana User",
    "manaPoints": 500,
    "hasSecretAccess": false,
    "language": "en",
    "createdAt": "2024-01-01T00:00:00Z"
  },
  "isNewUser": true,
  "welcomeBonus": 500
}
```

**Response (existing user):**
```json
{
  "token": "eyJhbGciOiJIUzI1NiIs...",
  "user": {
    "id": 1,
    "email": "user@example.com",
    "name": "John Doe",
    "manaPoints": 1500,
    "hasSecretAccess": true,
    "language": "en",
    "createdAt": "2024-01-01T00:00:00Z"
  },
  "isNewUser": false
}
```

---

## Get Current User

Get authenticated user info with stats.

**Endpoint:** `GET /auth/me`

**Headers:** `Authorization: Bearer <token>` (required)

**Response:**
```json
{
  "id": 1,
  "email": "user@example.com",
  "name": "John Doe",
  "avatar": "",
  "manaPoints": 500,
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

## Logout

Logout current session (client-side token invalidation).

**Endpoint:** `POST /auth/logout`

**Headers:** `Authorization: Bearer <token>` (required)

**Response:**
```json
{
  "message": "Logged out successfully"
}
```

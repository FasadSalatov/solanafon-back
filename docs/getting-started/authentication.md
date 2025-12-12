# Authentication

Solafon uses two types of authentication:

1. **User Authentication (JWT)** - For end-user actions (browsing apps, profile, etc.)
2. **API Token Authentication** - For app developers to send/receive messages

## User Authentication

### 1. Request OTP

Send an OTP code to the user's email address.

**Endpoint:** `POST /auth/email/request`

```bash
curl -X POST https://api.solafon.com/api/v1/auth/email/request \
  -H "Content-Type: application/json" \
  -d '{"email": "user@example.com"}'
```

**Response:**
```json
{
  "message": "OTP sent to your email",
  "email": "user@example.com"
}
```

### 2. Verify OTP

Verify the OTP code and receive a JWT token.

**Endpoint:** `POST /auth/email/verify`

```bash
curl -X POST https://api.solafon.com/api/v1/auth/email/verify \
  -H "Content-Type: application/json" \
  -d '{"email": "user@example.com", "code": "123456"}'
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

New users receive **500 Mana Points** as a welcome bonus.

### 3. Using JWT Token

Include the token in the Authorization header for all protected endpoints:

```bash
curl https://api.solafon.com/api/v1/profile \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIs..."
```

### 4. Get Current User

**Endpoint:** `GET /auth/me`

```bash
curl https://api.solafon.com/api/v1/auth/me \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

### 5. Logout

**Endpoint:** `POST /auth/logout`

Token invalidation is handled client-side.

---

## API Token Authentication (For Developers)

When you create an app, you receive an API token. Use this token for the Developer API.

### Using API Token

```bash
curl https://api.solafon.com/api/v1/bot/getMe \
  -H "Authorization: Bearer YOUR_API_TOKEN"
```

Or without the "Bearer" prefix:

```bash
curl https://api.solafon.com/api/v1/bot/getMe \
  -H "Authorization: YOUR_API_TOKEN"
```

### Regenerating API Token

If your token is compromised, regenerate it:

**Endpoint:** `POST /apps/:id/regenerate-token`

```bash
curl -X POST https://api.solafon.com/api/v1/apps/123/regenerate-token \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

**Response:**
```json
{
  "message": "API token regenerated successfully",
  "apiToken": "new-token-abc123..."
}
```

## Security Best Practices

1. **Never expose tokens in client-side code** - Always use server-side requests
2. **Store tokens securely** - Use environment variables or secret managers
3. **Regenerate compromised tokens immediately** - Use the regenerate endpoint
4. **Use HTTPS only** - All API calls should use HTTPS

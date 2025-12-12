# API Overview

The Solafon API allows you to interact with the platform programmatically.

## Base URL

```
https://api.solafon.com/api/v1
```

For local development:
```
http://localhost:8080/api/v1
```

## Authentication

Most endpoints require authentication via JWT token:

```
Authorization: Bearer <jwt_token>
```

See [Authentication Guide](../getting-started/authentication.md) for details.

## Response Format

All responses are JSON with consistent structure:

### Success Response
```json
{
  "data": {...}
}
```

Or for collections:
```json
{
  "items": [...],
  "total": 42
}
```

### Error Response
```json
{
  "error": "Error message description"
}
```

## HTTP Status Codes

| Code | Description |
|------|-------------|
| 200 | Success |
| 201 | Created |
| 400 | Bad Request - Invalid input or validation error |
| 401 | Unauthorized - Missing or invalid token |
| 403 | Forbidden - No access rights |
| 404 | Not Found |
| 500 | Internal Server Error |

## Rate Limiting

Current limits:
- 100 requests per minute for authenticated users
- 20 requests per minute for unauthenticated endpoints

Rate limit headers:
```
X-RateLimit-Limit: 100
X-RateLimit-Remaining: 95
X-RateLimit-Reset: 1704067200
```

## API Sections

### User API (JWT Authentication)

| Section | Description |
|---------|-------------|
| [Auth](auth.md) | Authentication and user management |
| [Apps](apps.md) | Browse and interact with apps |
| [Messages](messages.md) | Chat with apps |
| [Categories](categories.md) | App categories |
| [Profile](profile.md) | User profile management |
| [Mana](mana.md) | Mana Points system |
| [Secret](secret.md) | Secret Login feature |

### Developer API (API Token Authentication)

| Section | Description |
|---------|-------------|
| [Overview](../developer-api/overview.md) | Developer API introduction |
| [Send Messages](../developer-api/send-message.md) | Send messages to users |
| [Receive Messages](../developer-api/receive-messages.md) | Get user messages |
| [Webhooks](../developer-api/webhooks.md) | Real-time message delivery |
| [Commands](../developer-api/commands.md) | Manage app commands |

## CORS

The API supports CORS for browser-based applications. Allowed origins can be configured on the server.

## Versioning

The API uses URL versioning. Current version: `v1`

Breaking changes will be introduced in new versions (e.g., `/api/v2`).

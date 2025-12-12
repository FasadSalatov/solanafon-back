# Developer API Authentication

## API Token

Every app has a unique API token used to authenticate Developer API requests.

### Getting Your Token

**Option 1: Via Dev Studio**
```bash
curl -X POST https://api.solafon.com/api/v1/apps/devstudio/message \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{"content": "/token"}'
```

**Option 2: Via Apps API**
```bash
curl https://api.solafon.com/api/v1/apps/YOUR_APP_ID/settings \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

### Using the Token

Include it in the Authorization header:

```bash
# With Bearer prefix
curl https://api.solafon.com/api/v1/bot/getMe \
  -H "Authorization: Bearer abc123def456..."

# Without prefix (also works)
curl https://api.solafon.com/api/v1/bot/getMe \
  -H "Authorization: abc123def456..."
```

### Regenerating Token

If your token is compromised:

```bash
curl -X POST https://api.solafon.com/api/v1/apps/YOUR_APP_ID/regenerate-token \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

Response:
```json
{
  "message": "API token regenerated successfully",
  "apiToken": "new-token-xyz..."
}
```

The old token is immediately invalidated.

## Verify Token Works

Test your token with getMe:

```bash
curl https://api.solafon.com/api/v1/bot/getMe \
  -H "Authorization: Bearer YOUR_API_TOKEN"
```

**Success:**
```json
{
  "ok": true,
  "result": {
    "id": 123,
    "is_bot": true,
    "title": "My App",
    "username": "myawesomeapp"
  }
}
```

**Invalid Token:**
```json
{
  "ok": false,
  "error_code": 401,
  "description": "Unauthorized"
}
```

## Security Best Practices

1. **Never commit tokens to git** - Use environment variables
2. **Don't expose in client code** - Keep tokens server-side only
3. **Rotate regularly** - Regenerate tokens periodically
4. **Monitor usage** - Check webhook logs for suspicious activity
5. **Use HTTPS only** - Never send tokens over HTTP

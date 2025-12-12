# Secret Login API

Secret Login provides enhanced privacy through virtual phone numbers.

## Get Available Numbers

List available virtual numbers for purchase.

**Endpoint:** `GET /secret/numbers`

**Headers:** `Authorization: Bearer <token>` (required)

**Response:**
```json
{
  "numbers": [
    {
      "id": 1,
      "number": "+999 (764) 123-45-67",
      "isPremium": false,
      "priceMP": 50,
      "isAvailable": true
    },
    {
      "id": 3,
      "number": "+999 (555) 000-11-22",
      "isPremium": true,
      "priceMP": 100,
      "isAvailable": true
    }
  ],
  "total": 8,
  "userManaPoints": 500
}
```

---

## Activate Secret Number

Purchase and activate a secret number using Mana Points.

**Endpoint:** `POST /secret/activate`

**Headers:** `Authorization: Bearer <token>` (required)

**Request:**
```json
{
  "numberId": 1
}
```

**Response:**
```json
{
  "message": "Secret access activated successfully",
  "secretAccess": {
    "id": 1,
    "userId": 1,
    "secretNumberId": 1,
    "activatedAt": "2024-01-01T00:00:00Z",
    "isActive": true
  },
  "newManaBalance": 450,
  "spent": 50
}
```

**Error (insufficient MP):**
```json
{
  "error": "Insufficient Mana Points",
  "required": 50,
  "balance": 30
}
```

---

## Get Secret Access Status

Check current secret access status.

**Endpoint:** `GET /secret/status`

**Headers:** `Authorization: Bearer <token>` (required)

**Response (with active access):**
```json
{
  "hasAccess": true,
  "manaPoints": 450,
  "activatedAt": "2024-01-01T00:00:00Z",
  "number": "+999 (764) 123-45-67",
  "isPremium": false,
  "benefits": [
    "Full anonymity and confidentiality",
    "Receive one-time SMS codes",
    "No link to real phone number",
    "Access to secret apps"
  ]
}
```

**Response (no access):**
```json
{
  "hasAccess": false,
  "manaPoints": 500
}
```

---

## Deactivate Secret Access

Deactivate current secret access (number becomes available again).

**Endpoint:** `DELETE /secret/deactivate`

**Headers:** `Authorization: Bearer <token>` (required)

**Response:**
```json
{
  "message": "Secret access deactivated"
}
```

## Number Pricing

| Type | Price (MP) |
|------|------------|
| Standard | 50 MP |
| Premium | 100-150 MP |

Premium numbers have memorable or special patterns.

## Benefits of Secret Login

- Complete anonymity
- Receive SMS verification codes
- No connection to personal phone
- Access to exclusive apps

# Mana Points API

Mana Points (MP) is the internal currency of Solafon platform.

## Get MP Balance & History

Get current balance and transaction history.

**Endpoint:** `GET /mana`

**Headers:** `Authorization: Bearer <token>` (required)

**Response:**
```json
{
  "balance": 500,
  "transactions": [
    {
      "id": 1,
      "userId": 1,
      "amount": 500,
      "type": "reward",
      "description": "Welcome bonus for new user",
      "createdAt": "2024-01-01T00:00:00Z"
    },
    {
      "id": 2,
      "userId": 1,
      "amount": -50,
      "type": "purchase",
      "description": "Secret number activation: +999 (764) 123-45-67",
      "createdAt": "2024-01-02T00:00:00Z"
    }
  ],
  "total": 2
}
```

---

## Top Up MP

Add Mana Points to your balance.

**Endpoint:** `POST /mana/topup`

**Headers:** `Authorization: Bearer <token>` (required)

**Request:**
```json
{
  "amount": 100
}
```

**Response:**
```json
{
  "message": "Top up successful",
  "newBalance": 600,
  "added": 100
}
```

## Transaction Types

| Type | Description |
|------|-------------|
| reward | Bonus or reward (welcome, referral, etc.) |
| purchase | Spending MP on features |
| topup | Adding MP to balance |
| refund | Returned MP |

## MP Pricing

| MP Amount | USD Value |
|-----------|-----------|
| 100 MP | $2.30 |
| 500 MP | $11.50 |
| 1000 MP | $23.00 |

## Uses for Mana Points

- **Secret Login** - Activate anonymous virtual numbers (50-150 MP)
- **Premium Features** - Access exclusive app features
- **Verified Badge** - Apply for app verification

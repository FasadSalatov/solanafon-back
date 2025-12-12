# Categories API

Endpoints for app categories.

## Get All Categories

List all available categories.

**Endpoint:** `GET /categories`

**Headers:** `Authorization: Bearer <token>` (required)

**Response:**
```json
{
  "categories": [
    {"id": 1, "name": "AI", "slug": "ai", "icon": "🤖", "order": 1},
    {"id": 2, "name": "Games", "slug": "games", "icon": "🎮", "order": 2},
    {"id": 3, "name": "Trading", "slug": "trading", "icon": "📊", "order": 3},
    {"id": 4, "name": "DePIN", "slug": "depin", "icon": "🌐", "order": 4},
    {"id": 5, "name": "DeFi", "slug": "defi", "icon": "💎", "order": 5},
    {"id": 6, "name": "NFT", "slug": "nft", "icon": "🖼️", "order": 6},
    {"id": 7, "name": "Staking", "slug": "staking", "icon": "🔒", "order": 7},
    {"id": 8, "name": "Services", "slug": "services", "icon": "🛠️", "order": 8}
  ],
  "total": 8
}
```

---

## Get Apps by Category

Get all apps in a specific category.

**Endpoint:** `GET /categories/:slug/apps`

**Headers:** `Authorization: Bearer <token>` (required)

**Example:** `GET /categories/games/apps`

**Response:**
```json
{
  "category": {
    "id": 2,
    "name": "Games",
    "slug": "games",
    "icon": "🎮"
  },
  "apps": [
    {
      "id": 1,
      "title": "Crypto Quest",
      "subtitle": "Play and earn",
      "icon": "🎮",
      "isVerified": true,
      "usersCount": 45000
    }
  ],
  "total": 5
}
```

## Available Categories

| Slug | Name | Description |
|------|------|-------------|
| ai | AI | AI-powered applications |
| games | Games | Play-to-earn games and entertainment |
| trading | Trading | Trading and market analysis tools |
| depin | DePIN | Decentralized Physical Infrastructure |
| defi | DeFi | Decentralized Finance applications |
| nft | NFT | NFT marketplaces and collections |
| staking | Staking | Staking and yield farming |
| services | Services | Various utility services |

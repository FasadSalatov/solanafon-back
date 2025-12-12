# Sol Mobile Backend

Backend API for Sol Mobile platform with mini-apps.

## Tech Stack

- Go 1.21+
- PostgreSQL
- Fiber (web framework)
- GORM (ORM)
- JWT authentication

## Quick Start

### Prerequisites

- Go 1.21 or higher
- PostgreSQL 15
- Docker (optional)

### Installation

1. Clone the repository:
```bash
git clone <repo-url>
cd solanafon-back
```

2. Copy environment variables:
```bash
cp .env.example .env
```

3. Configure `.env` file:
```env
DATABASE_URL=postgresql://postgres:postgres@localhost:5432/solanafon?sslmode=disable
JWT_SECRET=your-super-secret-jwt-key
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USER=your-email@gmail.com
SMTP_PASSWORD=your-app-password
```

4. Start PostgreSQL (with Docker):
```bash
make docker-up
```

Or install PostgreSQL locally and create database:
```bash
createdb solanafon
```

5. Install dependencies:
```bash
make install
```

6. Seed the database:
```bash
make seed
```

7. Start the server:
```bash
make run
```

Server will be available at: `http://localhost:8080`

## Available Commands

```bash
make help          # Show all commands
make install       # Install dependencies
make run          # Start server
make seed         # Seed database
make dev          # Start with hot-reload (requires air)
make build        # Build application
make test         # Run tests
make docker-up    # Start Docker services
make docker-down  # Stop Docker services
```

## Project Structure

```
solanafon-back/
├── cmd/
│   ├── server/          # Application entry point
│   └── seed/            # Database seeding script
├── internal/
│   ├── config/          # Configuration
│   ├── database/        # Database connection
│   ├── handlers/        # HTTP handlers
│   ├── middleware/      # Middleware
│   ├── models/          # Database models
│   ├── routes/          # Routing
│   └── utils/           # Utilities (JWT, OTP, Email)
├── docs/                # API documentation (GitBook)
├── .env.example         # Environment variables example
├── docker-compose.yml   # Docker configuration
├── Makefile            # Development commands
└── README.md
```

## API Documentation

Full API documentation is available in [docs/](docs/) folder (GitBook format).

### Main Endpoints

**Authentication:**
- `POST /api/v1/auth/email/request` - Request OTP code
- `POST /api/v1/auth/email/verify` - Verify OTP and get token
- `GET /api/v1/auth/me` - Get current user

**Apps:**
- `GET /api/v1/apps` - List all apps
- `GET /api/v1/apps/:id` - App details
- `POST /api/v1/apps` - Create new app
- `GET /api/v1/apps/my` - My apps

**Dev Studio:**
- `POST /api/v1/apps/devstudio/message` - Interactive app management

**Developer API:**
- `GET /api/v1/bot/getMe` - Get app info
- `POST /api/v1/bot/sendMessage` - Send message to user
- `GET /api/v1/bot/getUpdates` - Get pending messages
- `POST /api/v1/bot/setWebhook` - Set webhook URL

**Categories:**
- `GET /api/v1/categories` - List categories
- `GET /api/v1/categories/:slug/apps` - Apps by category

**Profile:**
- `GET /api/v1/profile` - Get profile
- `PATCH /api/v1/profile` - Update profile

**Mana Points:**
- `GET /api/v1/mana` - Get MP balance and history
- `POST /api/v1/mana/topup` - Top up MP

**Secret Login:**
- `GET /api/v1/secret/numbers` - Available virtual numbers
- `POST /api/v1/secret/activate` - Activate number
- `GET /api/v1/secret/status` - Access status

## Usage Examples

### 1. Authentication

```bash
# Request OTP
curl -X POST http://localhost:8080/api/v1/auth/email/request \
  -H "Content-Type: application/json" \
  -d '{"email": "user@example.com"}'

# Verify OTP and get token
curl -X POST http://localhost:8080/api/v1/auth/email/verify \
  -H "Content-Type: application/json" \
  -d '{"email": "user@example.com", "code": "123456"}'
```

### 2. Create an App via Dev Studio

```bash
curl -X POST http://localhost:8080/api/v1/apps/devstudio/message \
  -H "Authorization: Bearer <your-token>" \
  -H "Content-Type: application/json" \
  -d '{"content": "/newapp"}'
```

### 3. Send Message (Developer API)

```bash
curl -X POST http://localhost:8080/api/v1/bot/sendMessage \
  -H "Authorization: Bearer <api-token>" \
  -H "Content-Type: application/json" \
  -d '{"chat_id": 123, "text": "Hello!"}'
```

## Development

### Hot Reload

For development with auto-reload install [Air](https://github.com/cosmtrek/air):

```bash
go install github.com/cosmtrek/air@latest
make dev
```

### Testing

```bash
make test
```

## Production

### Build

```bash
make build
./bin/server
```

### Docker

```bash
docker build -t sol-mobile-backend .
docker run -p 8080:8080 --env-file .env sol-mobile-backend
```

## Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `PORT` | Server port | `8080` |
| `DATABASE_URL` | PostgreSQL connection URL | - |
| `JWT_SECRET` | JWT secret key | - |
| `SMTP_HOST` | SMTP server | `smtp.gmail.com` |
| `SMTP_PORT` | SMTP port | `587` |
| `SMTP_USER` | SMTP user | - |
| `SMTP_PASSWORD` | SMTP password | - |
| `OTP_EXPIRY_MINUTES` | OTP lifetime | `10` |

## License

MIT

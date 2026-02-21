# MCP Server Setup

Solafon MCP Server даёт AI-ассистентам (Claude, Cursor, VS Code и др.) прямой доступ к инструментам разработки на платформе Solafon.

## Что такое MCP?

[Model Context Protocol (MCP)](https://modelcontextprotocol.io/) — открытый протокол, который позволяет AI-ассистентам подключаться к внешним инструментам и данным. Solafon MCP Server превращает вашего AI-помощника в полноценного разработчика на платформе Solafon.

## Установка

```bash
# Из директории проекта
cd mcp
npm install
npm run build
```

## Настройка Claude Desktop

Добавьте в `~/Library/Application Support/Claude/claude_desktop_config.json` (macOS) или `%APPDATA%/Claude/claude_desktop_config.json` (Windows):

```json
{
  "mcpServers": {
    "solafon": {
      "command": "node",
      "args": ["/absolute/path/to/solanafon-back/mcp/dist/index.js"],
      "env": {
        "SOLAFON_BOT_TOKEN": "ваш-api-токен",
        "SOLAFON_API_URL": "https://api.solafon.com/api/v1"
      }
    }
  }
}
```

## Настройка Cursor

Добавьте `.cursor/mcp.json` в ваш проект:

```json
{
  "mcpServers": {
    "solafon": {
      "command": "node",
      "args": ["./mcp/dist/index.js"],
      "env": {
        "SOLAFON_BOT_TOKEN": "ваш-api-токен"
      }
    }
  }
}
```

## Настройка VS Code (Claude Code)

Добавьте `.vscode/mcp.json`:

```json
{
  "servers": {
    "solafon": {
      "command": "node",
      "args": ["./mcp/dist/index.js"],
      "env": {
        "SOLAFON_BOT_TOKEN": "ваш-api-токен"
      }
    }
  }
}
```

## Получение API токена

1. **Через Dev Studio**: отправьте `/token` в чат с Dev Studio
2. **Через API**: `GET /apps/:id/settings` (требуется JWT)
3. **При создании приложения**: `POST /apps` возвращает `apiToken`

## Переменные окружения

| Переменная | Обязательна | По умолчанию | Описание |
|-----------|------------|-------------|----------|
| `SOLAFON_BOT_TOKEN` | Для Bot API | — | API токен вашего бота |
| `SOLAFON_API_URL` | Нет | `https://api.solafon.com/api/v1` | URL API |

## Проверка подключения

После настройки попросите AI-ассистента:
> "Проверь подключение к Solafon API"

Он использует `solafon_health_check` и `solafon_get_bot_info` для проверки.

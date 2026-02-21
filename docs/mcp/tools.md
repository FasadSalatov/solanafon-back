# MCP Tools Reference

Полный список инструментов, доступных AI-ассистенту через Solafon MCP Server.

## Инструменты документации

### solafon_list_docs
Показать все доступные топики документации.

### solafon_read_docs
Прочитать документацию по конкретному топику.

**Параметры:**
| Параметр | Тип | Описание |
|----------|-----|----------|
| `topic` | string | Ключ топика (introduction, quick-start, authentication, bot-api, app-management, dev-studio, webhooks-guide, message-types, mana-points, api-overview, code-examples, project-architecture, secret-login, mcp-setup) |

### solafon_search_docs
Поиск по всей документации.

**Параметры:**
| Параметр | Тип | Описание |
|----------|-----|----------|
| `query` | string | Поисковый запрос |

---

## Инструменты Bot API

### solafon_get_bot_info
Получить информацию о боте (id, title, username, webhook).
Требует `SOLAFON_BOT_TOKEN`.

### solafon_send_message
Отправить сообщение пользователю.

**Параметры:**
| Параметр | Тип | Обязателен | Описание |
|----------|-----|-----------|----------|
| `chat_id` | number | Да | ID пользователя |
| `text` | string | Да | Текст сообщения |
| `message_type` | string | Нет | `text`, `image` или `button` |
| `metadata` | string | Нет | JSON-строка с доп. данными |

**Примеры metadata:**
```json
// Для image:
{"image_url": "https://example.com/pic.png"}

// Для button:
{"buttons": [{"text": "Купить", "action": "buy"}, {"text": "Отмена", "action": "cancel"}]}
```

### solafon_get_updates
Получить непрочитанные сообщения (long-polling). Сообщения помечаются как прочитанные после получения.

### solafon_set_webhook
Установить URL вебхука для получения сообщений в реальном времени.

**Параметры:**
| Параметр | Тип | Описание |
|----------|-----|----------|
| `url` | string | HTTPS URL для получения обновлений |

### solafon_delete_webhook
Удалить текущий вебхук. Переключиться на режим polling.

### solafon_get_webhook_info
Получить информацию о текущем вебхуке (URL, количество ожидающих обновлений).

### solafon_set_commands
Определить команды бота. Команды с полем `response` отвечают автоматически.

**Параметры:**
| Параметр | Тип | Описание |
|----------|-----|----------|
| `commands` | array | Массив объектов команд |

**Объект команды:**
```json
{
  "command": "help",           // без / в начале
  "description": "Показать помощь",
  "response": "Доступные команды: /start, /help"  // опционально
}
```

### solafon_get_commands
Получить все определённые команды бота.

---

## Инструменты разработки

### solafon_api_request
Выполнить произвольный API запрос к любому эндпоинту Solafon.

**Параметры:**
| Параметр | Тип | Обязателен | Описание |
|----------|-----|-----------|----------|
| `method` | string | Да | GET, POST, PUT, PATCH, DELETE |
| `path` | string | Да | Путь API (напр. `/apps`, `/bot/getMe`) |
| `body` | object | Нет | Тело запроса (для POST/PUT/PATCH) |
| `token` | string | Нет | Переопределить токен |

### solafon_health_check
Проверить доступность Solafon API.

### solafon_scaffold_bot
Сгенерировать шаблон кода бота.

**Параметры:**
| Параметр | Тип | Описание |
|----------|-----|----------|
| `template` | string | Ключ шаблона (см. solafon_list_templates) |

**Доступные шаблоны:**
- `node-polling` — Node.js бот с polling
- `node-webhook` — Node.js бот с Express webhook
- `python-polling` — Python бот с polling
- `python-webhook` — Python бот с Flask webhook
- `go-webhook` — Go бот с net/http webhook
- `node-ai-bot` — AI бот с OpenAI (Node.js)

### solafon_list_templates
Показать все доступные шаблоны кода.

---

## Ресурсы (Resources)

MCP-ресурсы — это данные, которые AI может прочитать по URI:

| URI | Описание |
|-----|----------|
| `solafon://docs/introduction` | Введение в платформу |
| `solafon://docs/quick-start` | Быстрый старт |
| `solafon://docs/authentication` | Аутентификация |
| `solafon://docs/bot-api` | Bot API справочник |
| `solafon://docs/app-management` | Управление приложениями |
| `solafon://docs/dev-studio` | Dev Studio |
| `solafon://docs/webhooks-guide` | Руководство по вебхукам |
| `solafon://docs/message-types` | Типы сообщений |
| `solafon://docs/mana-points` | Mana Points |
| `solafon://docs/api-overview` | Обзор API |
| `solafon://docs/code-examples` | Примеры кода |
| `solafon://docs/project-architecture` | Архитектура проекта |
| `solafon://docs/secret-login` | Secret Login |
| `solafon://docs/mcp-setup` | Настройка MCP |
| `solafon://docs/full` | Полная документация |
| `solafon://templates/*` | Шаблоны кода |

---

## Промпты (Prompts)

Предустановленные сценарии для AI-ассистента:

### create-solafon-bot
Пошаговое создание нового бота на платформе Solafon.

**Параметры:**
- `language`: javascript, python, go
- `architecture`: polling, webhook

### debug-solafon-bot
Диагностика проблем с ботом.

**Параметры:**
- `issue`: описание проблемы

### solafon-api-explorer
Интерактивное исследование API Solafon.

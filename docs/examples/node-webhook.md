# Node.js Webhook Bot

Express.js бот с webhook для получения сообщений в реальном времени.

## Требования

- Node.js 18+
- API токен Solafon
- HTTPS-сервер (для production)

## Установка

```bash
mkdir my-solafon-bot && cd my-solafon-bot
npm init -y
npm install express
```

## Код

```javascript
// bot.js
import express from 'express';

const app = express();
app.use(express.json());

const TOKEN = process.env.SOLAFON_BOT_TOKEN;
const API = 'https://api.solafon.com/api/v1';

async function sendMessage(chatId, text, options = {}) {
  await fetch(`${API}/bot/sendMessage`, {
    method: 'POST',
    headers: {
      'Authorization': `Bearer ${TOKEN}`,
      'Content-Type': 'application/json'
    },
    body: JSON.stringify({ chat_id: chatId, text, ...options })
  });
}

app.post('/webhook', async (req, res) => {
  try {
    const { message } = req.body;
    const chatId = message.from.id;
    const text = message.text || '';

    if (text.startsWith('/')) {
      const command = text.split(' ')[0].slice(1);
      const args = text.split(' ').slice(1).join(' ');

      switch (command) {
        case 'start':
          await sendMessage(chatId, 'Welcome to my bot!');
          break;
        case 'help':
          await sendMessage(chatId, 'Available commands: /start, /help');
          break;
        default:
          await sendMessage(chatId, `Unknown command: /${command}`);
      }
    } else {
      await sendMessage(chatId, `Received: ${text}`);
    }

    res.json({ ok: true });
  } catch (error) {
    console.error('Webhook error:', error);
    res.status(500).json({ ok: false });
  }
});

const PORT = process.env.PORT || 3000;
app.listen(PORT, () => {
  console.log(`Webhook server running on port ${PORT}`);
});
```

## Запуск

```bash
SOLAFON_BOT_TOKEN=your-token node bot.js
```

## Регистрация webhook

```bash
curl -X POST https://api.solafon.com/api/v1/bot/setWebhook \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"url": "https://your-domain.com/webhook"}'
```

## Формат payload

Solafon отправляет POST запрос на ваш webhook:

```json
{
  "update_id": 12345,
  "message": {
    "message_id": 67,
    "from": {"id": 1, "email": "user@mail.com", "name": "User", "language": "en"},
    "chat": {"id": 1, "type": "private"},
    "date": 1704067200,
    "text": "Hello!"
  }
}
```

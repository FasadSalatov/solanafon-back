# Node.js Polling Bot

Простой бот на Node.js, использующий polling для получения сообщений.

## Требования

- Node.js 18+
- API токен Solafon

## Установка

```bash
mkdir my-solafon-bot && cd my-solafon-bot
npm init -y
```

## Код

```javascript
// bot.js
const TOKEN = process.env.SOLAFON_BOT_TOKEN;
const API = 'https://api.solafon.com/api/v1';

async function sendMessage(chatId, text, options = {}) {
  const body = { chat_id: chatId, text, ...options };
  const res = await fetch(`${API}/bot/sendMessage`, {
    method: 'POST',
    headers: {
      'Authorization': `Bearer ${TOKEN}`,
      'Content-Type': 'application/json'
    },
    body: JSON.stringify(body)
  });
  return res.json();
}

async function handleMessage(msg) {
  const chatId = msg.from.id;
  const text = msg.text || '';

  if (text === '/start') {
    await sendMessage(chatId, 'Welcome! I am your bot. Send me anything!');
  } else if (text === '/help') {
    await sendMessage(chatId, 'Commands:\n/start - Start\n/help - Show help');
  } else {
    await sendMessage(chatId, `You said: ${text}`);
  }
}

async function poll() {
  console.log('Bot started polling...');
  while (true) {
    try {
      const res = await fetch(`${API}/bot/getUpdates`, {
        headers: { 'Authorization': `Bearer ${TOKEN}` }
      });
      const data = await res.json();
      if (data.ok && data.result?.length > 0) {
        for (const msg of data.result) {
          await handleMessage(msg);
        }
      }
    } catch (error) {
      console.error('Polling error:', error.message);
    }
    await new Promise(resolve => setTimeout(resolve, 2000));
  }
}

poll();
```

## Запуск

```bash
SOLAFON_BOT_TOKEN=your-token node bot.js
```

## Как это работает

1. Бот каждые 2 секунды запрашивает новые сообщения через `GET /bot/getUpdates`
2. Полученные сообщения обрабатываются в `handleMessage`
3. Ответы отправляются через `POST /bot/sendMessage`
4. Сообщения помечаются как прочитанные автоматически после `getUpdates`

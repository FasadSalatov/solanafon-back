# AI Bot with OpenAI

AI-бот с памятью разговоров, использующий OpenAI GPT.

## Требования

- Node.js 18+
- API токен Solafon
- API ключ OpenAI

## Установка

```bash
npm install openai
```

## Код

```javascript
// bot.js
import OpenAI from 'openai';

const TOKEN = process.env.SOLAFON_BOT_TOKEN;
const API = 'https://api.solafon.com/api/v1';
const openai = new OpenAI();

const conversations = new Map();

async function sendMessage(chatId, text) {
  await fetch(`${API}/bot/sendMessage`, {
    method: 'POST',
    headers: {
      'Authorization': `Bearer ${TOKEN}`,
      'Content-Type': 'application/json'
    },
    body: JSON.stringify({ chat_id: chatId, text })
  });
}

async function handleMessage(msg) {
  const chatId = msg.from.id;
  const text = msg.text || '';

  if (text === '/start') {
    conversations.delete(chatId);
    await sendMessage(chatId, 'Hello! I am an AI assistant. Ask me anything!');
    return;
  }

  if (!conversations.has(chatId)) {
    conversations.set(chatId, [{
      role: 'system',
      content: 'You are a helpful assistant. Be concise and friendly.'
    }]);
  }

  const history = conversations.get(chatId);
  history.push({ role: 'user', content: text });

  // Keep last 20 messages
  if (history.length > 21) {
    history.splice(1, history.length - 21);
  }

  try {
    const completion = await openai.chat.completions.create({
      model: 'gpt-4o-mini',
      messages: history,
      max_tokens: 500
    });

    const reply = completion.choices[0].message.content;
    history.push({ role: 'assistant', content: reply });
    await sendMessage(chatId, reply);
  } catch (error) {
    console.error('OpenAI error:', error);
    await sendMessage(chatId, 'Sorry, I encountered an error. Please try again.');
  }
}

async function poll() {
  console.log('AI Bot started...');
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
      console.error('Poll error:', error.message);
    }
    await new Promise(r => setTimeout(r, 2000));
  }
}

poll();
```

## Запуск

```bash
SOLAFON_BOT_TOKEN=your-token OPENAI_API_KEY=your-key node bot.js
```

## Кастомизация

### Изменение системного промпта

```javascript
conversations.set(chatId, [{
  role: 'system',
  content: 'You are a Solana trading assistant. You help users understand DeFi, token prices, and trading strategies. Always provide data-driven advice.'
}]);
```

### Использование Claude вместо GPT

```javascript
import Anthropic from '@anthropic-ai/sdk';
const anthropic = new Anthropic();

const response = await anthropic.messages.create({
  model: 'claude-sonnet-4-20250514',
  max_tokens: 500,
  messages: history.filter(m => m.role !== 'system'),
  system: history[0].content
});
```

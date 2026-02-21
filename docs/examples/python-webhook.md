# Python Webhook Bot

Flask бот с webhook.

## Требования

- Python 3.8+
- API токен Solafon

## Установка

```bash
pip install flask requests
```

## Код

```python
# bot.py
import os
import requests
from flask import Flask, request, jsonify

app = Flask(__name__)

TOKEN = os.environ["SOLAFON_BOT_TOKEN"]
API = "https://api.solafon.com/api/v1"
HEADERS = {"Authorization": f"Bearer {TOKEN}", "Content-Type": "application/json"}


def send_message(chat_id: int, text: str, **kwargs):
    payload = {"chat_id": chat_id, "text": text, **kwargs}
    requests.post(f"{API}/bot/sendMessage", headers=HEADERS, json=payload)


@app.route("/webhook", methods=["POST"])
def webhook():
    data = request.json
    message = data.get("message", {})
    chat_id = message["from"]["id"]
    text = message.get("text", "")

    if text.startswith("/"):
        command = text.split()[0][1:]
        if command == "start":
            send_message(chat_id, "Welcome to my bot!")
        elif command == "help":
            send_message(chat_id, "Commands: /start, /help")
        else:
            send_message(chat_id, f"Unknown command: /{command}")
    else:
        send_message(chat_id, f"Received: {text}")

    return jsonify({"ok": True})


if __name__ == "__main__":
    app.run(host="0.0.0.0", port=int(os.environ.get("PORT", 3000)))
```

## Запуск

```bash
SOLAFON_BOT_TOKEN=your-token python bot.py
```

## Регистрация webhook

```bash
curl -X POST https://api.solafon.com/api/v1/bot/setWebhook \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -d '{"url": "https://your-server.com/webhook"}'
```

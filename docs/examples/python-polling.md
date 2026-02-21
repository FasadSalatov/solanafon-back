# Python Polling Bot

Простой бот на Python с polling.

## Требования

- Python 3.8+
- API токен Solafon

## Установка

```bash
pip install requests
```

## Код

```python
# bot.py
import os
import time
import requests

TOKEN = os.environ["SOLAFON_BOT_TOKEN"]
API = "https://api.solafon.com/api/v1"
HEADERS = {"Authorization": f"Bearer {TOKEN}", "Content-Type": "application/json"}


def send_message(chat_id: int, text: str, **kwargs):
    payload = {"chat_id": chat_id, "text": text, **kwargs}
    return requests.post(f"{API}/bot/sendMessage", headers=HEADERS, json=payload).json()


def handle_message(msg: dict):
    chat_id = msg["from"]["id"]
    text = msg.get("text", "")

    if text == "/start":
        send_message(chat_id, "Welcome! I am your Solafon bot.")
    elif text == "/help":
        send_message(chat_id, "Commands:\n/start - Start\n/help - Help")
    else:
        send_message(chat_id, f"You said: {text}")


def poll():
    print("Bot started polling...")
    while True:
        try:
            res = requests.get(f"{API}/bot/getUpdates", headers=HEADERS)
            data = res.json()
            if data.get("ok") and data.get("result"):
                for msg in data["result"]:
                    handle_message(msg)
        except Exception as e:
            print(f"Error: {e}")
        time.sleep(2)


if __name__ == "__main__":
    poll()
```

## Запуск

```bash
SOLAFON_BOT_TOKEN=your-token python bot.py
```

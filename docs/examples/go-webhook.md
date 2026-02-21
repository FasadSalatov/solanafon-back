# Go Webhook Bot

Бот на Go с net/http webhook.

## Требования

- Go 1.21+
- API токен Solafon

## Код

```go
// main.go
package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
)

var (
	botToken = os.Getenv("SOLAFON_BOT_TOKEN")
	apiURL   = "https://api.solafon.com/api/v1"
)

type Update struct {
	UpdateID int `json:"update_id"`
	Message  struct {
		MessageID int `json:"message_id"`
		From      struct {
			ID       int    `json:"id"`
			Email    string `json:"email"`
			Name     string `json:"name"`
			Language string `json:"language"`
		} `json:"from"`
		Chat struct {
			ID   int    `json:"id"`
			Type string `json:"type"`
		} `json:"chat"`
		Date int    `json:"date"`
		Text string `json:"text"`
	} `json:"message"`
}

func sendMessage(chatID int, text string) error {
	body, _ := json.Marshal(map[string]interface{}{
		"chat_id": chatID,
		"text":    text,
	})
	req, _ := http.NewRequest("POST", apiURL+"/bot/sendMessage", bytes.NewBuffer(body))
	req.Header.Set("Authorization", "Bearer "+botToken)
	req.Header.Set("Content-Type", "application/json")
	_, err := http.DefaultClient.Do(req)
	return err
}

func webhookHandler(w http.ResponseWriter, r *http.Request) {
	var update Update
	if err := json.NewDecoder(r.Body).Decode(&update); err != nil {
		http.Error(w, "Bad request", 400)
		return
	}

	chatID := update.Message.From.ID
	text := update.Message.Text

	if strings.HasPrefix(text, "/") {
		cmd := strings.Split(text, " ")[0][1:]
		switch cmd {
		case "start":
			sendMessage(chatID, "Welcome to my Solafon bot!")
		case "help":
			sendMessage(chatID, "Commands: /start, /help")
		default:
			sendMessage(chatID, fmt.Sprintf("Unknown command: /%s", cmd))
		}
	} else {
		sendMessage(chatID, fmt.Sprintf("You said: %s", text))
	}

	json.NewEncoder(w).Encode(map[string]bool{"ok": true})
}

func main() {
	http.HandleFunc("/webhook", webhookHandler)
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}
	log.Printf("Webhook server running on :%s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
```

## Запуск

```bash
SOLAFON_BOT_TOKEN=your-token go run main.go
```

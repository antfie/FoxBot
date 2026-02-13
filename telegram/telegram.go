package telegram

import (
	"fmt"
	"github.com/antfie/FoxBot/db"
	"github.com/antfie/FoxBot/types"
	"github.com/antfie/FoxBot/utils"
	"log"
	"net/url"
	"strings"
	"time"
)

type Telegram struct {
	url      string
	chatID   string
	db       *db.DB
	duration *types.TimeDuration
}

var telegramHeaders = map[string]string{
	"Content-Type": "application/x-www-form-urlencoded",
}

func NewTelegram(config *types.Telegram, db *db.DB) *Telegram {
	t := &Telegram{
		url:      fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", config.Token),
		chatID:   config.ChatID,
		db:       db,
		duration: config.Duration,
	}

	go t.processor()

	return t
}

func (t *Telegram) processor() {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if t.duration != nil && !utils.IsWithinDuration(time.Now(), *t.duration) {
				continue
			}

			messages := t.db.ConsumeTelegramNotificationQueue()

			if len(messages) > 0 {
				message := strings.Join(messages, "\n")
				t.notify(message)
			}
		}
	}
}

func (t *Telegram) notify(message string) {
	form := url.Values{}
	form.Add("chat_id", t.chatID)
	form.Add("text", message)

	response := utils.HttpRequest("POST", t.url, telegramHeaders, strings.NewReader(form.Encode()))

	if response == nil {
		log.Print("Could not connect to Telegram API")
		return
	}

	response.Body.Close()
}

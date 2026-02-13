package integrations

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/antfie/FoxBot/bayes"
	"github.com/antfie/FoxBot/db"
	"github.com/antfie/FoxBot/types"
	"github.com/antfie/FoxBot/utils"
)

type Telegram struct {
	apiBase  string
	chatID   string
	db       *db.DB
	duration *types.TimeDuration
	bayes    *bayes.Classifier
}

var telegramHeaders = map[string]string{
	"Content-Type": "application/x-www-form-urlencoded",
}

type telegramResponse struct {
	OK     bool             `json:"ok"`
	Result []telegramUpdate `json:"result"`
}

type telegramUpdate struct {
	UpdateID      int                    `json:"update_id"`
	CallbackQuery *telegramCallbackQuery `json:"callback_query"`
}

type telegramCallbackQuery struct {
	ID   string `json:"id"`
	Data string `json:"data"`
}

func NewTelegram(config *types.Telegram, db *db.DB, classifier *bayes.Classifier) *Telegram {
	t := &Telegram{
		apiBase:  fmt.Sprintf("https://api.telegram.org/bot%s", config.Token),
		chatID:   config.ChatID,
		db:       db,
		duration: config.Duration,
		bayes:    classifier,
	}

	go t.processor()
	go t.feedbackProcessor()

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

	response := utils.HttpRequest("POST", t.apiBase+"/sendMessage", telegramHeaders, strings.NewReader(form.Encode()))

	if response == nil {
		log.Print("Could not connect to Telegram API")
		return
	}

	if err := response.Body.Close(); err != nil {
		log.Print(err)
	}
}

func (t *Telegram) SendWithFeedback(message, articleHash string) {
	if t.duration != nil && !utils.IsWithinDuration(time.Now(), *t.duration) {
		return
	}

	replyMarkup := fmt.Sprintf(`{"inline_keyboard":[[{"text":"👍","callback_data":"r:%s"},{"text":"👎","callback_data":"i:%s"}]]}`, articleHash, articleHash)

	form := url.Values{}
	form.Add("chat_id", t.chatID)
	form.Add("text", message)
	form.Add("reply_markup", replyMarkup)

	response := utils.HttpRequest("POST", t.apiBase+"/sendMessage", telegramHeaders, strings.NewReader(form.Encode()))

	if response == nil {
		log.Print("Could not connect to Telegram API for feedback message")
		return
	}

	if err := response.Body.Close(); err != nil {
		log.Print(err)
	}
}

func (t *Telegram) feedbackProcessor() {
	ticker := time.NewTicker(30 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			t.pollFeedback()
		}
	}
}

func (t *Telegram) pollFeedback() {
	offsetStr := t.db.GetTelegramState("update_offset")

	offset := 0
	if len(offsetStr) > 0 {
		var err error
		offset, err = strconv.Atoi(offsetStr)

		if err != nil {
			log.Printf("Invalid telegram update offset: %s", offsetStr)
			offset = 0
		}
	}

	requestURL := fmt.Sprintf("%s/getUpdates?timeout=0", t.apiBase)
	if offset > 0 {
		requestURL = fmt.Sprintf("%s&offset=%d", requestURL, offset)
	}

	response := utils.HttpRequest("GET", requestURL, nil, nil)

	if response == nil {
		return
	}

	defer func() {
		if err := response.Body.Close(); err != nil {
			log.Print(err)
		}
	}()

	body, err := io.ReadAll(response.Body)

	if err != nil {
		log.Printf("Could not read Telegram getUpdates response: %v", err)
		return
	}

	var result telegramResponse
	if err = json.Unmarshal(body, &result); err != nil {
		log.Printf("Could not parse Telegram getUpdates response: %v", err)
		return
	}

	if !result.OK {
		return
	}

	maxUpdateID := offset - 1

	for _, update := range result.Result {
		if update.UpdateID > maxUpdateID {
			maxUpdateID = update.UpdateID
		}

		if update.CallbackQuery == nil {
			continue
		}

		t.processCallback(update.CallbackQuery)
	}

	if maxUpdateID >= offset {
		t.db.SetTelegramState("update_offset", strconv.Itoa(maxUpdateID+1))
	}
}

func (t *Telegram) processCallback(query *telegramCallbackQuery) {
	data := query.Data

	if len(data) < 3 || data[1] != ':' {
		return
	}

	prefix := data[0]
	hash := data[2:]

	feedGroup, title, found := t.db.BayesGetArticle(hash)

	if !found {
		log.Printf("Bayes article not found for hash: %s", hash)
		t.answerCallback(query.ID)
		return
	}

	relevant := prefix == 'r'
	t.bayes.Train(feedGroup, title, relevant)

	label := "irrelevant"
	if relevant {
		label = "relevant"
	}
	log.Printf("Bayes trained [%s] as %s: %s", feedGroup, label, title)

	t.answerCallback(query.ID)
}

func (t *Telegram) answerCallback(callbackQueryID string) {
	form := url.Values{}
	form.Add("callback_query_id", callbackQueryID)

	response := utils.HttpRequest("POST", t.apiBase+"/answerCallbackQuery", telegramHeaders, strings.NewReader(form.Encode()))

	if response == nil {
		return
	}

	if err := response.Body.Close(); err != nil {
		log.Print(err)
	}
}

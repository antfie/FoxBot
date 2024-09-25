package slack

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

type Slack struct {
	url string
	db  *db.DB
}

var slackHeaders = map[string]string{
	"Content-Type": "application/x-www-form-urlencoded",
}

func NewSlack(config *types.Slack, db *db.DB) *Slack {
	formattedUrl := fmt.Sprintf("https://slack.com/api/chat.postMessage?token=%s&channel=%s", config.Token, config.ChannelId)

	slack := &Slack{
		url: formattedUrl,
		db:  db,
	}

	go slack.processor()

	return slack
}

func (s *Slack) processor() {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			messages := s.db.ConsumeSlackNotificationQueue()
			
			if len(messages) > 0 {
				message := strings.Join(messages, "\n")
				s.notify(message)
			}
		}
	}
}

// TODO: Rich text messages - https://api.slack.com/reference/block-kit/blocks#rich_text

func (s *Slack) notify(message string) {
	form := url.Values{}
	form.Add("text", message)

	response := utils.HttpRequest("POST", s.url, slackHeaders, strings.NewReader(form.Encode()))

	if response == nil {
		log.Panic("Could not connect to Slack API")
	}
}

package integrations

import (
	"log"
	"strings"
	"time"

	"github.com/antfie/FoxBot/db"
	"github.com/antfie/FoxBot/types"
	"github.com/antfie/FoxBot/utils"
)

type Discord struct {
	webhookURL string
	db         *db.DB
	duration   *types.TimeDuration
}

var discordHeaders = map[string]string{
	"Content-Type": "application/json",
}

func NewDiscord(config *types.Discord, db *db.DB) *Discord {
	d := &Discord{
		webhookURL: config.WebhookURL,
		db:         db,
		duration:   config.Duration,
	}

	go d.processor()

	return d
}

func (d *Discord) processor() {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if d.duration != nil && !utils.IsWithinDuration(time.Now(), *d.duration) {
				continue
			}

			messages := d.db.ConsumeDiscordNotificationQueue()

			if len(messages) > 0 {
				message := strings.Join(messages, "\n")
				d.notify(message)
			}
		}
	}
}

func (d *Discord) notify(message string) {
	body := strings.NewReader(`{"content":` + jsonEscapeString(message) + `}`)

	response := utils.HttpRequest("POST", d.webhookURL, discordHeaders, body)

	if response == nil {
		log.Print("Could not connect to Discord webhook")
		return
	}

	if err := response.Body.Close(); err != nil {
		log.Print(err)
	}
}

func jsonEscapeString(s string) string {
	// Use strings.Builder to build a JSON-escaped string
	var b strings.Builder
	b.WriteByte('"')

	for _, c := range s {
		switch c {
		case '"':
			b.WriteString(`\"`)
		case '\\':
			b.WriteString(`\\`)
		case '\n':
			b.WriteString(`\n`)
		case '\r':
			b.WriteString(`\r`)
		case '\t':
			b.WriteString(`\t`)
		default:
			b.WriteRune(c)
		}
	}

	b.WriteByte('"')
	return b.String()
}

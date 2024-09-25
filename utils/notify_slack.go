package utils

import (
	"fmt"
	"github.com/antfie/FoxBot/types"
	"log"
	"net/url"
	"strings"
)

// TODO: Rich text messages - https://api.slack.com/reference/block-kit/blocks#rich_text

var slackHeaders = map[string]string{
	"Content-Type": "application/x-www-form-urlencoded",
}

func NotifySlack(creds *types.Slack, message string) {
	sendMessage(creds, message)
}

func sendMessage(creds *types.Slack, message string) {
	form := url.Values{}
	form.Add("text", message)

	formattedUrl := fmt.Sprintf("https://slack.com/api/chat.postMessage?token=%s&channel=%s", creds.Token, creds.ChannelId)

	response := HttpRequest("POST", formattedUrl, slackHeaders, strings.NewReader(form.Encode()))

	if response == nil {
		log.Panic("Could not connect to Slack webhook")
	}
}

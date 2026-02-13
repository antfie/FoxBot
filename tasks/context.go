package tasks

import (
	"github.com/antfie/FoxBot/db"
	"github.com/antfie/FoxBot/integrations/slack"
	"github.com/antfie/FoxBot/integrations/telegram"
	"github.com/antfie/FoxBot/types"
)

type Context struct {
	Config   *types.Config
	DB       *db.DB
	Slack    *slack.Slack
	Telegram *telegram.Telegram
}

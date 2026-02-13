package tasks

import (
	"github.com/antfie/FoxBot/bayes"
	"github.com/antfie/FoxBot/db"
	"github.com/antfie/FoxBot/integrations"
	"github.com/antfie/FoxBot/types"
)

type Context struct {
	Config   *types.Config
	DB       *db.DB
	Slack    *integrations.Slack
	Telegram *integrations.Telegram
	Bayes    *bayes.Classifier
}

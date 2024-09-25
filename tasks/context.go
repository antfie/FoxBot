package tasks

import (
	"github.com/antfie/FoxBot/db"
	"github.com/antfie/FoxBot/types"
)

type Context struct {
	Config *types.Config
	DB     *db.DB
}

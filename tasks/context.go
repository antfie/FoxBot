package tasks

import (
	"foxbot/db"
	"foxbot/types"
)

type Context struct {
	Config *types.Config
	DB     *db.DB
}

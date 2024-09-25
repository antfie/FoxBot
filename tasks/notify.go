package tasks

import "github.com/antfie/FoxBot/utils"

func (c *Context) Notify(message string) {
	if c.Config.Output.Console {
		utils.NotifyConsole(message)
	}

	if c.Config.Output.Slack != nil {
		c.DB.QueueSlackNotification(message)
	}
}

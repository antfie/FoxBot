package tasks

import "foxbot/utils"

func (c *Context) Notify(message string) {
	if c.Config.Output.Console {
		utils.NotifyConsole(message)
	}

	if c.Config.Output.Slack != nil {
		utils.NotifySlack(c.Config.Output.Slack, message)
	}
}

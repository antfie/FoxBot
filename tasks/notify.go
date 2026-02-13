package tasks

import "github.com/antfie/FoxBot/utils"

func (c *Context) Notify(message string) {
	if c.Config.Output.Console {
		utils.NotifyConsole(message)
	}

	if c.Config.Output.Slack != nil {
		c.DB.QueueSlackNotification(message)
	}

	if c.Config.Output.Telegram != nil {
		c.DB.QueueTelegramNotification(message)
	}
}

func (c *Context) NotifyGood(message string) {
	if c.Config.Output.Console {
		utils.NotifyConsoleGood(message)
	}

	if c.Config.Output.Slack != nil {
		c.DB.QueueSlackNotification(message)
	}

	if c.Config.Output.Telegram != nil {
		c.DB.QueueTelegramNotification(message)
	}
}

func (c *Context) NotifyBad(message string) {
	if c.Config.Output.Console {
		utils.NotifyConsoleBad(message)
	}

	if c.Config.Output.Slack != nil {
		c.DB.QueueSlackNotification(message)
	}

	if c.Config.Output.Telegram != nil {
		c.DB.QueueTelegramNotification(message)
	}
}

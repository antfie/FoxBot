package tasks

import (
	"fmt"
	"github.com/antfie/FoxBot/utils"
	"time"
)

const firstRunReminderIndex = -1

var currentReminderIndex = firstRunReminderIndex

func (c *Context) Reminders() {
	if c.Config.Reminders.Check.Duration != nil && !utils.IsWithinDuration(time.Now(), *c.Config.Reminders.Check.Duration) {
		return
	}

	// Shuffle the list for the first run
	if currentReminderIndex == firstRunReminderIndex {
		utils.ShuffleStringArray(c.Config.Reminders.Reminders)
	}

	currentReminderIndex++

	if currentReminderIndex > len(c.Config.Reminders.Reminders)-1 {
		currentReminderIndex = 0
	}

	c.Notify(fmt.Sprintf("🧘 %s", c.Config.Reminders.Reminders[currentReminderIndex]))
}

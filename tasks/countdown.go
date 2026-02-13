package tasks

import (
	"fmt"
	"github.com/antfie/FoxBot/utils"
	"sort"
	"sync"
	"time"
)

func (c *Context) Countdown() {
	c.countdown(time.Now())
}

var countdownOnce sync.Once

func (c *Context) countdown(now time.Time) {
	if c.Config.Countdown.Check.Duration != nil && !utils.IsWithinDuration(time.Now(), *c.Config.Countdown.Check.Duration) {
		return
	}

	timers := c.Config.Countdown.Timers

	countdownOnce.Do(func() {
		// Sort the timers by order of due date
		sort.Slice(timers, func(i, j int) bool {
			return timers[i].Date.Before(timers[j].Date)
		})
	})

	for i, x := range timers {
		formattedValue := utils.FormatHumanReadableDuration(now, x.Date)

		if x.LastFormattedValue != formattedValue {
			timers[i].LastFormattedValue = formattedValue
			c.Notify(fmt.Sprintf("⏲️ %s: %s", x.Name, formattedValue))
		}
	}
}

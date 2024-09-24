package types

import "time"

type Countdown struct {
	Check  TimeFrequencyAndDuration
	Timers []CountdownTimer
}

type CountdownTimer struct {
	Name               string
	Date               time.Time
	LastFormattedValue string
}

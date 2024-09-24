package types

import "time"

type TimeFrequencyAndDuration struct {
	Frequency time.Duration
	Duration  *TimeDuration
}

type TimeDuration struct {
	From time.Time
	To   time.Time
}

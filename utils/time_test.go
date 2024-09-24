package utils

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestHumanReadableDuration(t *testing.T) {
	start := time.Date(1784, 5, 15, 8, 12, 54, 0, time.UTC)

	finish := time.Date(2024, 11, 8, 4, 22, 17, 0, time.UTC)
	result := FormatHumanReadableDuration(start, finish)
	assert.Equal(t, "240 years, 5 months, 3 weeks", result)

	finish = start
	result = FormatHumanReadableDuration(start, finish)
	assert.Equal(t, "now", result)

	finish = time.Date(start.Year(), start.Month(), start.Day()+1, start.Hour()+2, start.Minute(), start.Second()+8, 0, time.UTC)
	result = FormatHumanReadableDuration(start, finish)
	assert.Equal(t, "1 day, 2 hours, 8 seconds", result)

	finish = time.Date(start.Year(), start.Month(), start.Day(), start.Hour()+2, start.Minute(), start.Second(), 0, time.UTC)
	result = FormatHumanReadableDuration(start, finish)
	assert.Equal(t, "2 hours", result)

	finish = time.Date(start.Year(), start.Month(), start.Day(), start.Hour()+2, start.Minute()+10, start.Second(), 0, time.UTC)
	result = FormatHumanReadableDuration(start, finish)
	assert.Equal(t, "2 hours, 10 minutes", result)

	finish = time.Date(start.Year(), start.Month(), start.Day(), start.Hour()+1, start.Minute()+1, start.Second()+23, 0, time.UTC)
	result = FormatHumanReadableDuration(start, finish)
	assert.Equal(t, "1 hour, 1 minute, 23 seconds", result)

	finish = time.Date(start.Year(), start.Month(), start.Day(), start.Hour()-1, start.Minute(), start.Second()-1, 0, time.UTC)
	result = FormatHumanReadableDuration(start, finish)
	assert.Equal(t, "1 hour, 1 second ago", result)

	finish = time.Date(start.Year()-10, start.Month(), start.Day(), start.Hour(), start.Minute(), start.Second(), 0, time.UTC)
	result = FormatHumanReadableDuration(start, finish)
	assert.Equal(t, "10 years ago", result)
}

package utils

import (
	"github.com/antfie/FoxBot/types"
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

func TestParseTimeFromString(t *testing.T) {
	result := ParseTimeFromString("08:30")
	assert.Equal(t, 8, result.Hour())
	assert.Equal(t, 30, result.Minute())

	result = ParseTimeFromString("00:00")
	assert.Equal(t, 0, result.Hour())
	assert.Equal(t, 0, result.Minute())

	result = ParseTimeFromString("23:59")
	assert.Equal(t, 23, result.Hour())
	assert.Equal(t, 59, result.Minute())
}

func TestParseDateFromString(t *testing.T) {
	result := ParseDateFromString("25/12/2024")
	assert.Equal(t, 25, result.Day())
	assert.Equal(t, time.December, result.Month())
	assert.Equal(t, 2024, result.Year())
}

func TestParseDurationFromString(t *testing.T) {
	assert.Equal(t, time.Hour, ParseDurationFromString("hourly"))
	assert.Equal(t, time.Hour, ParseDurationFromString("Hourly"))
	assert.Equal(t, 30*time.Minute, ParseDurationFromString("half_hourly"))
	assert.Equal(t, 30*time.Minute, ParseDurationFromString("Half_Hourly"))
	assert.Equal(t, 24*time.Hour, ParseDurationFromString("daily"))
	assert.Equal(t, 24*time.Hour, ParseDurationFromString("Daily"))
}

func TestIsWithinDuration(t *testing.T) {
	duration := types.TimeDuration{
		From: ParseTimeFromString("08:00"),
		To:   ParseTimeFromString("17:00"),
	}

	// Within range
	now := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)
	assert.True(t, IsWithinDuration(now, duration))

	// At start boundary
	now = time.Date(2024, 1, 1, 8, 0, 0, 0, time.UTC)
	assert.True(t, IsWithinDuration(now, duration))

	// At end boundary
	now = time.Date(2024, 1, 1, 17, 0, 0, 0, time.UTC)
	assert.True(t, IsWithinDuration(now, duration))

	// Before range
	now = time.Date(2024, 1, 1, 7, 59, 0, 0, time.UTC)
	assert.False(t, IsWithinDuration(now, duration))

	// After range
	now = time.Date(2024, 1, 1, 17, 1, 0, 0, time.UTC)
	assert.False(t, IsWithinDuration(now, duration))

	// Well before
	now = time.Date(2024, 1, 1, 3, 0, 0, 0, time.UTC)
	assert.False(t, IsWithinDuration(now, duration))

	// Well after
	now = time.Date(2024, 1, 1, 22, 0, 0, 0, time.UTC)
	assert.False(t, IsWithinDuration(now, duration))
}

func TestShuffleStringArray(t *testing.T) {
	original := []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j"}
	shuffled := make([]string, len(original))
	copy(shuffled, original)

	ShuffleStringArray(shuffled)

	// Same length
	assert.Equal(t, len(original), len(shuffled))

	// Same elements
	assert.ElementsMatch(t, original, shuffled)
}

func TestShuffleStringArrayEmpty(t *testing.T) {
	empty := []string{}
	ShuffleStringArray(empty)
	assert.Empty(t, empty)
}

func TestShuffleStringArraySingle(t *testing.T) {
	single := []string{"only"}
	ShuffleStringArray(single)
	assert.Equal(t, []string{"only"}, single)
}

package utils

import (
	"fmt"
	"github.com/antfie/FoxBot/types"
	"log"
	"strings"
	"time"
)

func ParseTimeFromString(t string) time.Time {
	value, err := time.Parse("15:04", t)

	if err != nil {
		log.Panic(err)
	}

	return value
}

func ParseDateFromString(d string) time.Time {
	value, err := time.Parse("02/01/2006", d)

	if err != nil {
		log.Panic(err)
	}

	return value
}

func ParseRSSTimestampFromString(d string) time.Time {
	value, err := time.Parse(time.RFC822, d)

	if err == nil {
		return value
	}

	value, err = time.Parse(time.RFC822Z, d)

	if err == nil {
		return value
	}

	value, err = time.Parse(time.RFC1123, d)

	if err == nil {
		return value
	}

	value, err = time.Parse(time.RFC1123Z, d)

	if err == nil {
		return value
	}

	log.Printf("RSS: Could not parse date: \"%s\"", d)
	return value
}

func ParseDurationFromString(d string) time.Duration {
	if strings.ToLower(d) == "hourly" {
		return time.Hour
	}

	if strings.ToLower(d) == "half_hourly" {
		return time.Minute * 30
	}

	if strings.ToLower(d) == "daily" {
		return time.Hour * 24
	}

	log.Panic("Invalid duration")
	return time.Hour
}

func FormatHumanReadableDuration(start, end time.Time) string {
	years, months, weeks, days, hours, minutes, seconds := dateDiff(start, end)

	var parts []string

	if years == 1 {
		parts = append(parts, "1 year")
	} else if years > 1 {
		parts = append(parts, fmt.Sprintf("%d years", years))
	}

	if months == 1 {
		parts = append(parts, "1 month")
	} else if months > 1 {
		parts = append(parts, fmt.Sprintf("%d months", months))
	}

	if weeks == 1 {
		parts = append(parts, "1 week")
	} else if weeks > 1 {
		parts = append(parts, fmt.Sprintf("%d weeks", weeks))
	}

	if years+months+weeks < 1 {
		if days == 1 {
			parts = append(parts, "1 day")
		} else if days > 1 {
			parts = append(parts, fmt.Sprintf("%d days", days))
		}

		if hours == 1 {
			parts = append(parts, "1 hour")
		} else if hours > 1 {
			parts = append(parts, fmt.Sprintf("%d hours", hours))
		}

		if days < 2 {
			if minutes == 1 {
				parts = append(parts, "1 minute")
			} else if minutes > 1 {
				parts = append(parts, fmt.Sprintf("%d minutes", minutes))
			}

			if seconds == 1 {
				parts = append(parts, "1 second")
			} else if seconds > 1 {
				parts = append(parts, fmt.Sprintf("%d seconds", seconds))
			}
		}
	}

	if len(parts) == 0 {
		return "now"
	}

	additional := ""

	if start.After(end) {
		additional = " ago"
	}

	return strings.Join(parts, ", ") + additional
}

func dateDiff(start, end time.Time) (years, months, weeks, days, hours, minutes, seconds int) {
	if end.Before(start) {
		start, end = end, start
	}

	// Calculate years and months
	years = end.Year() - start.Year()
	months = int(end.Month()) - int(start.Month())

	// Adjust if the end month is before the start month
	if months < 0 {
		years--
		months += 12
	}

	// Calculate days
	days = end.Day() - start.Day()
	if days < 0 {
		months--
		days += daysInMonth(start.Year(), start.Month())
	}

	// Convert days to weeks
	weeks = days / 7
	days = days % 7

	// Calculate hours
	hours = end.Hour() - start.Hour()
	if hours < 0 {
		days--
		hours += 24
	}

	// Calculate minutes
	minutes = end.Minute() - start.Minute()
	if minutes < 0 {
		hours--
		minutes += 60
	}

	// Calculate seconds
	seconds = end.Second() - start.Second()
	if seconds < 0 {
		minutes--
		seconds += 60
	}

	return
}

func daysInMonth(year int, month time.Month) int {
	return time.Date(year, month+1, 0, 0, 0, 0, 0, time.UTC).Day()
}

func IsWithinDuration(now time.Time, duration types.TimeDuration) bool {
	if now.Hour() < duration.From.Hour() {
		return false
	}

	if now.Hour() == duration.From.Hour() {
		if now.Minute() < duration.From.Minute() {
			return false
		}
	}

	if now.Hour() == duration.To.Hour() {
		if now.Minute() > duration.To.Minute() {
			return false
		}
	}

	if now.Hour() > duration.To.Hour() {
		return false
	}

	return true
}

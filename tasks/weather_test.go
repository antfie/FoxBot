package tasks

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWeatherCodeToDescription(t *testing.T) {
	assert.Equal(t, "Clear sky", weatherCodeToDescription(0))
	assert.Equal(t, "Mainly clear", weatherCodeToDescription(1))
	assert.Equal(t, "Partly cloudy", weatherCodeToDescription(2))
	assert.Equal(t, "Overcast", weatherCodeToDescription(3))
	assert.Equal(t, "Fog", weatherCodeToDescription(45))
	assert.Equal(t, "Fog", weatherCodeToDescription(48))
	assert.Equal(t, "Light drizzle", weatherCodeToDescription(51))
	assert.Equal(t, "Heavy rain", weatherCodeToDescription(65))
	assert.Equal(t, "Freezing rain", weatherCodeToDescription(66))
	assert.Equal(t, "Heavy snow", weatherCodeToDescription(75))
	assert.Equal(t, "Thunderstorm", weatherCodeToDescription(95))
	assert.Equal(t, "Thunderstorm with hail", weatherCodeToDescription(99))
	assert.Equal(t, "Unknown", weatherCodeToDescription(999))
}

func TestWeatherCodeToEmoji(t *testing.T) {
	assert.Equal(t, "☀️", weatherCodeToEmoji(0))
	assert.Equal(t, "🌤️", weatherCodeToEmoji(2))
	assert.Equal(t, "🌫️", weatherCodeToEmoji(45))
	assert.Equal(t, "🌦️", weatherCodeToEmoji(51))
	assert.Equal(t, "🌧️", weatherCodeToEmoji(65))
	assert.Equal(t, "🌨️", weatherCodeToEmoji(75))
	assert.Equal(t, "⛈️", weatherCodeToEmoji(95))
}

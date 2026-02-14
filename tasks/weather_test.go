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

func TestFormatWeatherForecast(t *testing.T) {
	data := openMeteoResponse{}
	data.Daily.TemperatureMax = []float64{18}
	data.Daily.TemperatureMin = []float64{12}
	data.Daily.PrecipitationProbMax = []int{30}
	data.Daily.WeatherCode = []int{0}
	data.Daily.WindSpeedMax = []float64{25}

	// 24 hours of hourly data
	hourlyTemps := make([]float64, 24)
	hourlyCodes := make([]int, 24)

	// Set specific values for morning (8), afternoon (13), evening (19)
	hourlyTemps[8] = 14
	hourlyCodes[8] = 2 // Partly cloudy
	hourlyTemps[13] = 18
	hourlyCodes[13] = 0 // Clear sky
	hourlyTemps[19] = 15
	hourlyCodes[19] = 63 // Moderate rain

	data.Hourly.Temperature = hourlyTemps
	data.Hourly.WeatherCode = hourlyCodes

	result := formatWeatherForecast("Manchester", data)

	expected := "☀️ Manchester: 12°C to 18°C\n" +
		"  Morning: 🌤️ Partly cloudy, 14°C\n" +
		"  Afternoon: ☀️ Clear sky, 18°C\n" +
		"  Evening: 🌧️ Moderate rain, 15°C\n" +
		"  Wind: up to 25 km/h | Rain: 30% chance"

	assert.Equal(t, expected, result)
}

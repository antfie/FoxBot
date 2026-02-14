package tasks

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/antfie/FoxBot/types"
	"github.com/antfie/FoxBot/utils"
)

type openMeteoResponse struct {
	Daily struct {
		TemperatureMax       []float64 `json:"temperature_2m_max"`
		TemperatureMin       []float64 `json:"temperature_2m_min"`
		PrecipitationProbMax []int     `json:"precipitation_probability_max"`
		WeatherCode          []int     `json:"weather_code"`
	} `json:"daily"`
}

func (c *Context) Weather() {
	if c.Config.Weather.Check.Duration != nil && !utils.IsWithinDuration(time.Now(), *c.Config.Weather.Check.Duration) {
		return
	}

	for _, location := range c.Config.Weather.Locations {
		c.fetchWeather(location)
	}
}

func (c *Context) fetchWeather(location types.WeatherLocation) {
	if c.DB.HasWeatherBeenNotifiedToday(location.Name) {
		return
	}

	url := fmt.Sprintf(
		"https://api.open-meteo.com/v1/forecast?latitude=%.4f&longitude=%.4f&daily=temperature_2m_max,temperature_2m_min,precipitation_probability_max,weather_code&timezone=auto&forecast_days=1",
		location.Latitude,
		location.Longitude,
	)

	response := utils.HttpRequest("GET", url, nil, nil)

	if response == nil {
		c.NotifyBad(fmt.Sprintf("Weather: Could not query API for %s", location.Name))
		return
	}

	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		c.NotifyBad(fmt.Sprintf("Weather: API returned status %s for %s", response.Status, location.Name))
		return
	}

	body, err := io.ReadAll(response.Body)

	if err != nil {
		log.Print(err)
		return
	}

	var data openMeteoResponse
	err = json.Unmarshal(body, &data)

	if err != nil {
		log.Printf("Weather: Could not parse response for %s: %v", location.Name, err)
		return
	}

	if len(data.Daily.TemperatureMax) < 1 || len(data.Daily.TemperatureMin) < 1 ||
		len(data.Daily.PrecipitationProbMax) < 1 || len(data.Daily.WeatherCode) < 1 {
		log.Printf("Weather: Incomplete data for %s", location.Name)
		return
	}

	condition := weatherCodeToDescription(data.Daily.WeatherCode[0])
	emoji := weatherCodeToEmoji(data.Daily.WeatherCode[0])

	message := fmt.Sprintf("%s %s: %.0f°C / %.0f°C, %s, %d%% chance of rain",
		emoji,
		location.Name,
		data.Daily.TemperatureMax[0],
		data.Daily.TemperatureMin[0],
		condition,
		data.Daily.PrecipitationProbMax[0],
	)

	c.Notify(message)
	c.DB.SetWeatherNotified(location.Name)
}

func weatherCodeToDescription(code int) string {
	switch code {
	case 0:
		return "Clear sky"
	case 1:
		return "Mainly clear"
	case 2:
		return "Partly cloudy"
	case 3:
		return "Overcast"
	case 45, 48:
		return "Fog"
	case 51:
		return "Light drizzle"
	case 53:
		return "Moderate drizzle"
	case 55:
		return "Dense drizzle"
	case 56, 57:
		return "Freezing drizzle"
	case 61:
		return "Slight rain"
	case 63:
		return "Moderate rain"
	case 65:
		return "Heavy rain"
	case 66, 67:
		return "Freezing rain"
	case 71:
		return "Slight snow"
	case 73:
		return "Moderate snow"
	case 75:
		return "Heavy snow"
	case 77:
		return "Snow grains"
	case 80:
		return "Slight rain showers"
	case 81:
		return "Moderate rain showers"
	case 82:
		return "Violent rain showers"
	case 85:
		return "Slight snow showers"
	case 86:
		return "Heavy snow showers"
	case 95:
		return "Thunderstorm"
	case 96, 99:
		return "Thunderstorm with hail"
	default:
		return "Unknown"
	}
}

func weatherCodeToEmoji(code int) string {
	switch {
	case code == 0:
		return "☀️"
	case code <= 3:
		return "🌤️"
	case code <= 48:
		return "🌫️"
	case code <= 57:
		return "🌦️"
	case code <= 67:
		return "🌧️"
	case code <= 77:
		return "🌨️"
	case code <= 82:
		return "🌦️"
	case code <= 86:
		return "🌨️"
	default:
		return "⛈️"
	}
}

package types

type Weather struct {
	Check     TimeFrequencyAndDuration
	Locations []WeatherLocation
}

type WeatherLocation struct {
	Name      string
	Latitude  float64
	Longitude float64
}

package shared

// Recipe represents a crop growth recipe.
type Recipe struct {
	Crop      string `json:"crop"`
	Condition struct {
		TemperatureMin       float64 `json:"temperature_min"`
		TemperatureMax       float64 `json:"temperature_max"`
		HumidityMin          float64 `json:"humidity_min"`
		HumidityMax          float64 `json:"humidity_max"`
		LiquidTemperatureMin float64 `json:"liquid_temperature_min"`
		LiquidTemperatureMax float64 `json:"liquid_temperature_max"`
		LiquidLevel          float64 `json:"liquid_level"`
		LightMin             float64 `json:"light_min"`
		LightMax             float64 `json:"light_max"`
		PHMin                float64 `json:"pH_min"`
		PHMax                float64 `json:"pH_max"`
		ECMin                float64 `json:"ec_min"`
		ECMax                float64 `json:"ec_max"`
		CO2Min               float64 `json:"co2_min"`
		CO2Max               float64 `json:"co2_max"`
	} `json:"condition"`
}

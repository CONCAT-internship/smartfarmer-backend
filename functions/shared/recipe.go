package shared

import "reflect"

// Recipe represents a crop growth recipe.
type Recipe struct {
	Crop      string `json:"crop"`
	Condition struct {
		TemperatureMin         float64 `json:"temperature_min"`
		TemperatureMax         float64 `json:"temperature_max"`
		HumidityMin            float64 `json:"humidity_min"`
		HumidityMax            float64 `json:"humidity_max"`
		LiquidTemperatureMin   float64 `json:"liquid_temperature"`
		TrayLiquidLevel        float64 `json:"tray_liquid_level"`
		LightMin               float64 `json:"light"`
		LightTime              float64 `json:"light_time"`
		PHMin                  float64 `json:"pH_min"`
		PHMax                  float64 `json:"pH_max"`
		ECMin                  float64 `json:"ec_min"`
		ECMax                  float64 `json:"ec_max"`
		PlantingDistanceWidth  float64 `json:"planting_distance_width"`
		PlantingDistanceHeight float64 `json:"planting_distnace_height"`
	} `json:"condition"`
}

func (r Recipe) ToMap() Document {
	doc := make(Document)
	val := reflect.ValueOf(r)
	typ := val.Type()
	for i := 0; i < typ.NumField(); i++ {
		tagname := typ.Field(i).Tag.Get("json")
		doc[tagname] = val.Field(i).Interface()
	}
	return doc
}

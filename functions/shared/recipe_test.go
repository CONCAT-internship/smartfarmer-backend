package shared

import (
	"reflect"
	"testing"
)

func TestRecipe_FromMap(t *testing.T) {
	var recipe = new(Recipe)
	recipe.FromMap(map[string]interface{}{
		"crop":                         "basil",
		"ec_max":                       1.5,
		"ec_min":                       int64(1),
		"humidity_max":                 int64(60),
		"humidity_min":                 int64(50),
		"light":                        int64(70),
		"light_time":                   int64(16),
		"liquid_temperature":           int64(20),
		"pH_max":                       6.5,
		"pH_min":                       int64(6),
		"planting_distance_min_width":  int64(20),
		"planting_distance_min_height": int64(20),
		"planting_distance_max_width":  int64(25),
		"planting_distance_max_height": int64(25),
		"temperature_max":              int64(30),
		"temperature_min":              int64(25),
		"tray_liquid_level":            int64(10),
	})
	var target = Recipe{
		Crop:                      "basil",
		ECMax:                     1.5,
		ECMin:                     1.0,
		HumidityMax:               60.0,
		HumidityMin:               50.0,
		Light:                     70.0,
		LightTime:                 16.0,
		LiquidTemperature:         20.0,
		PHMax:                     6.5,
		PHMin:                     6.0,
		PlantingDistanceMinWidth:  20.0,
		PlantingDistanceMinHeight: 20.0,
		PlantingDistanceMaxWidth:  25.0,
		PlantingDistanceMaxHeight: 25.0,
		TemperatureMax:            30.0,
		TemperatureMin:            25.0,
		TrayLiquidLevel:           10.0,
	}
	if !reflect.DeepEqual(*recipe, target) {
		t.Errorf("value mismatch: %v", *recipe)
	}
}

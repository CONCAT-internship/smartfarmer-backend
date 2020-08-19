package recipe

import "reflect"

// Recipe represents a crop growth recipe.
type Recipe struct {
	Crop                      string  `json:"crop"`
	FarmName                  string  `json:"farm_name"`
	TemperatureMin            float64 `json:"temperature_min"`
	TemperatureMax            float64 `json:"temperature_max"`
	HumidityMin               float64 `json:"humidity_min"`
	HumidityMax               float64 `json:"humidity_max"`
	LiquidTemperature         float64 `json:"liquid_temperature"`
	TrayLiquidLevel           float64 `json:"tray_liquid_level"`
	Light                     float64 `json:"light"`
	LightTime                 float64 `json:"light_time"`
	PHMin                     float64 `json:"pH_min"`
	PHMax                     float64 `json:"pH_max"`
	ECMin                     float64 `json:"ec_min"`
	ECMax                     float64 `json:"ec_max"`
	PlantingDistanceMinWidth  float64 `json:"planting_distance_min_width"`
	PlantingDistanceMinHeight float64 `json:"planting_distance_min_height"`
	PlantingDistanceMaxWidth  float64 `json:"planting_distance_max_width"`
	PlantingDistanceMaxHeight float64 `json:"planting_distance_max_height"`
}

// ToMap converts r to a Firestore document.
func (r Recipe) ToMap() map[string]interface{} {
	var doc = make(map[string]interface{})
	var val = reflect.ValueOf(r)
	var typ = val.Type()
	for i := 0; i < typ.NumField(); i++ {
		var tagname = typ.Field(i).Tag.Get("json")
		doc[tagname] = val.Field(i).Interface()
	}
	return doc
}

// FromMap binds a Firestore document to r.
func (r *Recipe) FromMap(doc map[string]interface{}) {
	var val = reflect.ValueOf(r).Elem()
	var typ = val.Type()
	for i := 0; i < typ.NumField(); i++ {
		var tagname = typ.Field(i).Tag.Get("json")
		switch doc[tagname].(type) {
		case int64: // avoid type mismatch
			val.Field(i).Set(reflect.ValueOf(float64(doc[tagname].(int64))))
		default:
			val.Field(i).Set(reflect.ValueOf(doc[tagname]))
		}
	}
}

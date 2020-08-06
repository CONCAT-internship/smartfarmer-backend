package shared

import "reflect"

// Document represents a single Firestore document.
type Document map[string]interface{}

// ToSensorData converts d to a sensor data.
func (d Document) ToSensorData() SensorData {
	var s = new(SensorData)
	var obj = reflect.ValueOf(s).Elem()
	var typ = obj.Type()
	for i := 0; i < typ.NumField(); i++ {
		var tagname = typ.Field(i).Tag.Get("json")
		obj.Field(i).Set(reflect.ValueOf(d[tagname]))
	}
	return *s
}

func (d Document) ToRecipe() Recipe {
	var r = new(Recipe)
	var obj = reflect.ValueOf(r).Elem()
	var typ = obj.Type()
	for i := 0; i < typ.NumField(); i++ {
		var tagname = typ.Field(i).Tag.Get("json")
		obj.Field(i).Set(reflect.ValueOf(d[tagname]))
	}
	return *r
}

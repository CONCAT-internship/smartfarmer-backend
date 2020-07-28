package shared

import "reflect"

// Document represents a single Firestore document.
type Document map[string]interface{}

// ToStruct converts d to sensord data struct.
func (d Document) ToStruct() SensorData {
	s := new(SensorData)
	obj := reflect.ValueOf(s).Elem()
	typ := obj.Type()
	for i := 0; i < typ.NumField(); i++ {
		tagname := typ.Field(i).Tag.Get("json")
		obj.Field(i).Set(reflect.ValueOf(d[tagname]))
	}
	return *s
}

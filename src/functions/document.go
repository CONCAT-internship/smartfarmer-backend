// Package functions defines sensor data model and implements related operations.
package functions

import "reflect"

// document represents a single Firestore document.
type document map[string]interface{}

// toStruct converts d to sensord data struct.
func (d document) toStruct() sensorData {
	s := new(sensorData)
	obj := reflect.ValueOf(s).Elem()
	typ := obj.Type()
	for i := 0; i < typ.NumField(); i++ {
		tagname := typ.Field(i).Tag.Get("json")
		obj.Field(i).Set(reflect.ValueOf(d[tagname]))
	}
	return *s
}

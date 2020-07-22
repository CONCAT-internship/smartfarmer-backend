// Package functions defines sensor data model and implements its CRUD operations.
package functions

import (
	"errors"
	"fmt"
	"reflect"
	"time"
)

// SensorData represents a single sensor data.
type SensorData struct {
	// unique id of arduino equipment
	UUID        string  `json:"uuid"`
	Temperature float64 `json:"temperature"`
	Humidity    float64 `json:"humidity"`
	PH          float64 `json:"pH"`
	// electrical conductivity
	EC float64 `json:"ec"`
	// light intensity
	Light             float64 `json:"light"`
	LiquidTemperature float64 `json:"liquid_temperature"`
	LiquidFlowRate    float64 `json:"liquid_flow_rate"`
	LiquidLevel       bool    `json:"liquid_level"`
	Valve             bool    `json:"valve"`
	LED               bool    `json:"led"`
	Fan               bool    `json:"fan"`
	// creation time of document
	UnixTime  int64     `json:"unix_time"`
	LocalTime time.Time `json:"local_time"`
}

// Document represents a single document of sensor data collection.
type Document map[string]interface{}

// validate checks whether there's any invalid value in s or not.
func (s SensorData) validate() error {
	var msg string
	if s.PH < 0 || s.PH > 14 {
		msg += fmt.Sprintf("Invalid value in pH: %f", s.PH)
	}
	if s.EC < 0 || s.EC > 2 {
		msg += fmt.Sprintf("Invalid value in ec: %f", s.EC)
	}
	if s.Light < 0 || s.Light > 100 {
		msg += fmt.Sprintf("Invalid value in light intensity: %f", s.Light)
	}
	if len(msg) > 0 {
		return errors.New(msg)
	}
	return nil
}

// setTime writes creation time to s.
func (s *SensorData) setTime() {
	s.LocalTime = time.Now()
	s.UnixTime = s.LocalTime.Unix()
}

// toMap converts s to Firestore document.
func (s SensorData) toMap() map[string]interface{} {
	doc := make(map[string]interface{})
	// copy values of s to doc
	val := reflect.ValueOf(s)
	typ := val.Type()
	for i := 0; i < typ.NumField(); i++ {
		tagname := typ.Field(i).Tag.Get("json")
		doc[tagname] = val.Field(i).Interface()
	}
	return doc
}

func (s SensorData) toLog(err error) map[string]string {
	doc := make(map[string]string)
	doc["uuid"] = s.UUID
	doc["message"] = err.Error()
	return doc
}

// toStruct converts d to sensord data struct.
func (d Document) toStruct() SensorData {
	s := new(SensorData)
	obj := reflect.ValueOf(s).Elem()
	typ := obj.Type()
	for i := 0; i < typ.NumField(); i++ {
		tagname := typ.Field(i).Tag.Get("json")
		obj.Field(i).Set(reflect.ValueOf(d[tagname]))
	}
	return *s
}

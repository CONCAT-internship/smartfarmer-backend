// Package functions defines sensor data model and implements related operations.
package functions

import (
	"errors"
	"fmt"
	"reflect"
	"time"
)

// sensorData represents a smart farm sensor data.
type sensorData struct {
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
	// data transmission time
	UnixTime  int64     `json:"unix_time"`
	LocalTime time.Time `json:"local_time"`
}

// setTime sets transmission time of s.
func (s *sensorData) setTime() {
	s.LocalTime = time.Now()
	s.UnixTime = s.LocalTime.Unix()
}

// validate checks whether the sensor works properly.
func (s sensorData) validate() error {
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

// appropriate checks whether the environment is suitable for crop growth.
func (s sensorData) appropriate() error {
	// TODO: change to make an error in case of inappropriate data
	return nil
}

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

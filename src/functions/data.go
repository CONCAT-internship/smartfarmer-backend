/*
Package functions defines sensor data model
and implements CRUD operations of smart farm.
*/
package functions

import (
	"errors"
	"fmt"
	"time"
)

// [Start fs_sensor_data_class]

// SensorData represents a single document of smart farm sensor data collection.
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
	UnixTime  int64     `json:"unix_time,omitempty"`
	LocalTime time.Time `json:"local_time,omitempty"`
}

// [End fs_sensor_data_class]

// validate checks whether there's any invalid value in s or not.
func (s SensorData) validate() error {
	var msg string
	if s.PH < 0 || s.PH > 14 {
		msg += fmt.Sprintf("Invalid value in pH: %f\n", s.PH)
	}
	if s.EC < 0 || s.EC > 2 {
		msg += fmt.Sprintf("Invalid value in ec: %f\n", s.EC)
	}
	if s.Light < 0 || s.Light > 100 {
		msg += fmt.Sprintf("Invalid value in light intensity: %f\n", s.Light)
	}
	if len(msg) > 0 {
		return errors.New(msg)
	}
	return nil
}

// toDocument converts sensor data to Firestore document.
func (s SensorData) toDocument() map[string]interface{} {
	doc := make(map[string]interface{})

	// copy data of s to doc
	doc["uuid"] = s.UUID
	doc["temperature"] = s.Temperature
	doc["humidity"] = s.Humidity
	doc["pH"] = s.PH
	doc["ec"] = s.EC
	doc["light"] = s.Light
	doc["liquid_temperature"] = s.LiquidTemperature
	doc["liquid_flow_rate"] = s.LiquidFlowRate
	doc["liquid_level"] = s.LiquidLevel
	doc["valve"] = s.Valve
	doc["led"] = s.LED
	doc["fan"] = s.Fan
	// update creation time of doc
	doc["local_time"] = time.Now()
	doc["unix_time"] = doc["local_time"].(time.Time).Unix()

	return doc
}

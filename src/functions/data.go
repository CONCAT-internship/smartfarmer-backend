// Package functions defines sensor data model
// and implements CRUD operations of smart farm.
package functions

import (
	"errors"
	"fmt"
	"time"
)

// [Start fs_sensor_data_class]

// sensorData represents a single document of smart farm sensor data collection.
type sensorData struct {
	uuid               string
	liquid_temperature float64
	temperature        float64
	humidity           float64
	liquid_flow_rate   float64
	pH                 float64
	ec                 float64
	light              float64
	liquid_level       bool
	valve              bool
	led                bool
	fan                bool
	unix_time          int64
	local_time         time.Time
}

// [End fs_sensor_data_class]

// setTime designates current time of s.
func (s *sensorData) setTime() {
	s.local_time = time.Now()
	s.unix_time = s.local_time.Unix()
}

// validate checks whether there's any invalid value in s or not.
func (s sensorData) validate() error {
	var msg string
	if s.pH < 0 || s.pH > 14 {
		msg += fmt.Sprintf("Invalid value in pH: %f\n", s.pH)
	}
	if s.ec < 0 || s.ec > 2 {
		msg += fmt.Sprintf("Invalid value in ec: %f\n", s.ec)
	}
	if s.light < 0 || s.light > 100 {
		msg += fmt.Sprintf("Invalid value in light intensity: %f\n", s.light)
	}
	if len(msg) > 0 {
		return errors.New(msg)
	}
	return nil
}

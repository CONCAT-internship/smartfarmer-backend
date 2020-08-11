package shared

import (
	"reflect"
	"testing"
	"time"
)

func TestSensorData_FromMap(t *testing.T) {
	var sensor_data = new(SensorData)
	sensor_data.FromMap(map[string]interface{}{
		"uuid":               "756e6b776f000c04",
		"temperature":        int64(30),
		"humidity":           int64(60),
		"pH":                 6.5,
		"ec":                 1.5,
		"light":              int64(70),
		"liquid_temperature": int64(20),
		"liquid_level":       false,
		"led":                true,
		"fan":                true,
		"unix_time":          int64(0),
		"local_time":         time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC),
		"light_time":         int64(0),
		"dark_time":          int64(0),
	})
	var target = SensorData{
		UUID:              "756e6b776f000c04",
		Temperature:       30.0,
		Humidity:          60.0,
		PH:                6.5,
		EC:                1.5,
		Light:             70.0,
		LiquidTemperature: 20.0,
		LiquidLevel:       false,
		LED:               true,
		Fan:               true,
		UnixTime:          int64(0),
		LocalTime:         time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC),
		LightTime:         0.0,
		DarkTime:          0.0,
	}
	if !reflect.DeepEqual(*sensor_data, target) {
		t.Errorf("value mismatch: %v", *sensor_data)
	}
}

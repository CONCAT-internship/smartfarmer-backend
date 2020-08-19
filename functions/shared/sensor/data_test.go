package sensor

import (
	"reflect"
	"testing"
	"time"
)

func TestData(t *testing.T) {
	var data = new(Data)
	data.FromMap(map[string]interface{}{
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
	var target = Data{
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
	if !reflect.DeepEqual(*data, target) {
		t.Errorf("value mismatch: %v", *data)
	}
}

package shared

import (
	"reflect"
	"time"
)

const (
	PH_MIN = 0
	PH_MAX = 14

	EC_MIN = 0
	EC_MAX = 2

	LIGHT_MIN = 0
	LIGHT_MAX = 100

	PH_INC  = 1
	PH_KEEP = 0
	PH_DEC  = -1

	EC_INC  = 1
	EC_KEEP = 0
	// no way to decrease the value in the EC pump
)

// SensorData represents a smart farm sensor data.
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
	LiquidLevel       bool    `json:"liquid_level"`
	LED               bool    `json:"led"`
	Fan               bool    `json:"fan"`
	// data transmission time
	UnixTime  int64     `json:"unix_time"`
	LocalTime time.Time `json:"local_time"`
}

// SetTime sets transmission time of s.
func (s *SensorData) SetTime() {
	s.LocalTime = time.Now()
	s.UnixTime = s.LocalTime.Unix()
}

// Validate checks whether s works normally and the data of it is appropriate.
func (s SensorData) Validate(r Recipe) (errorcodes []int) {
	if s.PH < PH_MIN || s.PH > PH_MAX {
		errorcodes = append(errorcodes, CODE_PH_MALFUNC)
	} else if s.PH > r.Condition.PHMax {
		errorcodes = append(errorcodes, CODE_PH_IMPROPER_HIGH)
	} else if s.PH < r.Condition.PHMin {
		errorcodes = append(errorcodes, CODE_PH_IMPROPER_LOW)
	}
	if s.EC < EC_MIN || s.EC > EC_MAX {
		errorcodes = append(errorcodes, CODE_EC_MALFUNC)
	} else if s.EC > r.Condition.ECMax {
		errorcodes = append(errorcodes, CODE_EC_IMPROPER_HIGH)
	} else if s.EC < r.Condition.ECMin {
		errorcodes = append(errorcodes, CODE_EC_IMPROPER_LOW)
	}
	if s.Light < 0 || s.Light > 100 {
		errorcodes = append(errorcodes, CODE_LIGHT_MALFUNC)
	}
	if s.Temperature > r.Condition.TemperatureMax {
		errorcodes = append(errorcodes, CODE_TEMPERATURE_IMPROPER_HIGH)
	}
	if s.Temperature < r.Condition.TemperatureMin {
		errorcodes = append(errorcodes, CODE_TEMPERATURE_IMPROPER_LOW)
	}
	if s.Humidity > r.Condition.HumidityMax {
		errorcodes = append(errorcodes, CODE_HUMIDITY_IMPROPER_HIGH)
	}
	if s.Humidity < r.Condition.HumidityMin {
		errorcodes = append(errorcodes, CODE_HUMIDITY_IMPROPER_LOW)
	}
	return
}

// ToMap converts s to a Firestore document.
func (s SensorData) ToMap() Document {
	var doc = make(Document)
	var val = reflect.ValueOf(s)
	var typ = val.Type()
	for i := 0; i < typ.NumField(); i++ {
		var tagname = typ.Field(i).Tag.Get("json")
		doc[tagname] = val.Field(i).Interface()
	}
	return doc
}

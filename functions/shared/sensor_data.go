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
	UUID              string    `json:"uuid"` // unique id of arduino equipment
	Temperature       float64   `json:"temperature"`
	Humidity          float64   `json:"humidity"`
	PH                float64   `json:"pH"`
	EC                float64   `json:"ec"`    // electrical conductivity
	Light             float64   `json:"light"` // light intensity
	LiquidTemperature float64   `json:"liquid_temperature"`
	LiquidLevel       bool      `json:"liquid_level"`
	LED               bool      `json:"led"`
	Fan               bool      `json:"fan"`
	UnixTime          int64     `json:"unix_time"` // data transmission time
	LocalTime         time.Time `json:"local_time"`
}

// SetTime sets transmission time of s.
func (s *SensorData) SetTime() {
	s.LocalTime = time.Now()
	s.UnixTime = s.LocalTime.Unix()
}

// Validate checks whether s works normally and it is appropriate for crop growth.
func (s SensorData) Validate(r Recipe) ([]ErrorCode, map[string]interface{}) {
	var errorCodes []ErrorCode
	var requirements = map[string]interface{}{
		"pH_pump": PH_KEEP,
		"ec_pump": EC_KEEP,
	}
	if s.PH < PH_MIN || s.PH > PH_MAX {
		errorCodes = append(errorCodes, CODE_PH_MALFUNC)
	} else if s.PH > r.PHMax {
		errorCodes = append(errorCodes, CODE_PH_IMPROPER_HIGH)
		requirements["pH_pump"] = PH_DEC
	} else if s.PH < r.PHMin {
		errorCodes = append(errorCodes, CODE_PH_IMPROPER_LOW)
		requirements["pH_pump"] = PH_INC
	}
	if s.EC < EC_MIN || s.EC > EC_MAX {
		errorCodes = append(errorCodes, CODE_EC_MALFUNC)
	} else if s.EC > r.ECMax {
		errorCodes = append(errorCodes, CODE_EC_IMPROPER_HIGH)
	} else if s.EC < r.ECMin {
		errorCodes = append(errorCodes, CODE_EC_IMPROPER_LOW)
		requirements["ec_pump"] = EC_INC
	}
	if s.Light < 0 || s.Light > 100 {
		errorCodes = append(errorCodes, CODE_LIGHT_MALFUNC)
	}
	if s.Temperature > r.TemperatureMax {
		errorCodes = append(errorCodes, CODE_TEMPERATURE_IMPROPER_HIGH)
		requirements["fan"] = true
	}
	if s.Temperature < r.TemperatureMin {
		errorCodes = append(errorCodes, CODE_TEMPERATURE_IMPROPER_LOW)
		requirements["fan"] = false
	}
	if s.Humidity > r.HumidityMax {
		errorCodes = append(errorCodes, CODE_HUMIDITY_IMPROPER_HIGH)
	}
	if s.Humidity < r.HumidityMin {
		errorCodes = append(errorCodes, CODE_HUMIDITY_IMPROPER_LOW)
	}
	return errorCodes, requirements
}

// ToMap converts s to a Firestore document.
func (s SensorData) ToMap() map[string]interface{} {
	var doc = make(map[string]interface{})
	var val = reflect.ValueOf(s)
	var typ = val.Type()
	for i := 0; i < typ.NumField(); i++ {
		var tagname = typ.Field(i).Tag.Get("json")
		doc[tagname] = val.Field(i).Interface()
	}
	return doc
}

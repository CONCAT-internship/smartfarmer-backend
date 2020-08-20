package sensor

import (
	"reflect"
	"time"

	"github.com/maengsanha/smartfarmer-backend/functions/shared/code"
	"github.com/maengsanha/smartfarmer-backend/functions/shared/recipe"
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

	TRANSMISSION_CYCLE = 3 // data transmission cycle of from device
)

// Data represents a smart farm sensor data.
type Data struct {
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
	LightTime         float64   `json:"light_time"`
	DarkTime          float64   `json:"dark_time"`
}

// SetTime sets transmission time of d.
func (d *Data) SetTime() {
	d.LocalTime = time.Now()
	d.UnixTime = d.LocalTime.Unix()
}

// Validate checks whether d works normally and is appropriate for crop growth.
func (d Data) Validate(r recipe.Recipe) ([]code.Code, map[string]interface{}) {
	var codes []code.Code
	var requirements = map[string]interface{}{
		"pH_pump": PH_KEEP,
		"ec_pump": EC_KEEP,
	}
	if d.PH < PH_MIN || d.PH > PH_MAX {
		codes = append(codes, code.PH_MALFUNC)
	} else if d.PH > r.PHMax {
		codes = append(codes, code.PH_IMPROPER_HIGH)
		requirements["pH_pump"] = PH_DEC
	} else if d.PH < r.PHMin {
		codes = append(codes, code.PH_IMPROPER_LOW)
		requirements["pH_pump"] = PH_INC
	}
	if d.EC < EC_MIN || d.EC > EC_MAX {
		codes = append(codes, code.EC_MALFUNC)
	} else if d.EC > r.ECMax {
		codes = append(codes, code.EC_IMPROPER_HIGH)
	} else if d.EC < r.ECMin {
		codes = append(codes, code.EC_IMPROPER_LOW)
		requirements["ec_pump"] = EC_INC
	}
	if d.Light < LIGHT_MIN || d.Light > LIGHT_MAX {
		codes = append(codes, code.LIGHT_MALFUNC)
	}
	if d.Temperature > r.TemperatureMax {
		codes = append(codes, code.TEMPERATURE_IMPROPER_HIGH)
		requirements["fan"] = true
	}
	if d.Temperature < r.TemperatureMin {
		codes = append(codes, code.TEMPERATURE_IMPROPER_LOW)
		requirements["fan"] = false
	}
	if d.Humidity > r.HumidityMax {
		codes = append(codes, code.HUMIDITY_IMPROPER_HIGH)
	}
	if d.Humidity < r.HumidityMin {
		codes = append(codes, code.HUMIDITY_IMPROPER_LOW)
	}
	return codes, requirements
}

// ToMap converts d to a Firestore document.
func (d Data) ToMap() map[string]interface{} {
	var doc = make(map[string]interface{})
	var val = reflect.ValueOf(d)
	var typ = val.Type()
	for i := 0; i < typ.NumField(); i++ {
		var tagname = typ.Field(i).Tag.Get("json")
		doc[tagname] = val.Field(i).Interface()
	}
	return doc
}

// FromMap binds a Firestore document to d.
func (d *Data) FromMap(doc map[string]interface{}) {
	var val = reflect.ValueOf(d).Elem()
	var typ = val.Type()
	for i := 0; i < typ.NumField(); i++ {
		var tagname = typ.Field(i).Tag.Get("json")
		if tagname == "unix_time" {
			val.Field(i).Set(reflect.ValueOf(doc[tagname]))
		} else {
			switch doc[tagname].(type) { // avoid type mismatch
			case int64:
				val.Field(i).Set(reflect.ValueOf(float64(doc[tagname].(int64))))
			default:
				val.Field(i).Set(reflect.ValueOf(doc[tagname]))
			}
		}
	}
}

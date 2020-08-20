package code

/*
Code represents a code.
See https://github.com/CONCAT-internship/smartfarmer-backend/blob/master/README.md
*/
type Code int

const (
	// case of the device failed to send data
	DATA_EMPTY Code = 7000
	// case of the device did not measure properly
	PH_MALFUNC    Code = 4000
	EC_MALFUNC    Code = 4001
	LIGHT_MALFUNC Code = 4002
	// case of the device measured properly but the environment is improper
	PH_IMPROPER_HIGH          Code = 5000
	PH_IMPROPER_LOW           Code = 5001
	EC_IMPROPER_HIGH          Code = 5010
	EC_IMPROPER_LOW           Code = 5011
	TEMPERATURE_IMPROPER_HIGH Code = 5020
	TEMPERATURE_IMPROPER_LOW  Code = 5021
	HUMIDITY_IMPROPER_HIGH    Code = 5030
	HUMIDITY_IMPROPER_LOW     Code = 5031
)

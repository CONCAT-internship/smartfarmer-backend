package shared

const (
	// case of the device failed to send data
	CODE_DATA_EMPTY = 7000
	// case of the device did not measure properly
	CODE_PH_MALFUNC    = 4000
	CODE_EC_MALFUNC    = 4001
	CODE_LIGHT_MALFUNC = 4002
	// case of the device measured properly but the environment is improper
	CODE_PH_IMPROPER_HIGH          = 5000
	CODE_PH_IMPROPER_LOW           = 5001
	CODE_EC_IMPROPER_HIGH          = 5010
	CODE_EC_IMPROPER_LOW           = 5011
	CODE_TEMPERATURE_IMPROPER_HIGH = 5020
	CODE_TEMPERATURE_IMPROPER_LOW  = 5021
	CODE_HUMIDITY_IMPROPER_HIGH    = 5030
	CODE_HUMIDITY_IMPROPER_LOW     = 5031
)

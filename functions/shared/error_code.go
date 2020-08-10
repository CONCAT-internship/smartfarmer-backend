package shared

// ErrorCode represents an error code.
type ErrorCode int

const (
	// case of the device failed to send data
	CODE_DATA_EMPTY ErrorCode = 7000
	// case of the device did not measure properly
	CODE_PH_MALFUNC    ErrorCode = 4000
	CODE_EC_MALFUNC    ErrorCode = 4001
	CODE_LIGHT_MALFUNC ErrorCode = 4002
	// case of the device measured properly but the environment is improper
	CODE_PH_IMPROPER_HIGH          ErrorCode = 5000
	CODE_PH_IMPROPER_LOW           ErrorCode = 5001
	CODE_EC_IMPROPER_HIGH          ErrorCode = 5010
	CODE_EC_IMPROPER_LOW           ErrorCode = 5011
	CODE_TEMPERATURE_IMPROPER_HIGH ErrorCode = 5020
	CODE_TEMPERATURE_IMPROPER_LOW  ErrorCode = 5021
	CODE_HUMIDITY_IMPROPER_HIGH    ErrorCode = 5030
	CODE_HUMIDITY_IMPROPER_LOW     ErrorCode = 5031
)

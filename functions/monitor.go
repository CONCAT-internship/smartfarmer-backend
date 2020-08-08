package functions

import (
	"context"
	"net/http"
	"time"
)

// Monitor checks the records for twice the duration of the insertion cycle of sensor data.
// exported to https://asia-northeast1-superfarmers.cloudfunctions.net/Monitor
func Monitor(writer http.ResponseWriter, request *http.Request) {
	defer request.Body.Close()

	const cycle = 3
	var base = time.Now().Unix() - 2*cycle*60

	var cursor = client.Collection("sensor_data").
		Where("unix_time", ">=", base).
		Documents(context.Background())

	_ = cursor
}

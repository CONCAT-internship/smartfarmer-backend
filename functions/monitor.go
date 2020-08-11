package functions

import (
	"context"
	"net/http"
	"time"

	"github.com/joshua-dev/smartfarmer-backend/functions/shared"
)

// Monitor checks whether the devices are working properly.
// exported to https://asia-northeast1-superfarmers.cloudfunctions.net/Monitor
func Monitor(writer http.ResponseWriter, request *http.Request) {
	defer request.Body.Close()

	var base = time.Now().Unix() - 2*shared.TRANSMISSION_CYCLE*60

	var cursor = client.Collection("sensor_data").
		Where("unix_time", ">=", base).
		Documents(context.Background())

	_ = cursor
}

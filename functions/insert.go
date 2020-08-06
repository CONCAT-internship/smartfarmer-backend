package functions

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/joshua-dev/smartfarmer-backend/functions/shared"
)

// Insert stores a sensor data into Firestore.
// exported to https://asia-northeast1-superfarmers.cloudfunctions.net/Insert
func Insert(writer http.ResponseWriter, request *http.Request) {
	var data = new(shared.SensorData)
	if err := json.NewDecoder(request.Body).Decode(data); err != nil {
		http.Error(writer, fmt.Sprintf("json.Decode: %v", err), http.StatusInternalServerError)
		return
	}
	defer request.Body.Close()

	data.SetTime()

	if _, _, err := db.Collection("sensor_data").
		Add(context.Background(), data.ToMap()); err != nil {
		http.Error(writer, fmt.Sprintf("firestore.Add: %v", err), http.StatusInternalServerError)
		return
	}
}

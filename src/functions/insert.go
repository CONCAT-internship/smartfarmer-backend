// Package functions defines sensor data model and implements related operations.
package functions

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"cloud.google.com/go/firestore"
)

// Insert stores a sensor data into Firestore.
// exported to https://asia-northeast1-superfarmers.cloudfunctions.net/Insert
func Insert(writer http.ResponseWriter, req *http.Request) {
	ctx := context.Background()

	client, err := firestore.NewClient(ctx, PROJECT_ID)
	if err != nil {
		http.Error(writer, fmt.Sprintf("firestore.NewClient: %v", err), http.StatusInternalServerError)
		return
	}
	defer client.Close()

	// parse JSON -> sensor data
	data := new(sensorData)
	if err = json.NewDecoder(req.Body).Decode(data); err != nil {
		http.Error(writer, fmt.Sprintf("json.Decode: %v", err), http.StatusInternalServerError)
		return
	}
	defer req.Body.Close()

	// update transmission time
	data.setTime()

	// store data into Firestore
	if _, _, err = client.Collection("sensor_data").Add(ctx, data.toMap()); err != nil {
		http.Error(writer, fmt.Sprintf("firestore.Add: %v", err), http.StatusInternalServerError)
		return
	}
}

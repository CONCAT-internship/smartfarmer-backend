// Package functions defines sensor data model and implements its CRUD operations.
package functions

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"cloud.google.com/go/firestore"
)

const projectID = "superfarmers"

// Insert stores a sensor data into Firestore.
// exported to https://asia-northeast1-superfarmers.cloudfunctions.net/Insert
func Insert(writer http.ResponseWriter, req *http.Request) {
	ctx := context.Background()

	// create new firestore client and close it when query is done
	client, err := firestore.NewClient(ctx, projectID)
	if err != nil {
		fmt.Fprintf(writer, "firestore.NewClient: %v\n", err)
		return
	}
	defer client.Close()

	// create new sensor data and parse json from request body
	var data SensorData
	if err = json.NewDecoder(req.Body).Decode(&data); err != nil {
		fmt.Fprintf(writer, "json.Decode: %v\n", err)
		return
	}
	defer req.Body.Close()

	data.setTime()

	// validates data
	if err = data.validate(); err != nil {
		fmt.Fprintf(writer, "validation failed: %v\n", err)
	}

	// store data into collection
	if _, _, err = client.Collection("sensor_data").Add(ctx, data.toMap()); err != nil {
		fmt.Fprintf(writer, "firestore.Add: %v\n", err)
		return
	}
	fmt.Fprintln(writer, "Successfully stored to Firestore.")
}

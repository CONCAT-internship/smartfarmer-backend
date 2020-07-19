// Package functions defines sensor data model
// and implements CRUD operations of smart farm.
package functions

// [Start fs_functions_dependencies]
import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"cloud.google.com/go/firestore"
)

// [End fs_functions_dependencies]

const projectID = "superfarmers"

// [Start fs_functions_insert]

// insert stores a sensor data into Firestore.
// exported to https://asia-northeast1-superfarmers.cloudfunctions.net/insert
func insert(writer http.ResponseWriter, req *http.Request) {
	ctx := context.Background()

	// create new firestore client and close it when query is done
	client, err := firestore.NewClient(ctx, projectID)
	if err != nil {
		fmt.Fprintf(writer, "firestore.NewClient: %v", err)
		return
	}
	defer client.Close()

	// create new sensor data and parse json from request body
	var data sensorData
	if err = json.NewDecoder(req.Body).Decode(&data); err != nil {
		fmt.Fprintf(writer, "json.Decode: %v", err)
		return
	}
	// update creation time of sensor data
	data.setTime()

	// validates data
	if err = data.validate(); err != nil {
		fmt.Fprintf(writer, "validation failed: %v", err)
	}

	// store data into collection
	if _, _, err = client.Collection("sensor_data").Add(ctx, data); err != nil {
		fmt.Fprintf(writer, "firestore.Add: %v", err)
		return
	}
}

// [End fs_functions_insert]

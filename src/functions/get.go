// Package functions defines sensor data model and implements its CRUD operations.
package functions

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"cloud.google.com/go/firestore"
)

// Get brings records for the last week from Firestore with given uuid.
// exported to https://asia-northeast1-superfarmers.cloudfunctions.net/Get
func Get(writer http.ResponseWriter, req *http.Request) {
	ctx := context.Background()

	// create new firestore client
	client, err := firestore.NewClient(ctx, projectID)
	if err != nil {
		fmt.Fprintf(writer, "firestore.NewClient: %v", err)
		return
	}
	defer client.Close()

	// parse uuid from query string
	uuid := req.URL.Query().Get("uuid")
	defer req.Body.Close()

	// create response body
	var resp struct {
		Records []SensorData `json:"records"`
	}

	// get records for the last week from Firestore in descending order
	records, err := client.Collection("sensor_data").Where("uuid", "==", uuid).OrderBy("unix_time", firestore.Desc).Where("unix_time", ">=", time.Now().Unix()-7*24*60*60).Documents(ctx).GetAll()
	if err != nil {
		fmt.Fprintf(writer, "firestore.GetAll: %v", err)
	}

	for _, record := range records {
		resp.Records = append(resp.Records, Document(record.Data()).toStruct())
	}

	// notify that it's a JSON response
	writer.Header().Set("Content-Type", "application/json")
	if err = json.NewEncoder(writer).Encode(resp); err != nil {
		fmt.Fprintf(writer, "json.Encode: %v", err)
		return
	}
}

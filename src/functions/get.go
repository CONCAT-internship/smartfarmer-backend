/*
Package functions defines sensor data model
and implements CRUD operations of smart farm.
*/
package functions

// [Start fs_functions_dependencies]
import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"cloud.google.com/go/firestore"
)

// [End fs_functions_dependencies]

// [Start fs_functions_get]

/*
Get brings records for the last week from Firestore with given uuid.
exported to https://asia-northeast1-superfarmers.cloudfunctions.net/Get
*/
func Get(writer http.ResponseWriter, req *http.Request) {
	ctx := context.Background()

	// create new firestore client
	client, err := firestore.NewClient(ctx, projectID)
	if err != nil {
		fmt.Fprintf(writer, "firestore.NewClient: %v\n", err)
		return
	}
	defer client.Close()

	// parse uuid from req
	var rqst struct {
		UUID string `json:"uuid"`
	}
	if err = json.NewDecoder(req.Body).Decode(&rqst); err != nil {
		fmt.Fprintf(writer, "json.Decode: %v\n", err)
		return
	}
	defer req.Body.Close()

	// create response body
	var resp struct {
		Records []SensorData `json:"records"`
	}

	now := time.Now().Unix()
	const weekTime = 7 * 24 * 60 * 60
	data := new(SensorData)

	records, err := client.Collection("sensor_data").Where("uuid", "==", rqst.UUID).Where("unix_time", ">=", now-weekTime).Documents(ctx).GetAll()
	if err != nil {
		fmt.Fprintf(writer, "firestore.GetAll: %v\n", err)
	}

	for _, record := range records {
		data.fromMap(record.Data())
		resp.Records = append(resp.Records, *data)
	}
	data = nil

	// notify that it's a JSON response
	writer.Header().Set("Content-Type", "application/json")
	if err = json.NewEncoder(writer).Encode(resp); err != nil {
		fmt.Fprintf(writer, "json.Encode: %v\n", err)
		return
	}
}

// [End fs_functions_get]

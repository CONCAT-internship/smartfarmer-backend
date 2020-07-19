// Package functions defines sensor data model
// and implements CRUD operations of smart farm.
package functions

// [Start fs_functions_dependencies]
import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"google.golang.org/api/iterator"

	"cloud.google.com/go/firestore"
)

// [End fs_functions_dependencies]

// [Start fs_functions_get]

// get brings records for the last week from Firestore with given uuid.
// exported to https://asia-northeast1-superfarmers.cloudfunctions.net/get
func get(writer http.ResponseWriter, req *http.Request) {
	ctx := context.Background()

	// create new firestore client
	client, err := firestore.NewClient(ctx, projectID)
	if err != nil {
		fmt.Fprintf(writer, "firestore.NewClient: %v\n", err)
		return
	}
	defer client.Close()

	// parse uuid from query string
	uuid := req.URL.Query().Get("uuid")
	if len(uuid) < 1 {
		fmt.Fprintln(writer, "URL param uuid is missing")
		return
	}

	// create response body
	var resp struct {
		records []sensorData
	}

	now := time.Now().Unix()
	const weekTime = 7 * 24 * 60 * 60

	// create temporal pointer to store document
	tmp := new(sensorData)

	cursor := client.Collection("sensor_data").Where("uuid", "==", uuid).Where("unix_time", ">=", now-weekTime).Documents(ctx)
	for {
		doc, err := cursor.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			fmt.Fprintf(writer, "firestore.Next: %v\n", err)
		}
		doc.DataTo(tmp)
		resp.records = append(resp.records, *tmp)
	}
	tmp = nil

	// notify that it's a JSON response
	writer.Header().Set("Content-Type", "application/json")
	if err = json.NewEncoder(writer).Encode(resp); err != nil {
		fmt.Fprintf(writer, "json.Encode: %v\n", err)
		return
	}
}

// [End fs_functions_get]

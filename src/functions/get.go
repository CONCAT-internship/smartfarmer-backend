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

	// parse uuid from query string
	uuid := req.URL.Query().Get("uuid")
	if len(uuid) < 1 {
		fmt.Fprintln(writer, "URL param uuid is missing")
		return
	}

	// create response body
	var resp struct {
		records []SensorData
	}

	now := time.Now().Unix()
	const weekTime = 7 * 24 * 60 * 60

	records, err := client.Collection("sensor_data").Where("uuid", "==", uuid).Where("unix_time", ">=", now-weekTime).Documents(ctx).GetAll()
	if err != nil {
		fmt.Fprintf(writer, "firestore.GetAll: %v\n", err)
	}

	// notify that it's a JSON response
	writer.Header().Set("Content-Type", "application/json")
	if err = json.NewEncoder(writer).Encode(resp); err != nil {
		fmt.Fprintf(writer, "json.Encode: %v\n", err)
		return
	}
}

// [End fs_functions_get]

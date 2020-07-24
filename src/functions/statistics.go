// Package functions defines sensor data model and implements related operations.
package functions

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"cloud.google.com/go/firestore"
)

// DailyAverage calculates the average of the sensor data for each day of the week.
// exported to https://asia-northeast1-superfarmers.cloudfunctions.net/DailyAverage
func DailyAverage(writer http.ResponseWriter, req *http.Request) {
	ctx := context.Background()

	client, err := firestore.NewClient(ctx, projectID)
	if err != nil {
		fmt.Fprintf(writer, "firestore.NewClient: %v", err)
		return
	}
	defer client.Close()

	uuid := req.URL.Query().Get("uuid")
	base, err := strconv.Atoi(req.URL.Query().Get("unixtime"))
	if err != nil {
		fmt.Fprintf(writer, "strconv.Atoi: %v", err)
		return
	}
	defer req.Body.Close()

	snapshot, err := client.Collection("sensor_data").Where("uuid", "==", uuid).OrderBy("unix_time", firestore.Asc).Where("unix_time", ">=", base).Where("unix_time", "<", base+7*24*60*60).Documents(ctx).GetAll()
	if err != nil {
		fmt.Fprintf(writer, "firestore.GetAll: %v", err)
		return
	}
	_ = snapshot

	writer.Header().Set("Content-Type", "application/json")
	if err = json.NewEncoder(writer).Encode(nil); err != nil {
		fmt.Fprintf(writer, "json.Encode: %v", err)
		return
	}
}

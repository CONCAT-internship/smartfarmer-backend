package functions

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"google.golang.org/api/iterator"

	"cloud.google.com/go/firestore"
)

// RecentStatus returns the latest status of the farm.
// exported to https://asia-northeast1-superfarmers.cloudfunctions.net/RecentStatus
func RecentStatus(writer http.ResponseWriter, request *http.Request) {
	var uuid = request.URL.Query().Get("uuid")
	defer request.Body.Close()

	doc, err := client.Collection("sensor_data").
		Where("uuid", "==", uuid).
		OrderBy("unix_time", firestore.Desc).
		Limit(1).
		Documents(context.Background()).
		Next()

	if err != iterator.Done && err != nil {
		http.Error(writer, fmt.Sprintf("firestore.Next: %v", err), http.StatusInternalServerError)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	if err = json.NewEncoder(writer).Encode(doc.Data()); err != nil {
		http.Error(writer, fmt.Sprintf("json.Encode: %v", err), http.StatusInternalServerError)
		return
	}
}

package functions

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/joshua-dev/smartfarmer-backend/functions/shared"

	"google.golang.org/api/iterator"

	"cloud.google.com/go/firestore"
)

// RecentStatus returns recent device status.
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
	var data = shared.Document(doc.Data()).ToStruct()

	writer.Header().Set("Content-Type", "application/json")
	if err = json.NewEncoder(writer).Encode(map[string]bool{
		"valve": data.Valve,
		"led":   data.LED,
		"fan":   data.Fan,
	}); err != nil {
		http.Error(writer, fmt.Sprintf("json.Encode: %v", err), http.StatusInternalServerError)
		return
	}
}

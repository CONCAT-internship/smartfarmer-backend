package functions

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// DesiredStatus returns the desired status of the device.
// exported to https://asia-northeast1-superfarmers.cloudfunctions.net/DesiredStatus
func DesiredStatus(writer http.ResponseWriter, request *http.Request) {
	var uuid = request.URL.Query().Get("uuid")
	defer request.Body.Close()

	doc, err := client.Collection("desired_status").
		Doc(uuid).
		Get(context.Background())

	if err != nil {
		http.Error(writer, fmt.Sprintf("firestore.Get: %v", err), http.StatusInternalServerError)
		return
	}
	writer.Header().Set("Content-Type", "application/json")
	if err = json.NewEncoder(writer).Encode(doc.Data()); err != nil {
		http.Error(writer, fmt.Sprintf("json.Encode: %v", err), http.StatusInternalServerError)
		return
	}
}

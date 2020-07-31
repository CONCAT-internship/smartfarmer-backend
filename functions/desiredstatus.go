package functions

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/joshua-dev/smartfarmer-backend/functions/shared"
)

// DesiredStatus returns the desired status of the device.
// exported to https://asia-northeast1-superfarmers.cloudfunctions.net/DesiredStatus
func DesiredStatus(writer http.ResponseWriter, request *http.Request) {
	var uuid = request.URL.Query().Get("uuid")
	defer request.Body.Close()
	if doc, err := client.Doc("desired_status/" + uuid).
		Get(context.Background()); err != nil {
		http.Error(writer, fmt.Sprintf("firestore.Get: %v", err), http.StatusInternalServerError)
		return
	} else {
		writer.Header().Set("Content-Type", "application/json")
		var data = shared.Document(doc.Data()).ToStruct()
		if err = json.NewEncoder(writer).Encode(map[string]interface{}{
			"valve": data.Valve,
			"led":   data.LED,
			"fan":   data.Fan,
		}); err != nil {
			http.Error(writer, fmt.Sprintf("json.Encode: %v", err), http.StatusInternalServerError)
			return
		}
	}
}

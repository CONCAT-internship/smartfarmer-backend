package functions

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// RegisterDevice appends a new device information to the existing device list.
// exported to https://asia-northeast1-superfarmers.cloudfunctions.net/RegisterDevice
func RegisterDevice(writer http.ResponseWriter, request *http.Request) {
	var data = new(struct {
		Email string `json:"email"`
		UUID  string `json:"uuid"`
	})
	if err := json.NewDecoder(request.Body).Decode(data); err != nil {
		http.Error(writer, fmt.Sprintf("json.Decode: %v", err), http.StatusInternalServerError)
		return
	}
	defer request.Body.Close()

	db.Collection("farmers").
		Where("email", "==", data.Email).
		Documents(context.Background()).
		Next()
}

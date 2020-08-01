package functions

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// Control changes the device state by chaning the desired status.
// exported to https://asia-northeast1-superfarmers.cloudfunctions.net/Control
func Control(writer http.ResponseWriter, request *http.Request) {
	var desired = new(struct {
		UUID   string `json:"uuid"`
		Status struct {
			Valve bool `json:"valve"`
			LED   bool `json:"led"`
			Fan   bool `json:"fan"`
		} `json:"status"`
	})
	if err := json.NewDecoder(request.Body).Decode(desired); err != nil {
		http.Error(writer, fmt.Sprintf("json.Decode: %v", err), http.StatusBadRequest)
		return
	}
	defer request.Body.Close()

	if _, err := client.Collection("desired_status").
		Doc(desired.UUID).
		Set(context.Background(), map[string]bool{
			"valve": desired.Status.Valve,
			"led":   desired.Status.LED,
			"fan":   desired.Status.Fan,
		}); err != nil {
		http.Error(writer, fmt.Sprintf("firestore.Set: %v", err), http.StatusInternalServerError)
		return
	}
}

package functions

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/joshua-dev/smartfarmer-backend/functions/shared"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"
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

	doc, err := client.Collection("farmers").
		Where("email", "==", data.Email).
		Documents(context.Background()).
		Next()

	if err != iterator.Done && err != nil {
		http.Error(writer, fmt.Sprintf("firestore.Next: %v", err), http.StatusInternalServerError)
		return
	}
	if _, err = doc.Ref.Set(context.Background(),
		map[string][]string{
			"device_uuid": append(doc.Data()["device_uuid"].([]string), data.UUID),
		}, firestore.MergeAll); err != nil {
		http.Error(writer, fmt.Sprintf("firestore.Set: %v", err), http.StatusInternalServerError)
		return
	}
}

// RegisterRecipe writes a new recipe to the recipe collection.
// exported to https://asia-northeast1-superfarmers.cloudfunctions.net/RegisterRecipe
func RegisterRecipe(writer http.ResponseWriter, request *http.Request) {
	var data = new(struct {
		Email  string        `json:"email"`
		UUID   string        `json:"uuid"`
		Recipe shared.Recipe `json:"recipe"`
	})

	if err := json.NewDecoder(request.Body).Decode(data); err != nil {
		http.Error(writer, fmt.Sprintf("json.Decode: %v", err), http.StatusInternalServerError)
		return
	}
	defer request.Body.Close()

	doc, err := client.Collection("farmers").
		Where("email", "==", data.Email).
		Documents(context.Background()).
		Next()

	if err != iterator.Done && err != nil {
		http.Error(writer, fmt.Sprintf("firestore.Next: %v", err), http.StatusInternalServerError)
		return
	}
	if _, err = doc.Ref.Collection("recipe").
		Doc(data.UUID).
		Set(context.Background(), data.Recipe.ToMap()); err != nil { // if there was an existing recipe, overwrite it
		http.Error(writer, fmt.Sprintf("firestore.Set: %v", err), http.StatusInternalServerError)
		return
	}
}

// CheckDeviceOverlap checks duplication of device.
// exported to https://asia-northeast1-superfarmers.cloudfunctions.net/CheckDeviceOverlap
func CheckDeviceOverlap(writer http.ResponseWriter, request *http.Request) {
	var device = new(struct {
		UUID string `json:"uuid"`
	})

	if err := json.NewDecoder(request.Body).Decode(device); err != nil {
		http.Error(writer, fmt.Sprintf("json.Decode: %v", err), http.StatusInternalServerError)
		return
	}
	defer request.Body.Close()

	snapshot, err := client.Collection("farmers").
		Documents(context.Background()).
		GetAll()

	if err != nil {
		http.Error(writer, fmt.Sprintf("firestore.GetAll: %v", err), http.StatusInternalServerError)
		return
	}
	for _, doc := range snapshot {
		var uuids = doc.Data()["device_uuid"].([]string)
		for _, uuid := range uuids {
			if device.UUID == uuid {
				http.Error(writer, fmt.Sprintf("duplicate id: %s", uuid), http.StatusForbidden)
			}
		}
	}
}

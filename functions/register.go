package functions

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"cloud.google.com/go/firestore"
	"github.com/maengsanha/smartfarmer-backend/functions/shared/recipe"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// RegisterDevice appends a new device information to the existing device list.
func RegisterDevice(writer http.ResponseWriter, request *http.Request) {
	var data = new(struct {
		UID  string `json:"uid"`
		UUID string `json:"uuid"`
	})

	if err := json.NewDecoder(request.Body).Decode(data); err != nil {
		http.Error(writer, fmt.Sprintf("json.Decode: %v", err), http.StatusInternalServerError)
		return
	}
	defer request.Body.Close()

	doc, err := client.Collection("farmers").
		Doc(data.UID).
		Get(context.Background())

	if err != nil {
		http.Error(writer, fmt.Sprintf("firestore.Get: %v", err), http.StatusInternalServerError)
		return
	}
	if _, err = client.Collection("devices").
		Doc(data.UUID).
		Set(context.Background(),
			map[string]bool{
				"in_use": true,
			}); err != nil {
		http.Error(writer, fmt.Sprintf("firestore.Set: %v", err), http.StatusInternalServerError)
		return
	}
	if _, err = doc.Ref.Set(context.Background(),
		map[string][]interface{}{
			"device_uuid": append(doc.Data()["device_uuid"].([]interface{}), data.UUID),
		}, firestore.MergeAll); err != nil {
		http.Error(writer, fmt.Sprintf("firestore.Set: %v", err), http.StatusInternalServerError)
		return
	}
}

// RegisterRecipe writes a new recipe to the recipe collection.
// exported to https://asia-northeast1-superfarmers.cloudfunctions.net/RegisterRecipe
func RegisterRecipe(writer http.ResponseWriter, request *http.Request) {
	var data = new(struct {
		UID    string        `json:"uid"`
		UUID   string        `json:"uuid"`
		Recipe recipe.Recipe `json:"recipe"`
	})

	if err := json.NewDecoder(request.Body).Decode(data); err != nil {
		http.Error(writer, fmt.Sprintf("json.Decode: %v", err), http.StatusInternalServerError)
		return
	}
	defer request.Body.Close()

	doc, err := client.Collection("farmers").
		Doc(data.UID).
		Get(context.Background())

	if err != nil {
		http.Error(writer, fmt.Sprintf("firestore.Get: %v", err), http.StatusInternalServerError)
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

	doc, err := client.Collection("devices").
		Doc(device.UUID).
		Get(context.Background())

	if status.Code(err) == codes.NotFound {
		http.Error(writer, fmt.Sprintf("unregistered device: %v", err), http.StatusNotFound)
		return
	}

	if err != nil {
		http.Error(writer, fmt.Sprintf("firestore.Get: %v", err), http.StatusInternalServerError)
		return
	}

	if doc.Data()["in_use"].(bool) {
		http.Error(writer, fmt.Sprintf("duplicated device: %s", device.UUID), http.StatusForbidden)
		return
	}
}

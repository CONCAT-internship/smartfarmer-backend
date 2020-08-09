package functions

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"cloud.google.com/go/firestore"

	"google.golang.org/api/iterator"

	"github.com/joshua-dev/smartfarmer-backend/functions/shared"
)

// Insert stores a sensor data into Firestore.
// exported to https://asia-northeast1-superfarmers.cloudfunctions.net/Insert
func Insert(writer http.ResponseWriter, request *http.Request) {
	var data = new(shared.SensorData)
	if err := json.NewDecoder(request.Body).Decode(data); err != nil {
		http.Error(writer, fmt.Sprintf("json.Decode: %v", err), http.StatusInternalServerError)
		return
	}
	defer request.Body.Close()

	data.SetTime()

	farmer, err := client.Collection("farmers").
		Where("device_uuid", "array-contains", data.UUID).
		Documents(context.Background()).
		Next()

	if err != iterator.Done && err != nil {
		http.Error(writer, fmt.Sprintf("firestore.Next: %v", err), http.StatusInternalServerError)
		return
	}
	doc, err := farmer.Ref.Collection("recipe").
		Doc(data.UUID).
		Get(context.Background())
	if err != nil {
		http.Error(writer, fmt.Sprintf("firestore.Get: %v", err), http.StatusInternalServerError)
		return
	}
	var recipe = shared.Document(doc.Data()).ToRecipe()
	var errorcodes = data.Validate(recipe)
	var requirements = map[string]interface{}{
		"pH_pump": shared.PH_KEEP,
		"ec_pump": shared.EC_KEEP,
	}

	if len(errorcodes) > 0 { // something wrong in the data
		if _, _, err = client.Collection("abnormal").
			Add(context.Background(), map[string]interface{}{
				"uuid":   data.UUID,
				"errors": errorcodes,
			}); err != nil {
			http.Error(writer, fmt.Sprintf("firestore.Add: %v", err), http.StatusInternalServerError)
			return
		}
		for _, code := range errorcodes {
			switch code {
			case shared.CODE_PH_IMPROPER_HIGH:
				requirements["pH_pump"] = shared.PH_DEC
			case shared.CODE_PH_IMPROPER_LOW:
				requirements["pH_pump"] = shared.PH_INC
			case shared.CODE_EC_IMPROPER_HIGH:
				fallthrough
			case shared.CODE_EC_IMPROPER_LOW:
				requirements["ec_pump"] = shared.EC_INC
			case shared.CODE_TEMPERATURE_IMPROPER_HIGH:
				requirements["fan"] = true
			case shared.CODE_TEMPERATURE_IMPROPER_LOW:
				requirements["fan"] = false
			case shared.CODE_HUMIDITY_IMPROPER_HIGH:
				fallthrough
			case shared.CODE_HUMIDITY_IMPROPER_LOW:
				fallthrough
			default:
			}
		}
		if _, err = client.Collection("desired_status").
			Doc(data.UUID).
			Set(context.Background(), requirements, firestore.MergeAll); err != nil {
			http.Error(writer, fmt.Sprintf("firestore.Set: %v", err), http.StatusInternalServerError)
			return
		}
	}

	if _, _, err = client.Collection("sensor_data").
		Add(context.Background(), data.ToMap()); err != nil {
		http.Error(writer, fmt.Sprintf("firestore.Add: %v", err), http.StatusInternalServerError)
		return
	}
}

package functions

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/maengsanha/smartfarmer-backend/shared/recipe"
	"github.com/maengsanha/smartfarmer-backend/shared/sensor"
	"google.golang.org/api/iterator"
)

// Insert stores a sensor data into Firestore.
// exported to https://asia-northeast1-superfarmers.cloudfunctions.net/Insert
func Insert(writer http.ResponseWriter, request *http.Request) {
	var data = new(sensor.Data)
	if err := json.NewDecoder(request.Body).Decode(data); err != nil {
		http.Error(writer, fmt.Sprintf("json.Decode: %v", err), http.StatusInternalServerError)
		return
	}
	defer request.Body.Close()

	data.SetTime()

	// find recipe
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
	var recipe = new(recipe.Recipe)
	recipe.FromMap(doc.Data())

	lastDoc, err := client.Collection("sensor_data").
		Where("uuid", "==", data.UUID).
		OrderBy("unix_time", firestore.Desc).
		Limit(1).
		Documents(context.Background()).
		Next()

	if err != iterator.Done && err != nil {
		http.Error(writer, fmt.Sprintf("firestore.Next: %v", err), http.StatusInternalServerError)
		return
	}
	var lastData = new(sensor.Data)
	lastData.FromMap(lastDoc.Data())

	if data.LED { // if LED is on, increase daylight time and compare with the recipe
		data.LightTime = lastData.LightTime + sensor.TRANSMISSION_CYCLE
	} else {
		data.DarkTime = lastData.DarkTime + sensor.TRANSMISSION_CYCLE
	}
	var LEDoption bool
	if data.LightTime >= recipe.LightTime*60 { // if daylight time is exceeded, request to LED off
		LEDoption = false
		data.LightTime = 0
	}
	if data.DarkTime >= (24-recipe.LightTime)*60 { // if dark time is exceeded, request to LED on
		LEDoption = true
		data.DarkTime = 0
	}
	if _, err = client.Collection("desired_status").
		Doc(data.UUID).
		Set(context.Background(),
			map[string]interface{}{
				"led": LEDoption,
			}, firestore.MergeAll); err != nil {
		http.Error(writer, fmt.Sprintf("firestore.Set: %v", err), http.StatusInternalServerError)
		return
	}
	errorCodes, requirements := data.Validate(*recipe)

	if len(errorCodes) > 0 { // something wrong in the data
		if _, _, err = client.Collection("abnormal").
			Add(context.Background(),
				map[string]interface{}{
					"uuid":   data.UUID,
					"errors": errorCodes,
					"time":   time.Now().Unix(),
				}); err != nil {
			http.Error(writer, fmt.Sprintf("firestore.Add: %v", err), http.StatusInternalServerError)
			return
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

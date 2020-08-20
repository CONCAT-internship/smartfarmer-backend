package functions

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"net/http"

	"cloud.google.com/go/firestore"
	"github.com/maengsanha/smartfarmer-backend/functions/shared/sensor"
	"google.golang.org/api/iterator"
)

// RecentStatus returns the latest status of the farm.
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
	var data = new(sensor.Data)
	data.FromMap(doc.Data())

	writer.Header().Set("Content-Type", "application/json")
	if err = json.NewEncoder(writer).Encode(map[string]interface{}{
		"temperature":        round2SecondDecimal(data.Temperature),
		"humidity":           round2SecondDecimal(data.Humidity),
		"pH":                 round2SecondDecimal(data.PH),
		"ec":                 round2SecondDecimal(data.EC),
		"light":              round2SecondDecimal(data.Light),
		"liquid_temperature": round2SecondDecimal(data.LiquidTemperature),
		"liquid_level":       data.LiquidLevel,
		"led":                data.LED,
		"fan":                data.Fan,
		"unix_time":          data.UnixTime,
		"local_time":         data.LocalTime,
	}); err != nil {
		http.Error(writer, fmt.Sprintf("json.Encode: %v", err), http.StatusInternalServerError)
		return
	}
}

func round2SecondDecimal(data float64) float64 {
	return math.Round(data*10) / 10
}

// DesiredStatus returns the desired status of the device.
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

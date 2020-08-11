package functions

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/joshua-dev/smartfarmer-backend/functions/shared"
)

// Monitor checks whether the devices are working properly.
// exported to https://asia-northeast1-superfarmers.cloudfunctions.net/Monitor
func Monitor(writer http.ResponseWriter, request *http.Request) {
	defer request.Body.Close()

	var base = time.Now().Unix() - 2*shared.TRANSMISSION_CYCLE*60

	lastRecords, err := client.Collection("sensor_data").
		Where("unix_time", ">=", base).
		Documents(context.Background()).
		GetAll()

	if err != nil {
		http.Error(writer, fmt.Sprintf("firestore.GetAll: %v", err), http.StatusInternalServerError)
		return
	}
	farmers, err := client.Collection("farmers").
		Documents(context.Background()).
		GetAll()

	if err != nil {
		http.Error(writer, fmt.Sprintf("firestore.GetAll: %v", err), http.StatusInternalServerError)
		return
	}
	var recordedDevices map[string]struct{}

	for _, record := range lastRecords {
		recordedDevices[record.Data()["uuid"].(string)] = struct{}{}
	}
	var devices []string

	for _, doc := range farmers {
		var device_uuids = doc.Data()["device_uuid"].([]interface{})
		for _, device_uuid := range device_uuids {
			devices = append(devices, device_uuid.(string))
		}
	}
	for _, device := range devices {
		if _, exists := recordedDevices[device]; !exists {
			if _, _, err = client.Collection("abnormal").
				Add(context.Background(),
					map[string]interface{}{
						"errors": []shared.ErrorCode{shared.CODE_DATA_EMPTY},
						"time":   time.Now().Unix(),
						"uuid":   device,
					}); err != nil {
				http.Error(writer, fmt.Sprintf("firestore.Add: %v", err), http.StatusInternalServerError)
				return
			}
		}
	}
}

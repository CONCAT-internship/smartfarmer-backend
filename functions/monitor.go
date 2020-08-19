package functions

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/maengsanha/smartfarmer-backend/functions/shared"
)

// Monitor checks whether the devices are working properly.
// exported to https://asia-northeast1-superfarmers.cloudfunctions.net/Monitor
func Monitor(writer http.ResponseWriter, request *http.Request) {
	defer request.Body.Close()

	var base = time.Now().Unix() - 2*shared.TRANSMISSION_CYCLE*60

	history, err := client.Collection("sensor_data").
		Where("unix_time", ">=", base).
		Documents(context.Background()).
		GetAll()

	if err != nil {
		http.Error(writer, fmt.Sprintf("firestore.GetAll: %v", err), http.StatusInternalServerError)
		return
	}

	var idSet = make(map[string]struct{})
	for _, doc := range history {
		idSet[doc.Data()["uuid"].(string)] = struct{}{}
	}

	inUse, err := client.Collection("devices").
		Where("in_use", "==", true).
		Documents(context.Background()).
		GetAll()

	if err != nil {
		http.Error(writer, fmt.Sprintf("firestore.GetAll: %v", err), http.StatusInternalServerError)
		return
	}

	for _, device := range inUse {
		if _, exists := idSet[device.Ref.ID]; !exists {
			if _, _, err = client.Collection("abnormal").
				Add(context.Background(),
					map[string]interface{}{
						"errors": []shared.ErrorCode{shared.CODE_DATA_EMPTY},
						"time":   time.Now().Unix(),
						"uuid":   device.Ref.ID,
					}); err != nil {
				http.Error(writer, fmt.Sprintf("firestore.Add: %v", err), http.StatusInternalServerError)
				return
			}
		}
	}
}

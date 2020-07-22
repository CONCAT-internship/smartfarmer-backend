// Package functions defines sensor data model and implements its CRUD operations.
package functions

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"cloud.google.com/go/firestore"
)

// Monitor checks the records for twice the duration of the insertion cycle of the device with matching uuid.
func Monitor(writer http.ResponseWriter, req *http.Request) {
	ctx := context.Background()

	client, err := firestore.NewClient(ctx, projectID)
	if err != nil {
		fmt.Fprintf(writer, "firestore.NewClient: %v\n", err)
		return
	}
	defer client.Close()
	defer req.Body.Close()

	const cycle = 3
	base := time.Now().Unix() - 2*cycle*60

	records, err := client.Collection("sensor_data").Where("unix_time", ">=", base).Documents(ctx).GetAll()
	if err != nil {
		fmt.Fprintf(writer, "firestore.GetAll: %v", err)
		return
	}

	for _, record := range records {
		data := Document(record.Data()).toStruct()
		// TODO: change validate to something checking propriety of data
		if err = data.validate(); err != nil {
			if _, _, fail := client.Collection("abnormal").Add(ctx, data.toLog(err)); fail != nil {
				fmt.Fprintf(writer, "firestore.Add: %v", err)
				return
			}
		}
	}
}

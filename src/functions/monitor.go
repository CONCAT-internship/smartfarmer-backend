// Package functions defines sensor data model and implements related operations.
package functions

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"google.golang.org/api/iterator"

	"cloud.google.com/go/firestore"
)

// Monitor checks the records for twice the duration of the insertion cycle of sensor data.
func Monitor(writer http.ResponseWriter, req *http.Request) {
	ctx := context.Background()

	client, err := firestore.NewClient(ctx, PROJECT_ID)
	if err != nil {
		fmt.Fprintf(writer, "firestore.NewClient: %v", err)
		return
	}
	defer client.Close()
	defer req.Body.Close()

	const cycle = 3
	base := time.Now().Unix() - 2*cycle*60

	cursor := client.Collection("sensor_data").Where("unix_time", ">=", base).Documents(ctx)
	for {
		doc, err := cursor.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			fmt.Fprintf(writer, "firestore.Next: %v", err)
			return
		}
		data := document(doc.Data()).toStruct()
		if err = data.validate(); err != nil {
			client.Collection("abnormal").Add(ctx, map[string]interface{}{
				"uuid":    data.UUID,
				"message": err.Error(),
			})
		}
		if err = data.appropriate(); err != nil {
			client.Collection("abnormal").Add(ctx, map[string]interface{}{
				"uuid":    data.UUID,
				"message": err.Error(),
			})
		}
	}
}

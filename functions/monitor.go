package functions

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/joshua-dev/smartfarmer-backend/functions/shared"
	"google.golang.org/api/iterator"
)

// Monitor checks the records for twice the duration of the insertion cycle of sensor data.
func Monitor(writer http.ResponseWriter, request *http.Request) {
	defer request.Body.Close()

	const cycle = 3
	var base = time.Now().Unix() - 2*cycle*60

	var cursor = client.Collection("sensor_data").
		Where("unix_time", ">=", base).
		Documents(context.Background())

	for {
		doc, err := cursor.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			http.Error(writer, fmt.Sprintf("firestore.Next: %v", err), http.StatusInternalServerError)
			return
		}
		var data = shared.Document(doc.Data()).ToStruct()
		if err = data.Validate(); err != nil {
			client.Collection("abnormal").Add(context.Background(), shared.Document{
				"uuid":    data.UUID,
				"message": err.Error(),
			})
		}
		if err = data.Appropriate(); err != nil {
			client.Collection("abnormal").Add(context.Background(), shared.Document{
				"uuid":    data.UUID,
				"message": err.Error(),
			})
		}
	}
}
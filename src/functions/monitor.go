// Package functions defines sensor data model and implements related operations.
package functions

// Monitor checks the records for twice the duration of the insertion cycle of sensor data.
// func Monitor(writer http.ResponseWriter, req *http.Request) {
// 	ctx := context.Background()

// 	client, err := firestore.NewClient(ctx, projectID)
// 	if err != nil {
// 		fmt.Fprintf(writer, "firestore.NewClient: %v", err)
// 		return
// 	}
// 	defer client.Close()
// 	defer req.Body.Close()

// 	const cycle = 3
// 	base := time.Now().Unix() - 2*cycle*60

// 	cursor := client.Collection("sensor_data").Where("unix_time", ">=", base).Documents(ctx)

// }

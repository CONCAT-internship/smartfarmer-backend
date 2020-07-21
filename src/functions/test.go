/*
Package functions defines sensor data model
and implements CRUD operations of smart farm.
*/
package functions

// [Start fs_functions_dependencies]
import (
	"context"
	"fmt"
	"net/http"

	"cloud.google.com/go/firestore"
)

// [End fs_functions_dependencies]

// [Start fs_functions_test]

/*
Test is just a test function to test a single query.
exported to https://asia-northeast1-superfarmers.cloudfunctions.net/Test
*/
func Test(writer http.ResponseWriter, req *http.Request) {
	ctx := context.Background()
	client, err := firestore.NewClient(ctx, projectID)
	if err != nil {
		fmt.Fprintf(writer, "firestore.NewClient: %v\n", err)
		return
	}
	defer client.Close()
	defer req.Body.Close()

	snapshots, err := client.Collection("sensor_data").Where("uuid", "==", "test").OrderBy("unix_time", firestore.Desc).Documents(ctx).GetAll()
	if err != nil {
		fmt.Fprintf(writer, "firestore.GetAll: %v\n", err)
	}
	for _, snapshot := range snapshots {
		fmt.Fprintln(writer, snapshot.Data()["unix_time"])
	}
}

// [End fs_functions_insert]

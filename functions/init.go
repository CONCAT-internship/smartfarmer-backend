package functions

import (
	"context"
	"log"

	"cloud.google.com/go/firestore"
)

// db is a global Firestore client, initialized once per instance.
var db *firestore.Client

func init() {
	// err is pre-declared to avoid shadowing client.
	var err error

	// db is initialized with context.Background() because it should
	// persist between function invocations.
	db, err = firestore.NewClient(context.Background(), "superfarmers")
	if err != nil {
		log.Fatalf("firestore.NewClient: %v", err)
	}
}

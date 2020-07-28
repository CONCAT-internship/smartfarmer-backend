package functions

import (
	"context"
	"log"

	"cloud.google.com/go/firestore"
)

// client is a global Firestore client, initialized once per instance.
var client *firestore.Client

func init() {
	// err is pre-declared to avoid shadowing client.
	var err error

	// client is initialized with context.Background() because it should
	// persist between function invocations.
	client, err = firestore.NewClient(context.Background(), "superfarmers")
	if err != nil {
		log.Fatalf("firestore.NewClient: %v", err)
	}
}

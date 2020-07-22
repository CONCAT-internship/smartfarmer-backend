// Package functions defines sensor data model and implements its CRUD operations.
package functions

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func ExampleInsert() {
	sample := SensorData{
		UUID:              "test",
		Temperature:       33.3,
		Humidity:          76.9,
		PH:                4.5,
		EC:                3.0,
		Light:             67.8,
		LiquidTemperature: 35.0,
		LiquidFlowRate:    2.1,
		LiquidLevel:       false,
		Valve:             true,
		LED:               true,
		Fan:               false,
	}
	obj, err := json.Marshal(sample)
	if err != nil {
		fmt.Printf("json.Marshal: %v\n", err)
	}
	buff := bytes.NewBuffer(obj)
	resp, err := http.Post("https://asia-northeast1-superfarmers.cloudfunctions.net/Insert", "application/json", buff)
	if err != nil {
		fmt.Printf("http.Post: %v\n", err)
	}
	defer resp.Body.Close()
	if msg, err := ioutil.ReadAll(resp.Body); err != nil {
		fmt.Printf("ioutil.ReadAll: %v\n", err)
	} else {
		fmt.Println(string(msg))
	}
	// Output:
	// Successfully stored to Firestore.
}

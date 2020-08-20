package functions

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"cloud.google.com/go/firestore"
)

// LookupByPeriod returns records of last period seconds.
func LookupByPeriod(writer http.ResponseWriter, request *http.Request) {
	var data = new(struct {
		UUID   string `json:"uuid"`
		Period int64  `json:"period"`
		Key    string `json:"key"`
	})
	if err := json.NewDecoder(request.Body).Decode(data); err != nil {
		http.Error(writer, fmt.Sprintf("strconv.Atoi: %v", err), http.StatusBadRequest)
		return
	}
	defer request.Body.Close()

	docs, err := client.Collection("sensor_data").
		Where("uuid", "==", data.UUID).
		OrderBy("unix_time", firestore.Desc).
		Where("unix_time", ">=", time.Now().Unix()-data.Period).
		Documents(context.Background()).
		GetAll()

	if err != nil {
		http.Error(writer, fmt.Sprintf("firestore.GetAll: %v", err), http.StatusInternalServerError)
		return
	}
	records := make([]map[string]interface{}, len(docs))

	if len(data.Key) > 0 {
		for idx := range docs {
			records[idx] = map[string]interface{}{
				data.Key:     docs[len(docs)-1-idx].Data()[data.Key],
				"local_time": docs[len(docs)-1-idx].Data()["local_time"],
			}
		}
	} else {
		for idx := range docs {
			records[idx] = map[string]interface{}{
				"data":       docs[len(docs)-1-idx].Data(),
				"local_time": docs[len(docs)-1-idx].Data()["local_time"],
			}
		}
	}

	writer.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(writer).Encode(records); err != nil {
		http.Error(writer, fmt.Sprintf("json.Encode: %v", err), http.StatusInternalServerError)
		return
	}
}

// LookupByNumber returns a given number of recent records.
func LookupByNumber(writer http.ResponseWriter, request *http.Request) {
	var data = new(struct {
		UUID   string `json:"uuid"`
		Number int    `json:"number"`
		Key    string `json:"key"`
	})
	if err := json.NewDecoder(request.Body).Decode(data); err != nil {
		http.Error(writer, fmt.Sprintf("json.Decode: %v", err), http.StatusInternalServerError)
		return
	}
	defer request.Body.Close()

	docs, err := client.Collection("sensor_data").
		Where("uuid", "==", data.UUID).
		OrderBy("unix_time", firestore.Desc).
		Limit(data.Number).
		Documents(context.Background()).
		GetAll()

	if err != nil {
		http.Error(writer, fmt.Sprintf("firestore.GetAll: %v", err), http.StatusInternalServerError)
		return
	}
	records := make([]map[string]interface{}, len(docs))

	if len(data.Key) > 0 {
		for idx := range docs {
			records[idx] = map[string]interface{}{
				data.Key:     docs[len(docs)-1-idx].Data()[data.Key],
				"local_time": docs[len(docs)-1-idx].Data()["local_time"],
			}
		}
	} else {
		for idx := range docs {
			records[idx] = map[string]interface{}{
				"data":       docs[len(docs)-1-idx].Data(),
				"local_time": docs[len(docs)-1-idx].Data()["local_time"],
			}
		}
	}

	writer.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(writer).Encode(records); err != nil {
		http.Error(writer, fmt.Sprintf("json.Encode: %v", err), http.StatusInternalServerError)
		return
	}
}

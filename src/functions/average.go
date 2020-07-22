// Package functions defines sensor data model and implements its CRUD operations.
package functions

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"cloud.google.com/go/firestore"
)

// DailyAverage calculates daily average of sensor data with given uuid and unix time.
// exported to https://asia-northeast1-superfarmers.cloudfunctions.net/DailyAverage
func DailyAverage(writer http.ResponseWriter, req *http.Request) {
	ctx := context.Background()

	client, err := firestore.NewClient(ctx, projectID)
	if err != nil {
		fmt.Fprintf(writer, "firestore.NewClient: %v", err)
		return
	}
	defer client.Close()

	uuid := req.URL.Query().Get("uuid")
	unixtime, err := strconv.Atoi(req.URL.Query().Get("unixtime"))
	if err != nil {
		fmt.Fprintf(writer, "strconv.Atoi: %v", err)
		return
	}
	defer req.Body.Close()

	var resp struct {
		Average struct {
			Temperature       float64 `json:"temperature"`
			Humidity          float64 `json:"humidity"`
			PH                float64 `json:"pH"`
			EC                float64 `json:"ec"`
			Light             float64 `json:"light"`
			LiquidTemperature float64 `json:"liquid_temperature"`
			LiquidFlowRate    float64 `json:"liquid_flow_rate"`
		} `json:"average"`
	}

	records, err := client.Collection("sensor_data").Where("uuid", "==", uuid).OrderBy("unix_time", firestore.Desc).Where("unix_time", ">=", unixtime).Where("unix_time", "<", unixtime+24*60*60).Documents(ctx).GetAll()
	if err != nil {
		fmt.Fprintf(writer, "firestore.GetAll: %v", err)
		return
	}

	// prevent zero division error
	if len(records) > 0 {
		for _, record := range records {
			data := Document(record.Data()).toStruct()
			resp.Average.Temperature += data.Temperature
			resp.Average.Humidity += data.Humidity
			resp.Average.PH += data.PH
			resp.Average.EC += data.EC
			resp.Average.Light += data.Light
			resp.Average.LiquidTemperature += data.LiquidTemperature
			resp.Average.LiquidFlowRate += data.LiquidFlowRate
		}
		resp.Average.Temperature /= float64(len(records))
		resp.Average.Humidity /= float64(len(records))
		resp.Average.PH /= float64(len(records))
		resp.Average.EC /= float64(len(records))
		resp.Average.Light /= float64(len(records))
		resp.Average.LiquidTemperature /= float64(len(records))
		resp.Average.LiquidFlowRate /= float64(len(records))
	}

	writer.Header().Set("Content-Type", "application/json")
	if err = json.NewEncoder(writer).Encode(resp); err != nil {
		fmt.Fprintf(writer, "json.Encode: %v", err)
		return
	}
}

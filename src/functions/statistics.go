// Package functions defines sensor data model and implements related operations.
package functions

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"google.golang.org/api/iterator"

	"cloud.google.com/go/firestore"
)

// DailyAverage calculates the average of the sensor data for each day of the week.
// exported to https://asia-northeast1-superfarmers.cloudfunctions.net/DailyAverage
func DailyAverage(writer http.ResponseWriter, req *http.Request) {
	ctx := context.Background()

	client, err := firestore.NewClient(ctx, PROJECT_ID)
	if err != nil {
		fmt.Fprintf(writer, "firestore.NewClient: %v", err)
		return
	}
	defer client.Close()

	uuid := req.URL.Query().Get("uuid")
	base, err := strconv.Atoi(req.URL.Query().Get("unixtime"))
	if err != nil {
		fmt.Fprintf(writer, "strconv.Atoi: %v", err)
		return
	}
	defer req.Body.Close()

	const day_time = 24 * 60 * 60

	datas := make([][]sensorData, 7)
	tmp := new(sensorData)

	cursor := client.Collection("sensor_data").Where("uuid", "==", uuid).OrderBy("unix_time", firestore.Asc).Where("unix_time", ">=", base).Where("unix_time", "<", base+7*day_time).Documents(ctx)
	for {
		doc, err := cursor.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			fmt.Fprintf(writer, "firestore.Next: %v", err)
			return
		}
		doc.DataTo(tmp)
		idx := (int(tmp.UnixTime) - base) / day_time
		datas[idx] = append(datas[idx], *tmp)
	}

	writer.Header().Set("Content-Type", "application/json")
	if err = json.NewEncoder(writer).Encode(map[string]map[string]float64{
		"sun": average(datas[0]),
		"mon": average(datas[1]),
		"tue": average(datas[2]),
		"wed": average(datas[3]),
		"thu": average(datas[4]),
		"fri": average(datas[5]),
		"sat": average(datas[6]),
	}); err != nil {
		fmt.Fprintf(writer, "json.Encode: %v", err)
		return
	}
}

func average(datas []sensorData) map[string]float64 {
	avg := make(map[string]float64)

	for _, data := range datas {
		avg["temperature"] += data.Temperature
		avg["humidity"] += data.Humidity
		avg["pH"] += data.PH
		avg["ec"] += data.EC
		avg["light"] += data.Light
		avg["liquid_temperature"] += data.LiquidTemperature
		avg["liquid_flow_rate"] += data.LiquidFlowRate
	}
	avg["temperature"] /= float64(len(datas))
	avg["humidity"] /= float64(len(datas))
	avg["pH"] /= float64(len(datas))
	avg["ec"] /= float64(len(datas))
	avg["light"] /= float64(len(datas))
	avg["liquid_temperature"] /= float64(len(datas))
	avg["liquid_flow_rate"] /= float64(len(datas))

	return avg
}

package functions

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"strconv"

	"cloud.google.com/go/firestore"
	"github.com/joshua-dev/smartfarmer-backend/functions/shared"
	"google.golang.org/api/iterator"
)

// DailyAverage calculates the average of the sensor data for each day of the week.
// exported to https://asia-northeast1-superfarmers.cloudfunctions.net/DailyAverage
func DailyAverage(writer http.ResponseWriter, request *http.Request) {
	var uuid = request.URL.Query().Get("uuid")
	base, err := strconv.Atoi(request.URL.Query().Get("unixtime"))
	if err != nil {
		http.Error(writer, fmt.Sprintf("strconv.Atoi: %v", err), http.StatusBadRequest)
		return
	}
	defer request.Body.Close()

	const day_time = 24 * 60 * 60

	var datas = make([][]shared.SensorData, 7)

	var cursor = db.Collection("sensor_data").
		Where("uuid", "==", uuid).
		OrderBy("unix_time", firestore.Desc).
		Where("unix_time", ">=", base).
		Where("unix_time", "<", base+7*day_time).
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
		var idx = (int(data.UnixTime) - base) / day_time
		datas[idx] = append(datas[idx], data)
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
		http.Error(writer, fmt.Sprintf("json.Encode: %v", err), http.StatusInternalServerError)
		return
	}
}

func average(datas []shared.SensorData) map[string]float64 {
	var avg = make(map[string]float64)
	if len(datas) > 0 {
		for _, data := range datas {
			avg["temperature"] += data.Temperature
			avg["humidity"] += data.Humidity
			avg["pH"] += data.PH
			avg["ec"] += data.EC
			avg["light"] += data.Light
			avg["liquid_temperature"] += data.LiquidTemperature
			avg["liquid_flow_rate"] += data.LiquidFlowRate
		}
		avg["temperature"] = round2SecondDecimal(avg["temperature"] / float64(len(datas)))
		avg["humidity"] = round2SecondDecimal(avg["humidity"] / float64(len(datas)))
		avg["pH"] = round2SecondDecimal(avg["pH"] / float64(len(datas)))
		avg["ec"] = round2SecondDecimal(avg["ec"] / float64(len(datas)))
		avg["light"] = round2SecondDecimal(avg["light"] / float64(len(datas)))
		avg["liquid_temperature"] = round2SecondDecimal(avg["liquid_temperature"] / float64(len(datas)))
		avg["liquid_flow_rate"] = round2SecondDecimal(avg["liquid_flow_rate"] / float64(len(datas)))
	}
	return avg
}

func round2SecondDecimal(data float64) float64 {
	return math.Round(data*10) / 10
}

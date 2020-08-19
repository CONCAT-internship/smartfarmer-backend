package functions

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"strconv"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/maengsanha/smartfarmer-backend/shared/sensor"
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

	var datas = make([][]sensor.Data, 7)

	var cursor = client.Collection("sensor_data").
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
		var data = new(sensor.Data)
		data.FromMap(doc.Data())
		var idx = (int(data.UnixTime) - base) / day_time
		datas[idx] = append(datas[idx], *data)
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

func average(datas []sensor.Data) map[string]float64 {
	var avg = make(map[string]float64)
	if len(datas) > 0 {
		for _, data := range datas {
			avg["temperature"] += data.Temperature
			avg["humidity"] += data.Humidity
			avg["pH"] += data.PH
			avg["ec"] += data.EC
			avg["light"] += data.Light
			avg["liquid_temperature"] += data.LiquidTemperature
		}
		avg["temperature"] = round2SecondDecimal(avg["temperature"] / float64(len(datas)))
		avg["humidity"] = round2SecondDecimal(avg["humidity"] / float64(len(datas)))
		avg["pH"] = round2SecondDecimal(avg["pH"] / float64(len(datas)))
		avg["ec"] = round2SecondDecimal(avg["ec"] / float64(len(datas)))
		avg["light"] = round2SecondDecimal(avg["light"] / float64(len(datas)))
		avg["liquid_temperature"] = round2SecondDecimal(avg["liquid_temperature"] / float64(len(datas)))
	}
	return avg
}

func round2SecondDecimal(data float64) float64 {
	return math.Round(data*10) / 10
}

// Records returns records of last period seconds.
// exported to https://asia-northeast1-superfarmers.cloudfunctions.net/Records
func Records(writer http.ResponseWriter, request *http.Request) {
	var data = new(struct {
		UUID string `json:"uuid"`
		Time int64  `json:"time"`
		Key  string `json:"key"`
	})
	if err := json.NewDecoder(request.Body).Decode(data); err != nil {
		http.Error(writer, fmt.Sprintf("strconv.Atoi: %v", err), http.StatusBadRequest)
		return
	}
	defer request.Body.Close()

	var records []map[string]interface{}

	var cursor = client.Collection("sensor_data").
		Where("uuid", "==", data.UUID).
		OrderBy("unix_time", firestore.Asc).
		Where("unix_time", ">=", time.Now().Unix()-data.Time).
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
		var record = map[string]interface{}{
			data.Key:     doc.Data()[data.Key],
			"local_time": doc.Data()["local_time"],
		}
		records = append(records, record)
	}
	writer.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(writer).Encode(records); err != nil {
		http.Error(writer, fmt.Sprintf("json.Encode: %v", err), http.StatusInternalServerError)
		return
	}
}

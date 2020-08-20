package functions

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// ProfileFarmer returns farmer's profile.
func ProfileFarmer(writer http.ResponseWriter, request *http.Request) {
	uid := request.URL.Query().Get("uid")
	defer request.Body.Close()

	doc, err := client.Collection("farmers").
		Doc(uid).
		Get(context.Background())

	if err != nil {
		http.Error(writer, fmt.Sprintf("firestore.Get: %v", err), http.StatusInternalServerError)
		return
	}
	farmer_info := new(struct {
		Nickname string              `json:"nickname"`
		FarmInfo []map[string]string `json:"farm_info"`
	})

	farmer_info.Nickname = doc.Data()["nickname"].(string)
	recipes, err := doc.Ref.Collection("recipe").
		Documents(context.Background()).
		GetAll()

	if err != nil {
		http.Error(writer, fmt.Sprintf("firestore.GetAll: %v", err), http.StatusInternalServerError)
		return
	}
	for _, recipe := range recipes {
		farmer_info.FarmInfo = append(farmer_info.FarmInfo, map[string]string{
			"device_uuid": recipe.Ref.ID,
			"farm_name":   recipe.Data()["farm_name"].(string),
		})
	}
	writer.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(writer).Encode(farmer_info); err != nil {
		http.Error(writer, fmt.Sprintf("json.Encode: %v", err), http.StatusInternalServerError)
		return
	}
}

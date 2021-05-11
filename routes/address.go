package routes

import (
	"encoding/json"
	"fmt"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
)

func GetCities(w http.ResponseWriter, r *http.Request) {
	mdb := getDB(r)
	cursor, err := mdb.CitiesCollection().Find(r.Context(), bson.D{})
	if err != nil {
		fmt.Printf("Error on cities: %v\n", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	var results []map[string]interface{}
	err = cursor.All(r.Context(), &results)

	if err != nil {
		fmt.Printf("Error on cities: %v\n", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(results)
}

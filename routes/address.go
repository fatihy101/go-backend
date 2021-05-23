package routes

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"enstrurent.com/server/db"
	"github.com/go-chi/chi/v5"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func getCities(w http.ResponseWriter, r *http.Request) {
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

func addAddress(w http.ResponseWriter, r *http.Request) {
	var address db.Address
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	err := dec.Decode(&address)

	if err != nil {
		fmt.Println("Error on adding address: " + err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	address.CreatedAt = time.Now()
	collection := getDB(r).AddressCollection()
	result, err := collection.InsertOne(r.Context(), address)

	if err != nil {
		fmt.Println("Error on adding address: " + err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	json.NewEncoder(w).Encode(map[string]interface{}{"address_id": result.InsertedID})

}

func getAddress(w http.ResponseWriter, r *http.Request) {
	id := strings.Split(r.RequestURI, "/")[2] // Note: chi.URLParam does not work properly. So I handled manually
	hexID, _ := primitive.ObjectIDFromHex(id)
	result := getDB(r).AddressCollection().FindOne(r.Context(), bson.M{"_id": hexID})

	if result.Err() != nil {
		fmt.Printf("Error on getting address : %v\nID: %v\n", result.Err().Error(), hexID)
		http.Error(w, "there's no address with this ID", http.StatusBadRequest)
		return
	}
	var address db.Address
	result.Decode(&address)

	json.NewEncoder(w).Encode(address)
}

func deleteAddress(w http.ResponseWriter, r *http.Request) {
	idHex, err := primitive.ObjectIDFromHex(chi.URLParam(r, "address_id"))
	if err != nil {
		http.Error(w, "something went wrong", http.StatusBadRequest)
		return
	}
	_, err = getDB(r).AddressCollection().DeleteOne(r.Context(), bson.M{"_id": idHex})

	if err != nil {
		fmt.Println("Error on deleting address: " + err.Error())
		http.Error(w, "something went wrong", http.StatusBadRequest)
		return
	}
}

func updateAddress(w http.ResponseWriter, r *http.Request) {
	var newAddress db.Address
	json.NewDecoder(r.Body).Decode(&newAddress)
	fmt.Println(newAddress)
	newAddress.UpdatedAt = time.Now()
	res := getDB(r).AddressCollection().FindOneAndReplace(r.Context(), bson.M{"_id": newAddress.ID}, newAddress)

	if res.Err() != nil {
		http.Error(w, res.Err().Error(), http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(newAddress)
}

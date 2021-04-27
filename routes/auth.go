package routes

import (
	"encoding/json"
	"fmt"
	"net/http"

	"enstrurent.com/server/db"
)

const (
	ClientRole = "client"
	RenterRole = "renter"
)

type LoginInfo struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func login(w http.ResponseWriter, r *http.Request) {
	var info LoginInfo
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	if err := dec.Decode(&info); err != nil { // Decode to info
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// TODO Get role from context
}

func signUpClient(w http.ResponseWriter, r *http.Request) {
	// TODO check is client registered before.
	var clientInfo db.Client

	dec := json.NewDecoder(r.Body)
	dec.Decode(&clientInfo)

	db := getDB(r)
	savedID, err := db.Collection(ClientCollection).InsertOne(r.Context(), clientInfo)

	if err != nil {
		http.Error(w, "Error while signing up", http.StatusBadRequest)
	}
	fmt.Print(savedID.InsertedID)
}

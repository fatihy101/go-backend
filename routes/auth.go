package routes

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"enstrurent.com/server/db"
	"go.mongodb.org/mongo-driver/bson"
)

const (
	ClientRole   = "client"
	RenterRole   = "renter"
	ExpiresHours = time.Hour * 96 // 4 days
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

	getDB(r)

	// TODO Get role from context
}

func signUpClient(w http.ResponseWriter, r *http.Request) {
	// Decode request info.
	var body map[string]interface{}
	json.NewDecoder(r.Body).Decode(&body)
	email := fmt.Sprint(body["email"])
	password := fmt.Sprint(body["password"])

	mdb := getDB(r)
	// Check is user registered before or not
	var qResult bson.M

	mdb.GetCredsByEmail(r.Context(), email, &qResult)

	if len(qResult) != 0 {
		http.Error(w, "User already registered with this email!", http.StatusConflict)
		return
	}

	// Set credentials of user
	var currentRole string
	if body["store_name"] == nil {
		currentRole = ClientRole
	} else {
		currentRole = RenterRole
	}

	role := db.UserCredentials{
		Email:    email,
		Password: HashPassword(password),
		Role:     currentRole,
	}
	// Parse data to json again.
	jsonData, err := json.Marshal(body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	// create decoder
	dec := json.NewDecoder(bytes.NewReader(jsonData))
	// Save user
	if currentRole == ClientRole {
		var clientInfo db.Client
		dec.Decode(&clientInfo) // decode to struct
		mdb.SaveOne(db.UserCredsCollection, r.Context(), role)
		mdb.SaveOne(db.ClientCollection, r.Context(), clientInfo)
	} else {
		var renterInfo db.Renter
		dec.Decode(&renterInfo)
		mdb.SaveOne(db.UserCredsCollection, r.Context(), role)
		mdb.SaveOne(db.RenterCollection, r.Context(), renterInfo)
	}

	// Return jwt
	token, err := generateToken(email, currentRole, time.Now().Add(ExpiresHours).Unix())

	if err != nil {
		http.Error(w, "Something went wrong with authentication", http.StatusBadRequest)
		return
	}

	// Encode with json
	json.NewEncoder(w).Encode(map[string]string{"token": token})
}

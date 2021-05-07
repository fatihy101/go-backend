package routes

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"enstrurent.com/server/db"
)

const (
	ClientRole   = "client"
	RenterRole   = "renter"
	ExpiresHours = time.Hour * 96 // 4 days
)

type AuthResponse struct {
	Token string `json:"token"`
}

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
	mdb := getDB(r)

	creds := mdb.GetCredsByEmail(r.Context(), info.Email)

	if creds.Email == "" {
		http.Error(w, "email is not registered", http.StatusUnauthorized)
	}

	if CompareHashAndPassword(creds.Password, info.Password) {
		createResponse(creds, w, r.Context(), mdb)
	} else {
		http.Error(w, "password is not valid", http.StatusUnauthorized)
	}
}

// ADD created date to sign up
func signUp(w http.ResponseWriter, r *http.Request) {
	// Decode request info.
	var email, password string
	var body map[string]interface{}
	json.NewDecoder(r.Body).Decode(&body)
	if body["email"] == nil || body["password"] == nil {
		http.Error(w, "Email or password is null!", http.StatusBadRequest)
		return
	}
	email = fmt.Sprint(body["email"])
	password = fmt.Sprint(body["password"])

	mdb := getDB(r)
	// Check is user registered before or not

	creds := mdb.GetCredsByEmail(r.Context(), email)

	if len(creds.Email) != 0 {
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

	createResponse(role, w, r.Context(), mdb)
}

func createResponse(creds db.UserCredentials, w http.ResponseWriter, ctx context.Context, mdb *db.DBHandle) {
	token, err := generateToken(creds.Email, creds.Role, ExpiresHours)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	encoder := json.NewEncoder(w)
	encoder.Encode(AuthResponse{
		Token: token,
	})
}

package routes

import (
	"encoding/json"
	"net/http"
	"time"

	"enstrurent.com/server/db"
	"go.mongodb.org/mongo-driver/bson"
)

func getClientInfo(w http.ResponseWriter, r *http.Request) {
	email, err := validateClient(r)

	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	mdb := getDB(r)
	clientInfo := mdb.GetClientByEmail(r.Context(), email)

	json.NewEncoder(w).Encode(clientInfo)
}

func updateClientInfo(w http.ResponseWriter, r *http.Request) {
	email, err := validateClient(r)

	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	var clientUpdated db.Client
	json.NewDecoder(r.Body).Decode(&clientUpdated)

	if email != clientUpdated.Email {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	clientUpdated.UpdatedAt = time.Now()

	mdb := getDB(r)
	result := mdb.MongoDB().Collection(db.ClientCollection).
		FindOneAndReplace(r.Context(), bson.M{"email": email}, clientUpdated)
	if result.Err() != nil {
		http.Error(w, result.Err().Error(), http.StatusBadRequest)
		return
	}
	json.NewEncoder(w).Encode(clientUpdated)
}

func validateClient(r *http.Request) (string, error) {
	role := r.Context().Value(UserRoleContext)
	if role != ClientRole {
		return "", ErrUnauthorized
	}
	return r.Context().Value(UserEmailContext).(string), nil
}

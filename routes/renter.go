package routes

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"enstrurent.com/server/db"
	"go.mongodb.org/mongo-driver/bson"
)

func getRenterInfo(w http.ResponseWriter, r *http.Request) {
	role := r.Context().Value(UserRoleContext)
	if role != RenterRole {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	email := r.Context().Value(UserEmailContext).(string)

	mdb := getDB(r)
	renterInfo := mdb.GetRenterByEmail(r.Context(), email)

	json.NewEncoder(w).Encode(renterInfo)
}

func updateRenterInfo(w http.ResponseWriter, r *http.Request) {
	role := r.Context().Value(UserRoleContext)
	if role != RenterRole {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	email := r.Context().Value(UserEmailContext).(string)

	var renterInfo db.Renter
	json.NewDecoder(r.Body).Decode(&renterInfo)

	if renterInfo.Email != email {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	renterInfo.UpdatedAt = time.Now()

	result := getDB(r).RenterCollection().
		FindOneAndReplace(r.Context(), bson.M{"_id": renterInfo.ID}, renterInfo)

	if result.Err() != nil {
		fmt.Println(result.Err().Error())
		http.Error(w, "something went wrong", http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(renterInfo)
}

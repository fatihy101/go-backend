package routes

import (
	"encoding/json"
	"net/http"
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

}
func updateStoreHeader(w http.ResponseWriter, r *http.Request) {

}

func updateStorePicture(w http.ResponseWriter, r *http.Request) {

}

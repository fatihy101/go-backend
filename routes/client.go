package routes

import (
	"encoding/json"
	"net/http"
)

func getClientInfo(w http.ResponseWriter, r *http.Request) {
	role := r.Context().Value(UserRoleContext)
	if role != ClientRole {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	email := r.Context().Value(UserEmailContext).(string)

	mdb := getDB(r)
	clientInfo := mdb.GetClientByEmail(r.Context(), email)

	json.NewEncoder(w).Encode(clientInfo)
}

func updateClientInfo(w http.ResponseWriter, r *http.Request) {

}

func updateClientPicture(w http.ResponseWriter, r *http.Request) {

}

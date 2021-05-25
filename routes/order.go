package routes

import (
	"net/http"

	"enstrurent.com/server/db"
)

const (
	RentOrder     string = "rental"
	PurchaseOrder string = "purchase"
)

func getOrdersByEmail(w http.ResponseWriter, r *http.Request) {

}

func createOrder(w http.ResponseWriter, r *http.Request) {
	email, err := validateClient(r)

	if err != nil {
		http.Error(w, ErrUnauthorized.Error(), http.StatusUnauthorized)
		return
	}
	orderType := r.Header.Get("order_type")
	mdb := getDB(r)
	clientID := mdb.GetClientByEmail(r.Context(), email).ID.Hex()

	if orderType == RentOrder {
		var order db.RentOrder
		order.InitializeOrder(clientID, mdb, w, r)
	} else if orderType == PurchaseOrder {
		var order db.PurchaseOrder
		order.InitializeOrder(clientID, mdb, w, r)
	} else {
		http.Error(w, "order type is empty", http.StatusBadRequest)
		return
	}
}

func updateOrder(w http.ResponseWriter, r *http.Request) {

}

func cancelOrder(w http.ResponseWriter, r *http.Request) {
	// TODO Soft delete
}

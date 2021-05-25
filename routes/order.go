package routes

import (
	"encoding/json"
	"net/http"

	"enstrurent.com/server/db"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	RentOrder     string = "rental"
	PurchaseOrder string = "purchase"
)

func getOrdersByEmail(w http.ResponseWriter, r *http.Request) {
	email, err := validateClient(r)

	if err != nil {
		http.Error(w, ErrUnauthorized.Error(), http.StatusUnauthorized)
		return
	}

	mdb := getDB(r)
	client := mdb.GetClientByEmail(r.Context(), email)

	cursor, err := mdb.OrdersCollection().Find(r.Context(), bson.M{"client_id": client.ID.Hex()})

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	var results []map[string]interface{}
	err = cursor.All(r.Context(), &results)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	json.NewEncoder(w).Encode(results)
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
		OrdersCommon(&order, clientID, mdb, w, r)

	} else if orderType == PurchaseOrder {
		var order db.PurchaseOrder
		OrdersCommon(&order, clientID, mdb, w, r)

	} else {
		http.Error(w, "order type is empty", http.StatusBadRequest)
		return
	}
}
func OrdersCommon(order db.IOrder, clientID string, mdb *db.DBHandle, w http.ResponseWriter, r *http.Request) {
	order.InitializeOrder(clientID, mdb)
	json.NewDecoder(r.Body).Decode(&order)
	order.InitializeOrder(clientID, mdb)
	result, err := mdb.OrdersCollection().InsertOne(r.Context(), order)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	order.SetID(result.InsertedID.(primitive.ObjectID))

	json.NewEncoder(w).Encode(order)
}

func updateOrder(w http.ResponseWriter, r *http.Request) {

}

func cancelOrder(w http.ResponseWriter, r *http.Request) {
	// TODO Soft delete
}

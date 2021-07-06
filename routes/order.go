package routes

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"enstrurent.com/server/db"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	RentOrder     string = "rental"
	PurchaseOrder string = "purchase"
)

func getOrdersByEmail(w http.ResponseWriter, r *http.Request) {
	mdb := getDB(r)

	role := r.Context().Value(UserRoleContext)
	email := r.Context().Value(UserEmailContext).(string)
	var cursor *mongo.Cursor
	var err error
	if role == ClientRole {
		client := mdb.GetClientByEmail(r.Context(), email)
		cursor, err = mdb.OrdersCollection().Find(r.Context(), bson.M{"client_id": client.ID.Hex()})
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

	} else if role == RenterRole {
		renter := mdb.GetRenterByEmail(r.Context(), email)
		cursor, err = mdb.OrdersCollection().Find(r.Context(), bson.M{"renter_id": renter.ID})
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
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

func updateOrderStatus(w http.ResponseWriter, r *http.Request) { //FIXME
	_, err := validateUserForOrder(r)
	if err != nil {
		if errors.Is(err, ErrUnauthorized) {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		} else {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
	}
	var updatedOrder map[string]interface{}
	json.NewDecoder(r.Body).Decode(&updatedOrder)
	// collection := getDB(r).OrdersCollection()
	// collection.ReplaceOne()
}

func validateUserForOrder(r *http.Request) (interface{}, error) {
	orderID := r.Header.Get("order_id")
	id, err := primitive.ObjectIDFromHex(orderID)
	if err != nil {
		return nil, ErrWentWrong
	}
	userRole := r.Context().Value(UserRoleContext)
	userEmail := r.Context().Value(UserEmailContext).(string)
	mdb := getDB(r)
	result := mdb.OrdersCollection().FindOne(r.Context(), bson.M{"_id": id})

	if result.Err() != nil {
		fmt.Println(result.Err())
		return nil, ErrWentWrong
	}
	var order map[string]interface{}
	result.Decode(&order)

	if userRole == ClientRole {
		client := mdb.GetClientByEmail(r.Context(), userEmail)
		if client.ID.Hex() != order["client_id"] {
			return nil, ErrUnauthorized
		}
	} else if userRole == RenterRole {
		renter := mdb.GetRenterByEmail(r.Context(), userEmail)
		// We're getting the product of the order for validating renter's id.
		productID, _ := primitive.ObjectIDFromHex(fmt.Sprint(order["product_id"]))

		if product, err := mdb.GetProductByID(r.Context(), productID); err != nil {
			return nil, err
		} else if product.RenterID != renter.ID.Hex() {
			return nil, ErrUnauthorized
		}
	} else {
		return nil, ErrUnauthorized
	}
	return order, nil
}

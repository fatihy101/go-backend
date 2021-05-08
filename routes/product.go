package routes

import (
	"encoding/json"
	"net/http"
	"time"

	"enstrurent.com/server/db"
	"github.com/go-chi/chi/v5"
	"go.mongodb.org/mongo-driver/bson"
)

func getOneProduct(w http.ResponseWriter, r *http.Request) { // Guest OP
	collection := getCollection(r, db.ProductCollection)
	singleResult := collection.FindOne(r.Context(), bson.M{"_id": chi.URLParam(r, "id")})
	var product db.Product
	singleResult.Decode(&product)

	json.NewEncoder(w).Encode(product)
}

func getAllProducts(w http.ResponseWriter, r *http.Request) { // Guest OP
	// TODO do by location
	collection := getCollection(r, db.ProductCollection)
	mCursor, err := collection.Find(r.Context(), bson.M{"deleted_at": time.Time{}})

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var products []db.Product
	err = mCursor.Decode(&products)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	json.NewEncoder(w).Encode(products)
}

func getRenterProducts(w http.ResponseWriter, r *http.Request) { // Renter OP
	role := r.Context().Value(UserRoleContext)
	if role != RenterRole {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	email := r.Context().Value(UserEmailContext).(string)

	mdb := getDB(r)
	renter := mdb.GetRenterByEmail(r.Context(), email)
	filter := bson.M{"renterID": renter.ID.String()}
	mCursor, err := mdb.MongoDB().Collection(db.ProductCollection).Find(r.Context(), filter)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var products []db.Product
	mCursor.Decode(&products)
	json.NewEncoder(w).Encode(products)
}

func deleteProduct(w http.ResponseWriter, r *http.Request) { // Renter OP
	role := r.Context().Value(UserRoleContext)
	if role != RenterRole {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	email := r.Context().Value(UserEmailContext).(string)

	mdb := getDB(r)
	renter := mdb.GetRenterByEmail(r.Context(), email)
	id := chi.URLParam(r, "id")
	filter := bson.M{"_id": id}
	singleResult := mdb.MongoDB().Collection(db.ProductCollection).FindOne(r.Context(), filter) // FIXME Test the filter

	var product db.Product
	singleResult.Decode(&product)
	if renter.ID.String() != product.RenterID {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	_, err := mdb.MongoDB().Collection(db.ProductCollection).DeleteOne(r.Context(), filter)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}

func addProduct(w http.ResponseWriter, r *http.Request) { // Renter OP
	role := r.Context().Value(UserRoleContext)
	if role != RenterRole {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	email := r.Context().Value(UserEmailContext).(string)

	mdb := getDB(r)
	renter := mdb.GetRenterByEmail(r.Context(), email)

	var newProduct db.Product
	json.NewDecoder(r.Body).Decode(&newProduct)
	newProduct.RenterID = renter.ID.String()
	newProduct.CreatedAt = time.Now()
	newProduct.UpdatedAt = time.Now()
	_, err := mdb.SaveOne(db.ProductCollection, r.Context(), newProduct)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}

func updateProduct(w http.ResponseWriter, r *http.Request) { // Renter OP
	role := r.Context().Value(UserRoleContext)
	if role != RenterRole {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	email := r.Context().Value(UserEmailContext).(string)

	mdb := getDB(r)
	renter := mdb.GetRenterByEmail(r.Context(), email)

	var product db.Product

	json.NewDecoder(r.Body).Decode(&product)
	// Check the product's renterID and token renter's id.
	if product.RenterID != renter.ID.String() {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	_, err := mdb.MongoDB().Collection(db.ProductCollection).UpdateByID(r.Context(), product.ID, product)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}

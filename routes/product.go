package routes

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"enstrurent.com/server/db"
	"github.com/go-chi/chi/v5"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func getOneProduct(w http.ResponseWriter, r *http.Request) { // Guest OP
	collection := getDB(r).ProductCollection()
	id, _ := primitive.ObjectIDFromHex(chi.URLParam(r, "id"))
	singleResult := collection.FindOne(r.Context(), bson.M{"_id": id})
	if singleResult.Err() != nil {
		fmt.Println("Error on get product: " + singleResult.Err().Error())
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("there is no product with specified id"))
		return
	}
	var product db.Product
	err := singleResult.Decode(&product)

	if err != nil {
		fmt.Println("Error on get product: " + err.Error())
		return
	}
	json.NewEncoder(w).Encode(product)
}

func getAllProducts(w http.ResponseWriter, r *http.Request) { // Guest OP
	// TODO do by location and paginate
	collection := getDB(r).ProductCollection()
	mCursor, err := collection.Find(r.Context(), bson.D{})
	if mCursor.RemainingBatchLength() == 0 {
		http.Error(w, "No item in DB", http.StatusBadRequest)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var products []db.Product
	err = mCursor.All(r.Context(), &products)

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
	collection := mdb.ProductCollection()

	renter := mdb.GetRenterByEmail(r.Context(), email)
	filter := bson.M{"renter_id": renter.ID.Hex()}
	mCursor, err := collection.Find(r.Context(), filter)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	var products []db.Product
	err = mCursor.All(r.Context(), &products)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
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
	collection := mdb.ProductCollection()
	renter := mdb.GetRenterByEmail(r.Context(), email)
	id, _ := primitive.ObjectIDFromHex(chi.URLParam(r, "id"))
	filter := bson.M{"_id": id}
	singleResult := collection.FindOne(r.Context(), filter)

	var product db.Product
	singleResult.Decode(&product)
	if renter.ID.Hex() != product.RenterID {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	_, err := collection.DeleteOne(r.Context(), filter)

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
	newProduct.RenterID = renter.ID.Hex()
	newProduct.CreatedAt = time.Now()
	newProduct.UpdatedAt = time.Now()
	newProduct.City = renter.RenterAddress.City
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
	collection := mdb.ProductCollection()
	renter := mdb.GetRenterByEmail(r.Context(), email)

	var product db.Product

	json.NewDecoder(r.Body).Decode(&product)
	// Check the product's renterID and token renter's id.
	if product.RenterID != renter.ID.Hex() {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	product.UpdatedAt = time.Now()
	res := collection.FindOneAndReplace(r.Context(), bson.M{"_id": product.ID}, product)

	if res.Err() != nil {
		http.Error(w, res.Err().Error(), http.StatusBadRequest)
		return
	}
}

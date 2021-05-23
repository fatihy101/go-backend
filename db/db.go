package db

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	UserCredsCollection = "credentials"
	ClientCollection    = "clients"
	RenterCollection    = "renters"
	AddressCollection   = "addresses"
	OrderCollection     = "orders"
	PhotoCollection     = "photos"
	ProductCollection   = "products"
	CitiesCollection    = "cities"
)

func OpenConnection(conString string, dbName string) *DBHandle {
	ctx := context.TODO()
	clientOptions := options.Client().ApplyURI(conString)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	if err = client.Ping(ctx, nil); err != nil {
		log.Fatal(err)
	}

	return &DBHandle{mdb: client.Database(dbName)}
}

func (d *DBHandle) MongoDB() *mongo.Database {
	return d.mdb
}

func (d *DBHandle) SaveOne(collection string, ctx context.Context, data interface{}) (interface{}, error) {
	id, err := d.mdb.Collection(collection).InsertOne(ctx, data)
	if err != nil {
		return nil, err
	}
	return id.InsertedID, nil
}

func (d *DBHandle) GetCredsByEmail(ctx context.Context, email string) UserCredentials {
	var res UserCredentials
	d.mdb.Collection(UserCredsCollection).FindOne(ctx, bson.M{"email": email}).Decode(&res)
	return res
}

func (d *DBHandle) GetRenterByEmail(ctx context.Context, email string) Renter {
	var res Renter
	ErrorCheck(d.mdb.Collection(RenterCollection).FindOne(ctx, bson.M{"email": email}).Decode(&res))
	return res
}

func (d *DBHandle) GetClientByEmail(ctx context.Context, email string) Client {
	var res Client
	ErrorCheck(d.mdb.Collection(ClientCollection).FindOne(ctx, bson.M{"email": email}).Decode(&res))
	return res
}

func (d *DBHandle) ProductCollection() *mongo.Collection {
	return d.mdb.Collection(ProductCollection)
}

func (d *DBHandle) CitiesCollection() *mongo.Collection {
	return d.mdb.Collection(CitiesCollection)
}

func (d *DBHandle) AddressCollection() *mongo.Collection {
	return d.mdb.Collection(AddressCollection)
}

func (d *DBHandle) RenterCollection() *mongo.Collection {
	return d.mdb.Collection(RenterCollection)
}

func ErrorCheck(err error) {
	if err != nil {
		log.Fatal(fmt.Sprintf("Error on db get: %v", err))
	}
}

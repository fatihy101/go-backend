package db

import (
	"context"
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

func (d *DBHandle) SaveOne(collection string, ctx context.Context, data interface{}) interface{} {
	id, err := d.mdb.Collection(collection).InsertOne(ctx, data)
	if err != nil {
		log.Fatal(err)
	}
	return id.InsertedID
}

func (d *DBHandle) GetCredsByEmail(ctx context.Context, email string, qResult *bson.M) {
	d.mdb.Collection(UserCredsCollection).FindOne(ctx, bson.M{"email": email}).Decode(qResult)
}

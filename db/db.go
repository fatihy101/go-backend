package db

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (d *DBHandle) MongoDB() *mongo.Database {
	return d.db
}

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

	return &DBHandle{db: client.Database(dbName)}
}

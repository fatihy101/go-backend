package db

import (
	"log"

	"github.com/go-bongo/bongo"
)

func OpenConnection(conString string, dbName string) *DBHandle {
	config := &bongo.Config{
		ConnectionString: conString,
		Database:         dbName,
	}

	connection, err := bongo.Connect(config)

	if err != nil {
		log.Fatal(err)
	}

	return &DBHandle{db: connection}
}

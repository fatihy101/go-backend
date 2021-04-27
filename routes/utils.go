package routes

import (
	"net/http"

	"enstrurent.com/server/db"
	"go.mongodb.org/mongo-driver/mongo"
)

func getDB(r *http.Request) *mongo.Database {
	return r.Context().Value(DBContext).(*db.DBHandle).MongoDB()
}

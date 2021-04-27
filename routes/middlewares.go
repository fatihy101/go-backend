package routes

import (
	"context"
	"net/http"

	"enstrurent.com/server/db"
)

type key int // To supress a warning
const (
	DBContext key = iota
	UserEmailContext
	UserRoleContext
)

func DBMiddleware(db *db.DBHandle) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), DBContext, db)))
		})
	}
}

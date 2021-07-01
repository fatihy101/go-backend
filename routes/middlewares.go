package routes

import (
	"context"
	"fmt"
	"net/http"

	"enstrurent.com/server/db"
	"github.com/dgrijalva/jwt-go"
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

func JSONResponseMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		next.ServeHTTP(w, r)
	})
}

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get cookie from request.
		token_str := r.Header.Get("token")
		// Check token string is empty.
		if token_str == "" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		// Check token.
		token, err := checkToken(token_str)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		// Get email from token.
		claims := token.Claims.(jwt.MapClaims)
		email := claims["email"]
		if email == nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		mdb := getDB(r)
		creds := mdb.GetCredsByEmail(r.Context(), fmt.Sprint(email))
		role := claims["role"].(string)
		if creds.Role != role {
			http.Error(w, "role and email does not match", http.StatusUnauthorized)
			return
		}
		// Pass the email to context
		ctx := r.Context()
		ctx = context.WithValue(ctx, UserEmailContext, email)
		ctx = context.WithValue(ctx, UserRoleContext, role)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func AllowOrigin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET,PUT,POST,DELETE,OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Access-Control-Allow-Headers, Authorization, X-Requested-With, token")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		if r.Method == "OPTIONS" {
			return
		}
		next.ServeHTTP(w, r)
	})
}

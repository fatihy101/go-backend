package routes

import (
	"fmt"
	"net/http"
	"time"

	"enstrurent.com/server/db"
	"enstrurent.com/server/flags"
	"github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) string {
	val, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		return ""
	}
	return string(val)
}

func CompareHashAndPassword(hash string, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func getDB(r *http.Request) *db.DBHandle {
	return r.Context().Value(DBContext).(*db.DBHandle)
}

func getCollection(r *http.Request, collectionName string) *mongo.Collection {
	return r.Context().Value(DBContext).(*db.DBHandle).MongoDB().Collection(collectionName)
}

func generateToken(email string, role string, expires time.Duration) (token string, err error) {
	claims := jwt.MapClaims{}
	claims["email"] = email
	claims["role"] = role
	claims["exp"] = time.Now().Add(expires).Unix()
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(flags.GetConfig().JWT_KEY))
}

func checkToken(token string) (*jwt.Token, error) {
	return jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {

		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(flags.GetConfig().JWT_KEY), nil
	})
}

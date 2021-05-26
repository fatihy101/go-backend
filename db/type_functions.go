package db

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (order *OrderBase) InitializeOrder(clientID string, mdb *DBHandle) {
	order.ClientID = clientID
	order.OrderStatus = "Sipariş Alındı"
	order.CreatedAt = time.Now()
}

func (order *OrderBase) SetID(clientID primitive.ObjectID) {
	order.ID = clientID
}

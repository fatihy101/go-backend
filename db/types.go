package db

import (
	"time"

	"github.com/go-bongo/bongo"
)

type DBHandle struct {
	db *bongo.Connection
}

type Address struct {
	bongo.DocumentBase `bson:",inline"`
	Title,
	City,
	State,
	Zipcode,
	Street,
	BuildingNo,
	Description string
}

type Client struct {
	bongo.DocumentBase `bson:",inline"`
	BirthdayDate       time.Time
	Email,
	Name,
	Surname,
	PhoneNumber string
	ImagePath     Image
	ClientAddress Address
}

type Renter struct { // TODO: FUTURE Add geolocation on search from map.
	bongo.DocumentBase `bson:",inline"`
	Name,
	Surname,
	Email,
	StoreInfo,
	StoreName,
	PhoneNumber string
	Rating         float32
	ProfilePicPath Image
	HeaderPicPath  Image
	RenterAdress   Address
}

type Order struct {
	bongo.DocumentBase `bson:",inline"`
	Product            Product
	Client             Client
	Renter             Renter
	DeliveryType       string
	Address            Address
	IsRental           bool
	InitalImages       []Image
}

type Image struct {
	bongo.DocumentBase `bson:",inline"`
	ImageName          string
}

type Product struct {
	bongo.DocumentBase `bson:",inline"`
}

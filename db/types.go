package db

import (
	"time"

	"github.com/go-bongo/bongo"
)

type Address struct {
	bongo.DocumentBase `bson:",inline"`
	Title              string
	City               string
	State              string
	Zipcode            string
	Street             string
	BuildingNo         string
	Description        string
}

type Client struct {
	bongo.DocumentBase `bson:",inline"`
	BirthdayDate       time.Time
	Email              string
	Name               string
	Surname            string
	PhoneNumber        string
	ImagePath          Image
	ClientAddress      Address
}

type Renter struct { // TODO: FUTURE Add geolocation on search from map.
	bongo.DocumentBase `bson:",inline"`
	Name               string
	Surname            string
	Email              string
	StoreInfo          string
	StoreName          string
	PhoneNumber        string
	Rating             float32
	ProfilePicPath     Image
	HeaderPicPath      Image
	RenterAdress       Address
}

type User struct { // Renter + Client
	bongo.DocumentBase `bson:",inline"`
	Type               string
	UID                string
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

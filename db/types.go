package db

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type DBHandle struct {
	db *mongo.Database
}

type Base struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	CreatedAt time.Time          `bson:"created_at,omitempty"`
	UpdatedAt time.Time          `bson:"updated_at,omitempty"`
}

type UserRole struct {
	Base  `bson:",inline"`
	Email string `json:"email"`
	Role  string `json:"role"`
}

type Address struct {
	Base        `bson:",inline"`
	Title       string `json:"title"`
	City        string `json:"city"`
	State       string `json:"state"`
	Zipcode     string `json:"zipcode"`
	Street      string `json:"street"`
	BuildingNo  string `json:"building_no"`
	Description string `json:"description"`
}

type Client struct {
	Base          `bson:",inline"`
	BirthdayDate  time.Time
	Email         string  `json:"email"`
	Name          string  `json:"name"`
	Surname       string  `json:"surname"`
	PhoneNumber   string  `json:"phone_number"`
	ImagePath     Image   `json:"image_path"`
	ClientAddress Address `json:"client_address"`
}

type Renter struct { // TODO: FUTURE Add geolocation on search from map.
	Base           `bson:",inline"`
	Name           string  `json:"name"`
	Surname        string  `json:"surname"`
	Email          string  `json:"email"`
	StoreInfo      string  `json:"store_info"`
	StoreName      string  `json:"store_name"`
	PhoneNumber    string  `json:"phone_number"`
	Rating         float32 `json:"rating"`
	ProfilePicPath Image   `json:"profile_pic_path"`
	HeaderPicPath  Image   `json:"header_pic_path"`
	RenterAddress  Address `json:"renter_address"`
}

type Order struct {
	Base         `bson:",inline"`
	Product      Product
	Client       Client
	Renter       Renter
	DeliveryType string
	Address      Address
	IsRental     bool
	InitalImages []Image
}

type Image struct {
	Base      `bson:",inline"`
	ImageName string
}

type Product struct {
	Base `bson:",inline"`
}

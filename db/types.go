package db

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type DBHandle struct {
	mdb *mongo.Database
}

type Base struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	CreatedAt time.Time          `bson:"created_at,omitempty"`
	UpdatedAt time.Time          `bson:"updated_at,omitempty"`
}

type UserCredentials struct {
	Base     `bson:",inline"`
	Email    string `json:"email"`
	Role     string `json:"role"`
	Password string `json:"password"`
}

type Address struct {
	Base        `bson:",inline"`
	Title       string `json:"title" bson:"title"`
	City        string `json:"city" bson:"city"`
	State       string `json:"state" bson:"state"`
	Zipcode     string `json:"zipcode" bson:"zipcode"`
	Street      string `json:"street" bson:"street"`
	BuildingNo  string `json:"building_no" bson:"building_no"`
	Description string `json:"description" bson:"description"`
}

type Client struct {
	Base          `bson:",inline"`
	BirthdayDate  time.Time
	Email         string  `json:"email" bson:"email"`
	Name          string  `json:"name" bson:"name"`
	Surname       string  `json:"surname" bson:"surname"`
	PhoneNumber   string  `json:"phone_number" bson:"phone_number"`
	ImagePath     Image   `json:"image_path" bson:"image_path"`
	ClientAddress Address `json:"client_address" bson:"client_address"`
}

type Renter struct { // TODO: FUTURE Add geolocation on search from map.
	Base           `bson:",inline"`
	Name           string  `json:"name" bson:"name"`
	Surname        string  `json:"surname" bson:"surname"`
	Email          string  `json:"email" bson:"email"`
	StoreInfo      string  `json:"store_info" bson:"store_info"`
	StoreName      string  `json:"store_name" bson:"store_name"`
	PhoneNumber    string  `json:"phone_number" bson:"phone_number"`
	Rating         float32 `json:"rating" bson:"rating"`
	ProfilePicPath Image   `json:"profile_pic_path" bson:"profile_pic_path,omitempty"`
	HeaderPicPath  Image   `json:"header_pic_path" bson:"header_pic_path,omitempty"`
	RenterAddress  Address `json:"renter_address" bson:"renter_address,omitempty"`
}

type Order struct {
	Base         `bson:",inline"`
	Product      Product
	Client       Client
	Renter       Renter
	DeliveryType string `json:"delivery_type" bson:"delivery_type"`
	Address      Address
	IsRental     bool `json:"is_rental" bson:"is_rental"`
	InitalImages []Image
}

type Image struct {
	Base      `bson:",inline"`
	ImageName string
}

type Product struct {
	Base `bson:",inline"`
}

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
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	CreatedAt time.Time          `bson:"created_at,omitempty"`
	UpdatedAt time.Time          `bson:"updated_at,omitempty"`
	DeletedAt time.Time          `bson:"deleted_at,omitempty"`
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
	AddressLine string `json:"address_line" bson:"address_line"`
	Description string `json:"description" bson:"description"`
}

type Client struct {
	Base         `bson:",inline"`
	BirthdayDate time.Time
	Email        string `json:"email" bson:"email"`
	Name         string `json:"name" bson:"name"`
	Surname      string `json:"surname" bson:"surname"`
	PhoneNumber  string `json:"phone_number" bson:"phone_number"`
	ImageName    string `json:"image_name" bson:"image_name"`
	AddressID    string `json:"address_id" bson:"address_id"`
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
	ProfilePicName string  `json:"profile_pic_name" bson:"profile_pic_name,omitempty"`
	HeaderPicName  string  `json:"header_pic_name" bson:"header_pic_name,omitempty"`
	AddressID      string  `json:"address_id" bson:"address_id"`
}

type IOrder interface {
	InitializeOrder(clientID string, mdb *DBHandle)
	SetID(ClientID primitive.ObjectID)
}

type OrderBase struct {
	Base          `bson:",inline"`
	ProductID     string `json:"product_id" bson:"product_id"`
	ClientID      string `json:"client_id" bson:"client_id"`
	AddressID     string `json:"address_id" bson:"address_id"`
	RenterID      string `json:"renter_id" bson:"renter_id"`
	DeliveryType  string `json:"delivery_type" bson:"delivery_type"`
	PaymentMethod string `json:"payment_method" bson:"payment_method"`
	OrderStatus   string `json:"order_status" bson:"order_status"`
}

type RentOrder struct {
	OrderBase        `bson:",inline"`
	InitalImageNames []string  `json:"initial_image_names" bson:"initial_image_names"`
	FinalImageNames  []string  `json:"final_image_names" bson:"final_image_names"`
	DepositPrice     float32   `json:"deposit_price" bson:"deposit_price"`
	RentingPrice     float32   `json:"renting_price" bson:"renting_price"`
	RentedDateRange  DateRange `json:"rented_date_range" bson:"rented_date_range"`
}
type DateRange struct {
	Start string `json:"start" bson:"start"`
	End   string `json:"end" bson:"end"`
}

type PurchaseOrder struct {
	OrderBase `bson:",inline"`
	Price     float32 `json:"price" bson:"price"`
}

// TODO Split the data as "rentaL_fields" and "selling_fields" as two.
type Product struct {
	Base              `bson:",inline"`
	RenterID          string        `json:"renter_id" bson:"renter_id"`
	City              string        `json:"city" bson:"city"`
	Category          string        `json:"category" bson:"category"`
	Brand             string        `json:"brand" bson:"brand"`
	Model             string        `json:"model" bson:"model"`
	Info              string        `json:"info" bson:"info"`
	IsRental          bool          `json:"is_rental" bson:"is_rental"`
	IsDepositRequired bool          `json:"is_deposit_required" bson:"is_deposit_required"`
	IsOpenToSell      bool          `json:"is_open_to_sell" bson:"is_open_to_sell"`
	IsUsed            bool          `json:"is_used" bson:"is_used"`
	MaxRentalDays     int           `json:"max_rental_days" bson:"max_rental_days"`
	DailyPrice        float32       `json:"daily_price" bson:"daily_price"`
	FullPrice         float32       `json:"full_price" bson:"full_price"`
	DepositPrice      float32       `json:"deposit_price" bson:"deposit_price"`
	StockQuantity     int           `json:"stock_quantity" bson:"stock_quantity"`
	DeliveryTypes     []string      `json:"delivery_types" bson:"delivery_types"`
	ImageNames        []string      `json:"image_names" bson:"image_names"`
	ThumbnailNames    []string      `json:"thumbnail_names" bson:"thumbnail_names"`
	Tags              []string      `json:"tags" bson:"tags"`
	PaymentMethods    []string      `json:"payment_methods" bson:"payment_methods"`
	RentedDaysRanges  [][]time.Time `json:"rented_days_ranges" bson:"rented_days_ranges"`
	IsPublished       bool          `json:"is_published" bson:"is_published"` // products can be saved as draft
}

type City struct {
	ID        int    `json:"id" bson:"_id"`
	Name      string `json:"name" bson:"name"`
	Latitude  string `json:"latitude" bson:"latitude"`
	Longitude string `json:"longtitude" bson:"longtitude"`
}

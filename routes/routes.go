package routes

import (
	"encoding/json"
	"errors"
	"net/http"

	"enstrurent.com/server/db"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
)

var (
	ErrWentWrong    error = errors.New("something went wrong")
	ErrUnauthorized error = errors.New("unauthorized act")
)

func Routes(db *db.DBHandle) *chi.Mux {
	router := chi.NewRouter()
	router.Use(middleware.Logger)
	router.Use(AllowOrigin)
	router.Use(DBMiddleware(db))
	router.Use(JSONResponseMiddleware)
	router.Route("/", allRoutes)

	return router
}

func allRoutes(r chi.Router) {
	r.Get("/", func(rw http.ResponseWriter, r *http.Request) {
		json.NewEncoder(rw).Encode(map[string]string{"server": "active"})
	})
	r.Route("/auth", authRoutes)
	r.Route("/products", productRoutes)
	r.With(AuthMiddleware).Route("/clients", clientRoutes)
	r.With(AuthMiddleware).Route("/renters", renterRoutes)
	r.With(AuthMiddleware).Route("/orders", orderRoutes)
	r.Route("/addresses", addressRoutes)
	r.Route("/images", ImageRoutes)
}

func authRoutes(r chi.Router) {
	r.Post("/login", login)
	r.Post("/sign_up", signUp)
	// Token validation occurs in middleware. Function is placeholder.
	r.With(AuthMiddleware).Post("/validate", func(rw http.ResponseWriter, r *http.Request) {})
}

func productRoutes(r chi.Router) {
	r.With(AuthMiddleware).Post("/", addProduct)
	r.Get("/", getAllProducts) // TODO paginate
	r.Get("/{id}", getOneProduct)
	r.With(AuthMiddleware).Delete("/{id}", deleteProduct)
	r.With(AuthMiddleware).Put("/", updateProduct)
	r.With(AuthMiddleware).Get("/by_renter", getRenterProducts)
}

func clientRoutes(r chi.Router) {
	r.Get("/", getClientInfo)
	r.Put("/", updateClientInfo)
}

func renterRoutes(r chi.Router) {
	r.Get("/", getRenterInfo)
	r.Put("/", updateRenterInfo)
}

func orderRoutes(r chi.Router) {
	r.Get("/", getOrdersByEmail) // TODO paginate
	r.Post("/", createOrder)
	r.Put("/", updateOrderStatus)
	// r.Delete("/{order_id}", cancelOrder)
}

func addressRoutes(r chi.Router) {
	r.Get("/{address_id}", getAddress)
	r.With(AuthMiddleware).Delete("/{address_id}", deleteAddress)
	r.Get("/cities", getCities)
	r.Post("/", addAddress)
	r.With(AuthMiddleware).Put("/", updateAddress)

}

func ImageRoutes(r chi.Router) {
	r.With(AuthMiddleware).Post("/", addPhoto)
	r.With(AuthMiddleware).Delete("/{id}&{thumbnail}", deletePhoto)
	r.Get("/{photoName}", getPhoto)
}

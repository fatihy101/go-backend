package routes

import (
	"enstrurent.com/server/db"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
)

const (
	UserRoleCollection = "user_types"
	ClientCollection   = "clients"
	RenterCollection   = "renters"
	AddressCollection  = "addresses"
	OrderCollection    = "orders"
	PhotoCollection    = "photos"
	ProductCollection  = "products"
)

func Routes(db *db.DBHandle) *chi.Mux {
	router := chi.NewRouter()
	router.Use(middleware.Logger)
	router.Use(DBMiddleware(db))
	router.Route("/", allRoutes)

	return router
}

func allRoutes(r chi.Router) {
	r.Route("/auth", authRoutes)
	r.Route("/products", productRoutes)
	r.Route("/clients", clientRoutes)
	r.Route("/renters", renterRoutes)
	r.Route("/orders", orderRoutes)
	r.Route("/addresses", addressRoutes)
}

func authRoutes(r chi.Router) {
	r.Post("/login", login)
	r.Post("/sign_up_client", signUpClient)
}

func productRoutes(r chi.Router) {
}

func clientRoutes(r chi.Router) {
	r.Get("/", getOneClient)
}

func renterRoutes(r chi.Router) {
}

func orderRoutes(r chi.Router) {

}

func addressRoutes(r chi.Router) {

}

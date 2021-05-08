package routes

import (
	"encoding/json"
	"net/http"

	"enstrurent.com/server/db"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
)

func Routes(db *db.DBHandle) *chi.Mux {
	router := chi.NewRouter()
	router.Use(middleware.Logger)
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
}

func authRoutes(r chi.Router) {
	r.Post("/login", login)
	r.Post("/sign_up", signUp)
}

func productRoutes(r chi.Router) {
	r.With(AuthMiddleware).Post("/", addProduct)
	r.Get("/", getAllProducts) // TODO paginate
	r.Get("/{id}", getOneProduct)
	r.With(AuthMiddleware).Delete("/{id}", deleteProduct)
	r.With(AuthMiddleware).Put("/product", updateProduct)
	r.With(AuthMiddleware).Get("/by_renter", getRenterProducts)
}

func clientRoutes(r chi.Router) {
	r.Get("/", getClientInfo)
	r.Put("/", updateClientInfo)
	r.Put("/profile_picture", updateClientPicture)
}

func renterRoutes(r chi.Router) { // TODO Write a middleware for checking the is renters product belongings
	r.Get("/", getRenterInfo)
	r.Put("/header", updateStoreHeader)
	r.Put("/info", updateRenterInfo)
	r.Put("/profile_picture", updateStorePicture)
}

func orderRoutes(r chi.Router) {
	r.Get("/", getOrdersByEmail) // TODO paginate
	r.Post("/", createOrder)
	r.Put("/", updateOrder)
	r.Delete("/{order_id}", cancelOrder)
}

func addressRoutes(r chi.Router) {

}

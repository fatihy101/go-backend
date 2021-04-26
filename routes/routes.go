package routes

import (
	"enstrurent.com/server/db"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
)

func Routes(db *db.DBHandle) *chi.Mux {
	router := chi.NewRouter()
	router.Use(middleware.Logger)
	router.Route("/", allRoutes)

	return router
}

func allRoutes(r chi.Router) {
	// r.Route("/auth")
	// r.Route("/products")
	// r.Route("all_user_types")
	// r.Route("/clients")
	// r.Route("/renters")
	// r.Route("/orders")
	// r.Route("/addresses")
}

package main

import (
	"fmt"
	"log"
	"net/http"

	"enstrurent.com/server/db"
	"enstrurent.com/server/routes"
)

func main() {
	// TODO get configuration from config.json file
	var port = ":4002"
	dbCon := db.OpenConnection("localhost", "enstrurent")

	server := &http.Server{
		Addr:    port,
		Handler: routes.Routes(dbCon),
	}

	if err := server.ListenAndServe(); err != nil {
		log.Fatal(fmt.Sprintf("Error on listen and serve: %v", err))
	}
	fmt.Printf("Serving on %v\n", port)
}

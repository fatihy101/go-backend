package main

import (
	"fmt"
	"log"
	"net/http"

	"enstrurent.com/server/db"
	"enstrurent.com/server/flags"
	"enstrurent.com/server/routes"
)

// For local db change config con str as "CON_STR":"mongodb://localhost:27017/",
func main() {
	conf := flags.InitConfig()
	var port = conf.SERVER_PORT
	// var port = ":" + os.Getenv("PORT") // heroku version
	dbCon := db.OpenConnection(conf.CON_STR, conf.DBNAME)
	flags.InitCities(dbCon)
	server := &http.Server{
		Addr:    port,
		Handler: routes.Routes(dbCon),
	}
	fmt.Printf("Serving on %v\n", port)

	if err := server.ListenAndServe(); err != nil {
		log.Fatal(fmt.Sprintf("Error on listen and serve: %v", err))
	}
}

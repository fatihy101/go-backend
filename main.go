package main

import (
	"fmt"
	"log"
	"net/http"

	"enstrurent.com/server/db"
	"enstrurent.com/server/flags"
	"enstrurent.com/server/routes"
)

func main() {
	conf := flags.InitConfig()
	var port = conf.SERVER_PORT
	dbCon := db.OpenConnection(conf.CON_STR, conf.DBNAME)

	server := &http.Server{
		Addr:    port,
		Handler: routes.Routes(dbCon),
	}
	fmt.Printf("Serving on %v\n", port)

	if err := server.ListenAndServe(); err != nil {
		log.Fatal(fmt.Sprintf("Error on listen and serve: %v", err))
	}
}

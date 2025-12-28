package main

import (
	"backend/internal/config"
	"backend/internal/repository/sqlite"
	"log"
	"net/http"

	"backend/internal"
	"backend/internal/handler"
)

func main() {
	conf, err := config.Load()
	if err != nil {
		log.Fatalf("error loading config: %v", err)
	}

	database, err := sqlite.NewSQLiteDB(conf.Database)
	if err != nil {
		log.Fatalf("error opening repository: %v", err)
	}
	defer database.Close()

	handler := handler.NewAppHandler(database)
	app := internal.JSONMiddleware(handler)

	log.Println("Listening on :8080")
	if err := http.ListenAndServe(":8080", app); err != nil {
		log.Fatalf("error starting server: %v", err)
	}
}

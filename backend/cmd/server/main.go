package main

import (
	"backend/internal/repository/sqlite"
	"log"
	"net/http"

	"backend/internal"
	"backend/internal/handler"
)

func main() {
	database, err := sqlite.NewSQLiteDB()
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

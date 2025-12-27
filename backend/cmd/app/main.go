package main

import (
	"log"
	"net/http"

	"backend/internal"
)

func main() {
	db, err := internal.OpenDB()
	if err != nil {
		log.Fatalf("error opening db: %v", err)
	}
	defer db.Close()

	handler := internal.NewAppHandler(db)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /products", handler.GetProductsHandler)
	mux.HandleFunc("POST /products", handler.CreateProductHandler)
	mux.HandleFunc("POST /sales", handler.CreateSaleHandler)

	app := internal.JSONMiddleware(mux)

	log.Println("Listening on :8080")
	if err := http.ListenAndServe(":8080", app); err != nil {
		log.Fatalf("error starting server: %v", err)
	}
}

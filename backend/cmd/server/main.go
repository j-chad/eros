package main

import (
	"backend/internal/config"
	"backend/internal/repository/sqlite"
	"backend/internal/service"
	"fmt"
	"log"
	"net/http"

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

	authService := service.NewAuthService(conf.Auth, database)
	adminService := service.NewAdminService(database)

	app := handler.NewHandler(authService, adminService)

	serverAddr := fmt.Sprintf("%s:%d", conf.Server.Host, conf.Server.Port)
	server := &http.Server{
		Addr:         serverAddr,
		Handler:      app,
		ReadTimeout:  conf.Server.ReadTimeout,
		WriteTimeout: conf.Server.WriteTimeout,
		IdleTimeout:  conf.Server.IdleTimeout,
	}

	log.Println(fmt.Sprintf("Listening at %s", serverAddr))
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("error starting server: %v", err)
	}
}

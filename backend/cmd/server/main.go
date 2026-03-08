package main

import (
	"backend/internal"
	"backend/internal/config"
	"backend/internal/handler/middleware"
	"backend/internal/repository/sqlite"
	"backend/internal/repository/storage"
	"backend/internal/service"
	"fmt"
	"log"
	"log/slog"
	"net/http"

	"backend/internal/handler"
)

func main() {
	conf, err := config.Load()
	if err != nil {
		log.Fatalf("error loading config: %v", err)
	}

	internal.Init(conf.Logging)

	database, err := sqlite.NewSQLiteDB(conf.Database)
	if err != nil {
		log.Fatalf("error opening repository: %v", err)
	}
	defer database.Close()

	fileStore, err := storage.NewFileStore(conf.FileStorage)
	if err != nil {
		log.Fatalf("error initializing file storage: %v", err)
	}

	authService := service.NewAuthService(conf.Auth, database)
	adminService := service.NewAdminService(database, fileStore)
	favourService := service.NewFavourService(database)
	graphService := service.NewGraphService(database)
	fileService := service.NewFileService(database, fileStore)

	h := handler.NewHandler(authService, adminService, favourService, graphService, fileService)

	mux := http.NewServeMux()
	h.RegisterRoutes(mux)

	serverAddr := fmt.Sprintf("%s:%d", conf.Server.Host, conf.Server.Port)
	server := &http.Server{
		Addr:         serverAddr,
		Handler:      middleware.WithCORS(mux, conf.Server.CorsOrigins),
		ReadTimeout:  conf.Server.ReadTimeout,
		WriteTimeout: conf.Server.WriteTimeout,
		IdleTimeout:  conf.Server.IdleTimeout,
	}

	slog.Default().Info("starting server", "address", serverAddr)
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("error starting server: %v", err)
	}
}

package main

import (
	"backend/internal/config"
	"backend/internal/handler"
	"backend/internal/handler/middleware"
	"backend/internal/repository/sqlite"
	"backend/internal/repository/storage"
	"backend/internal/service"
	"fmt"
	"log"
	"log/slog"
	"net/http"
)

func runServer(conf *config.Config) {
	if conf == nil {
		log.Fatal("config is nil")
	}

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
	favourService := service.NewFavourService(database)
	fileService := service.NewFileService(database, fileStore)
	adminService := service.NewAdminService(database, fileStore, fileService)
	graphService := service.NewGraphService(database, fileService)

	h := handler.NewHandler(authService, adminService, favourService, graphService, fileService)

	mux := http.NewServeMux()
	h.RegisterRoutes(mux)

	serverAddr := fmt.Sprintf("%s:%d", conf.Server.Host, conf.Server.Port)
	server := &http.Server{
		Addr:         serverAddr,
		Handler:      middleware.WithLogging(middleware.WithCORS(mux, conf.Server.CorsOrigins)),
		ReadTimeout:  conf.Server.ReadTimeout,
		WriteTimeout: conf.Server.WriteTimeout,
		IdleTimeout:  conf.Server.IdleTimeout,
	}

	slog.Default().Info("starting server", "address", serverAddr)
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("error starting server: %v", err)
	}
}

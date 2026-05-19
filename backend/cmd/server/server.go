package main

import (
	"backend/internal/config"
	"backend/internal/handler"
	"backend/internal/handler/middleware"
	"backend/internal/repository/sqlite"
	"backend/internal/repository/storage"
	"backend/internal/scheduler"
	"backend/internal/service"
	"context"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func runServer(conf *config.Config) {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

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
	pushService := service.NewPushService(conf.Push, database)
	favourService := service.NewFavourService(database)
	fileService := service.NewFileService(database, fileStore)
	adminService := service.NewAdminService(database, fileStore, fileService)
	graphService := service.NewGraphService(database, fileService)

	if !conf.Scheduler.Disabled {
		sched, err := scheduler.New(conf.Scheduler, []scheduler.Task{
			{
				Name: "graph_notifications",
				Fn:   nil,
			},
		}...)
		if err != nil {
			log.Fatalf("error initializing scheduler: %v", err)
		}
		go sched.Run(ctx)
	}

	h := handler.NewHandler(authService, adminService, favourService, graphService, fileService, pushService)

	mux := http.NewServeMux()
	h.RegisterRoutes(mux)

	serverAddr := fmt.Sprintf("%s:%d", conf.Server.Host, conf.Server.Port)
	server := &http.Server{
		Addr:         serverAddr,
		Handler:      middleware.WithLogging(middleware.WithCTF(middleware.WithCORS(mux, conf.Server.CORS))),
		ReadTimeout:  conf.Server.ReadTimeout,
		WriteTimeout: conf.Server.WriteTimeout,
		IdleTimeout:  conf.Server.IdleTimeout,
	}

	go func() {
		slog.Default().Info("starting server", "address", serverAddr)
		if err := server.ListenAndServe(); err != nil {
			if errors.Is(err, http.ErrServerClosed) {
				return
			}
			
			log.Fatalf("error starting server: %v", err)
		}
	}()

	<-ctx.Done()
	log.Println("shutting down server...")
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()
	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("error shutting down server: %v", err)
	}
	log.Println("server gracefully stopped")
}

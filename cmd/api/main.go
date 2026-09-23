package main

import (
	"caronago/internal/modules/auth"
	"caronago/internal/modules/health"
	"caronago/internal/platform/config"
	"caronago/internal/platform/db"
	"caronago/internal/platform/middleware"
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg, err := config.Load()

	if err != nil {
		log.Fatalf("Configuration error: %v\n", err)
	}

	database, err := db.Connect(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Database connection error: %v", err)
	}

	log.Println("Connected to PostgreSQL with connection pooling")

	router := gin.New()

	router.Use(
		middleware.RequestID(),
		middleware.StructuredLogger(),
		gin.Recovery(),
	)

	router.HandleMethodNotAllowed = true
	router.NoRoute(middleware.NoRoute())
	router.NoMethod(middleware.NoMethod())

	health.RegisterRoutes(router)

	v1 := router.Group("/api/v1")
	auth.RegisterRoutes(v1, database, cfg)

	srv := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		log.Printf("🚀 Server running on http://localhost:%s", cfg.Port)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("Server failed to start: %v\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit
	log.Println("Shutting down server gracefully ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	if sqlDB, err := database.DB(); err == nil {
		if err := sqlDB.Close(); err != nil {
			log.Printf("Error closing database pool: %v", err)
		} else {
			log.Println("Database connection pool closed successfully")
		}
	}

	log.Println("Server exited cleanly 👋")

}

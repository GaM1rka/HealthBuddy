package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	httpHandler "auth_service/api/http"
	"auth_service/cmd/config"
	"auth_service/jwt"
	"auth_service/repository/db"
	authService "auth_service/usecases/service"

	"github.com/gin-gonic/gin"
)

func main() {
	// Load configuration
	dbCfg := config.LoadDBConfig()
	svcCfg := config.LoadServiceConfig()

	// Initialize JWT secret
	jwt.Init(svcCfg.Jwt_secret)

	// Initialize GORM DB
	gormDB, err := db.InitDB(dbCfg)
	if err != nil {
		log.Fatalf("failed to init DB: %v", err)
	}
	sqlDB, err := gormDB.DB()
	if err != nil {
		log.Fatalf("unable to get raw sql.DB: %v", err)
	}
	defer func() {
		if err := sqlDB.Close(); err != nil {
			log.Printf("error closing DB: %v", err)
		}
	}()

	// Verify DB connectivity
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := sqlDB.PingContext(ctx); err != nil {
		log.Fatalf("DB health check failed: %v", err)
	}
	log.Println("Database connection is healthy")

	// Set up Gin
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.Use(gin.Logger(), gin.Recovery())

	// Wire up repository, service, and handler
	userRepo := db.NewUserRepo(gormDB)
	svc := authService.NewAuthService(userRepo, svcCfg.Profile_service_utl)
	handler := httpHandler.NewAuthHandler(svc)
	handler.RegisterRoutes(router, svcCfg.ProfileServiceAuthToken)

	// Start HTTP server
	srv := &http.Server{
		Addr:    ":8083",
		Handler: router,
	}

	go func() {
		log.Printf("Auth service listening on %s", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen error: %v", err)
		}
	}()

	// Wait for shutdown signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit
	log.Println("Shutdown signal received, exiting...")

	// Graceful shutdown
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("Server shutdown failed: %v", err)
	}

	log.Println("Server exited cleanly")
}

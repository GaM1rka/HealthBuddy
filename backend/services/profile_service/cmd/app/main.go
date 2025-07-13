package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	handler "profile_service/api/http"
	"profile_service/cmd/config"
	"profile_service/repository/db"
	profileService "profile_service/usecases/service"

	"github.com/gin-gonic/gin"
)

func main() {
	// Config
	dbCfg := config.LoadDBConfig()
	svcCfg := config.LoafServiceCfg()

	// DB init
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

	// health db
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := sqlDB.PingContext(ctx); err != nil {
		log.Fatalf("DB health check failed: %v", err)
	}
	log.Println("Database connection is healthy")

	// gin
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.Use(gin.Logger(), gin.Recovery())

	repo := db.NewProfileRepo(gormDB)
	svc := profileService.NewProfileService(repo, svcCfg.Auth_service_url)
	h := handler.NewProfileHandler(svc)
	h.RegisterRoutes(router)

	srv := &http.Server{
		Addr:    ":8081",
		Handler: router,
	}
	go func() {
		log.Printf("Server started on %s", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %v", err)
		}
	}()

	// SIGINT/SIGTERM
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit
	log.Println("Shutdown signal received, exiting...")

	// Graceful shutdown
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("Server Shutdown failed: %v", err)
	}

	log.Println("Server exited cleanly")
}

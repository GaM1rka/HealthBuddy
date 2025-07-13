package main

import (
	"gateway_service/api/http"
	"gateway_service/cmd/config"
	"gateway_service/usecases/service"

	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	// Config load
	cfg := config.LoadConfig()

	// GatewayService
	svc := service.NewGatewayService(
		cfg.AuthServiceURL,
		cfg.ProfileServiceURL,
		cfg.FeedServiceURL,
		cfg.JWTSecret,
	)

	// Gin
	r := gin.Default()
	http.RegisterRoutes(r, svc)

	addr := ":8080"
	log.Printf("API Gateway listening on %s\n", addr)
	if err := r.Run(addr); err != nil {
		log.Fatalf("failed to run gateway: %v", err)
	}
}

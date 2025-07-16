package config

import (
	"log"
	"os"
)

type GatewayConfig struct {
	AuthServiceURL    string
	ProfileServiceURL string
	FeedServiceURL    string
	FrontURL          string
	JWTSecret         string
}

func LoadConfig() *GatewayConfig {
	mustEnv := func(key string) string {
		v := os.Getenv(key)
		if v == "" {
			log.Fatalf("environment variable %s is required", key)
		}
		return v
	}

	return &GatewayConfig{
		AuthServiceURL:    mustEnv("AUTH_SERVICE_URL"),
		ProfileServiceURL: mustEnv("PROFILE_SERVICE_URL"),
		FeedServiceURL:    mustEnv("FEED_SERVICE_URL"),
		FrontURL:          mustEnv("FRONT_URL"),
		JWTSecret:         mustEnv("JWT_SECRET"),
	}
}

package service

import (
	"gateway_service/usecases/helpers"

	"github.com/gin-gonic/gin"
)

// GatewayService
type GatewayService struct {
	AuthURL    string
	ProfileURL string
	FeedURL    string
	JWTSecret  string
}

// NewGatewayService
func NewGatewayService(authURL, profileURL, feedURL, jwtSecret string) *GatewayService {
	return &GatewayService{
		AuthURL:    authURL,
		ProfileURL: profileURL,
		FeedURL:    feedURL,
		JWTSecret:  jwtSecret,
	}
}

// Handlers

func (g *GatewayService) AuthProxy() gin.HandlerFunc {
	return helpers.ReverseProxy(g.AuthURL)
}

func (g *GatewayService) ProfileProxy() gin.HandlerFunc {
	return helpers.ReverseProxy(g.ProfileURL)
}

func (g *GatewayService) FeedProxy() gin.HandlerFunc {
	return helpers.ReverseProxy(g.FeedURL)
}

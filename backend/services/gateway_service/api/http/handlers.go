package http

import (
	"gateway_service/usecases/middleware"
	"gateway_service/usecases/service"

	"github.com/gin-gonic/gin"
)

// RegisterRoutes
func RegisterRoutes(r *gin.Engine, svc *service.GatewayService) {
	// 1) public auth
	auth := r.Group("/auth")
	{
		auth.GET("/health", svc.AuthProxy())
		auth.POST("/register", svc.AuthProxy())
		auth.POST("/login", svc.AuthProxy())
	}

	// 2) protected JWT-middleware
	protected := r.Group("/", middleware.JWTMiddleware(svc.JWTSecret))
	{
		protected.Any("/profile/*proxyPath", svc.ProfileProxy())
		protected.Any("/feed/*proxyPath", svc.FeedProxy())
	}
}

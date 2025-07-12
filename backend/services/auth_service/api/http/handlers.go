package http

import (
	"context"
	"net/http"

	"auth_service/domain"
	"auth_service/usecases"

	"github.com/gin-gonic/gin"
)

// AuthHandler handles HTTP requests for user registration, login and health check.
type AuthHandler struct {
	svc usecases.AuthService
}

// NewAuthHandler constructs a new AuthHandler.
func NewAuthHandler(svc usecases.AuthService) *AuthHandler {
	return &AuthHandler{svc: svc}
}

// RegisterRoutes registers auth routes on the Gin engine.
func (h *AuthHandler) RegisterRoutes(r *gin.Engine) {
	grp := r.Group("/auth")
	{
		grp.GET("/health", h.Health) // ← добавили health
		grp.POST("/register", h.Register)
		grp.POST("/login", h.Login)
	}
}

// Health handles GET /auth/health
func (h *AuthHandler) Health(c *gin.Context) {
	if healthSvc, ok := h.svc.(interface {
		Health(ctx context.Context) error
	}); ok {
		if err := healthSvc.Health(c.Request.Context()); err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{"status": "down"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
		return
	}
}

// Register handles POST /auth/register
func (h *AuthHandler) Register(c *gin.Context) {
	var creds domain.Credentials
	if err := c.ShouldBindJSON(&creds); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := creds.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := h.svc.Register(c.Request.Context(), creds)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"token": token})
}

// Login handles POST /auth/login
func (h *AuthHandler) Login(c *gin.Context) {
	var creds domain.Credentials
	if err := c.ShouldBindJSON(&creds); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := creds.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := h.svc.Login(c.Request.Context(), creds)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}

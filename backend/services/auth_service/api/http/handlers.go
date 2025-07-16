package http

import (
	"context"
	"net/http"

	"auth_service/api/http/apierrors"
	"auth_service/api/http/middleware"
	"auth_service/domain"
	"auth_service/usecases"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	svc usecases.AuthService
}

func NewAuthHandler(svc usecases.AuthService) *AuthHandler {
	return &AuthHandler{svc: svc}
}

func (h *AuthHandler) RegisterRoutes(r *gin.Engine) {
	grp := r.Group("/auth")
	{
		// оборачиваем каждый метод в наш middleware
		grp.GET("/health", middleware.ErrorHandlerMiddleware(h.Health))
		grp.POST("/register", middleware.ErrorHandlerMiddleware(h.Register))
		grp.POST("/login", middleware.ErrorHandlerMiddleware(h.Login))
		grp.GET("/user/:id", middleware.ErrorHandlerMiddleware(h.GetUserByID))
		grp.DELETE("/user/:id", middleware.ErrorHandlerMiddleware(h.DeleteUser))
	}
}

// Health
func (h *AuthHandler) Health(c *gin.Context) error {
	if healthSvc, ok := h.svc.(interface {
		Health(ctx context.Context) error
	}); ok {
		if err := healthSvc.Health(c.Request.Context()); err != nil {
			return apierrors.NewInternal(err)
		}
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
		return nil
	}
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
	return nil
}

func (h *AuthHandler) Register(c *gin.Context) error {
	var creds domain.RegistrCredentials
	if err := c.ShouldBindJSON(&creds); err != nil {
		return apierrors.NewBadRequest("invalid JSON", err)
	}
	if err := creds.Validate(); err != nil {
		return apierrors.NewBadRequest(err.Error(), err)
	}

	token, err := h.svc.Register(c.Request.Context(), creds)
	if err != nil {
		return apierrors.NewInternal(err)
	}

	c.JSON(http.StatusCreated, gin.H{"token": token})
	return nil
}

// Login JWT
func (h *AuthHandler) Login(c *gin.Context) error {
	var creds domain.LoginCredentials
	if err := c.ShouldBindJSON(&creds); err != nil {
		return apierrors.NewBadRequest("invalid JSON", err)
	}
	if err := creds.Validate(); err != nil {
		return apierrors.NewBadRequest(err.Error(), err)
	}

	token, err := h.svc.Login(c.Request.Context(), creds)
	if err != nil {
		return apierrors.NewForbidden(err.Error())
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
	return nil
}

// GetUserByID by ID
func (h *AuthHandler) GetUserByID(c *gin.Context) error {
	id := c.Param("id")
	user, err := h.svc.FindByID(c.Request.Context(), id)
	if err != nil {
		if err == usecases.ErrNotFound {
			return apierrors.NewNotFound("user not found")
		}
		return apierrors.NewInternal(err)
	}

	c.JSON(http.StatusOK, user)
	return nil
}

// DeleteUser by ID
func (h *AuthHandler) DeleteUser(c *gin.Context) error {
	id := c.Param("id")
	if err := h.svc.DeleteUser(c.Request.Context(), id); err != nil {
		if err == usecases.ErrNotFound {
			return apierrors.NewNotFound("user not found")
		}
		return apierrors.NewInternal(err)
	}

	c.Status(http.StatusNoContent)
	return nil
}

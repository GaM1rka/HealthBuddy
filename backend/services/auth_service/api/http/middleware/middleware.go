package middleware

import (
	"auth_service/api/http/apierrors"
	"auth_service/usecases"
	"errors"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type HandlerFunc func(c *gin.Context) error

func ErrorHandlerMiddleware(h HandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := h(c); err != nil {
			switch {
			case errors.Is(err, usecases.ErrEmailTaken):
				c.JSON(http.StatusConflict, gin.H{"error": "email or username already in use"})
			case errors.Is(err, usecases.ErrInvalidCredentials):
				c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
			case errors.Is(err, usecases.ErrUserNotFound):
				c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
			case errors.Is(err, usecases.ErrProfileServiceDown):
				c.JSON(http.StatusBadGateway, gin.H{"error": "profile service unavailable"})
			default:
				var apiErr apierrors.APIError
				if errors.As(err, &apiErr) {
					c.JSON(apiErr.Code, gin.H{"error": apiErr.Message})
				} else {
					log.Printf("unexpected error: %v\n", err)
					c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
				}
			}
			c.Abort()
		}
	}
}

func ServiceAuthMiddleware(expectedToken string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.GetHeader("X-Service-Token") != expectedToken {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "service unauthorized"})
			c.Abort()
			return
		}
		c.Next()
	}
}

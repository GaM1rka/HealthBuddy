package middleware

import (
	"auth_service/api/http/apierrors"
	"errors"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type HandlerFunc func(c *gin.Context) error

func ErrorHandlerMiddleware(h HandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := h(c); err != nil {
			var apiErr apierrors.APIError
			if ok := errors.As(err, &apiErr); ok {
				c.JSON(apiErr.Code, gin.H{"error": apiErr.Message})
			} else {
				log.Printf("unexpected error: %v\n", err)
				c.JSON(500, gin.H{"error": "internal server error"})
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

package middleware

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

func JWTMiddleware(secret, authServiceURL string) gin.HandlerFunc {
	// создаём клиент с небольшим таймаутом
	client := &http.Client{Timeout: 2 * time.Second}

	return func(c *gin.Context) {
		// 1) Парсим и верифицируем JWT
		authHeader := c.GetHeader("Authorization")
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing or malformed Authorization header"})
			return
		}

		token, err := jwt.ParseWithClaims(parts[1], &jwt.RegisteredClaims{}, func(t *jwt.Token) (interface{}, error) {
			return []byte(secret), nil
		})
		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			return
		}

		claims, ok := token.Claims.(*jwt.RegisteredClaims)
		if !ok || claims.Subject == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "token contains no subject"})
			return
		}
		userID := claims.Subject
		authURL := fmt.Sprintf("%s/auth/user/%s", authServiceURL, userID)
		req, err := http.NewRequestWithContext(c.Request.Context(), http.MethodGet, authURL, nil)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
			return
		}

		resp, err := client.Do(req)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "cannot verify user"})
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode == http.StatusNotFound {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "user not found"})
			return
		}
		if resp.StatusCode != http.StatusOK {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "user verification failed"})
			return
		}

		// 3) Всё ок, прокидываем userID дальше
		c.Request.Header.Set("X-User-ID", userID)
		c.Next()
	}
}

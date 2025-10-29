package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/example/golang-rest-boilerplate/internal/service"
	"github.com/example/golang-rest-boilerplate/pkg/response"
)

const userClaimsKey = "userClaims"

// AuthMiddleware validates JWT tokens and attaches claims to the request context.
func AuthMiddleware(authService *service.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.Error(c, http.StatusUnauthorized, "authorization header missing")
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			response.Error(c, http.StatusUnauthorized, "invalid authorization header")
			return
		}

		claims, err := authService.ParseToken(parts[1])
		if err != nil {
			response.Error(c, http.StatusUnauthorized, "invalid token")
			return
		}

		c.Set(userClaimsKey, claims)
		c.Next()
	}
}

// GetClaims extracts JWT claims from the context.
func GetClaims(c *gin.Context) *service.Claims {
	value, exists := c.Get(userClaimsKey)
	if !exists {
		return nil
	}
	if claims, ok := value.(*service.Claims); ok {
		return claims
	}
	return nil
}

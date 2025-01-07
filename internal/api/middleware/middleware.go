package middleware

import (
	"errors"
	"mahi-go-explorer/internal/api/response"
	"mahi-go-explorer/internal/config"
	userpkg "mahi-go-explorer/pkg/user"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

// Authenticate checks if user is authenticated
func Authenticate() gin.HandlerFunc {
	return func(c *gin.Context) {
		//get auth header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.LogAndErrorResponse(c, http.StatusUnauthorized, "Authorization header reqired", errors.New("Auth header required"))
			c.Abort()
			return
		}

		//get token part
		tokenParts := strings.Split(authHeader, "Bearer ")
		if len(tokenParts) != 2 {
			response.LogAndErrorResponse(c, http.StatusUnauthorized, "Invalid authorization header", errors.New("Invalid authorization header"))
			c.Abort()
			return
		}

		//verify by jwt secret
		secret := config.GetFromEnv("JWT_SECRET")

		token, err := jwt.ParseWithClaims(tokenParts[1], &userpkg.JwtClaims{}, func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("Invalid signing method")
			}
			return []byte(secret), nil
		})
		if err != nil {
			response.LogAndErrorResponse(c, http.StatusUnauthorized, "Invalid token", err)
			c.Abort()
			return
		}

		if !token.Valid {
			response.LogAndErrorResponse(c, http.StatusUnauthorized, "Invalid token", errors.New("Invalid token"))
			c.Abort()
			return
		}

		claims, ok := token.Claims.(*userpkg.JwtClaims)
		if !ok {
			response.LogAndErrorResponse(c, http.StatusUnauthorized, "Invalid token claims", errors.New("Invalid token claims"))
			c.Abort()
			return
		}

		//check token expiration
		if claims.ExpiresAt > jwt.TimeFunc().Unix() {
			response.LogAndErrorResponse(c, http.StatusUnauthorized, "Token expired", errors.New("Token expired"))
			c.Abort()
			return
		}

		user := &userpkg.UserContext{
			ID:        claims.ID,
			FirstName: claims.FirstName,
			LastName:  claims.LastName,
			Email:     claims.Email,
			Role:      claims.Role,
			Exp:       claims.ExpiresAt,
		}

		c.Set("user", user)

		c.Next()
	}
}

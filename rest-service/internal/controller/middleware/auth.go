package middleware

import (
	"calls-service/rest-service/internal/controller/apierrors"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, apierrors.Response{Error: "Authorization header is missing"})
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, apierrors.Response{Error: "Invalid authorization header format"})
			return
		}

		tokenStr := parts[1]

		token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (any, error) {
			if token.Method.Alg() != jwt.SigningMethodHS256.Alg() {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(os.Getenv("JWT_SECRET")), nil
		})
		if err != nil {
			if err == jwt.ErrSignatureInvalid {
				c.AbortWithStatusJSON(http.StatusUnauthorized, apierrors.Response{Error: "Invalid token signature"})
				return
			}
			c.AbortWithStatusJSON(http.StatusUnauthorized, apierrors.Response{Error: "Invalid token"})
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			userID, ok := claims["id"].(float64)
			if !ok {
				c.AbortWithStatusJSON(http.StatusUnauthorized, apierrors.Response{Error: "Invalid token claims"})
				return
			}
			c.Set("id", int64(userID))
		} else {
			c.AbortWithStatusJSON(http.StatusUnauthorized, apierrors.Response{Error: "Invalid token"})
			return
		}
		c.Next()
	}
}

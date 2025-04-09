package middleware

import (
	"money-transfer-demo/pkg/identity/token"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware(accountType token.AccountType) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing token"})
			c.Abort()
			return
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
		claims, err := token.ValidateJWT(tokenStr, accountType)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		c.Set("member_id", claims.UserID)
		c.Next()
	}
}

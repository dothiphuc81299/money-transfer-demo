package middleware

import (
	"context"
	"money-transfer-demo/pkg/identity/token"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware(accountType token.AccountType) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing token"})
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		authToken := strings.Split(authHeader, " ")
		if len(authToken) != 2 || authToken[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token format"})
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
		claims, err := token.ValidateJWT(tokenStr, accountType)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		expirationTime, err := claims.GetExpirationTime()
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token expiration"})
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		if float64(time.Now().UTC().Unix()) > float64(expirationTime.Unix()) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "token expired"})
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		var account = &token.AccountData{
			LoginName:   claims.LoginName,
			ID:          claims.UserID,
			AccountType: accountType,
		}
		
		ctx := context.WithValue(c.Request.Context(), "current_account", account)
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}

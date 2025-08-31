package middleware

import (
	"github.com/Citadelas/api-gateway/internal/lib/jwt"
	ssov1 "github.com/Citadelas/protos/golang/sso"
	"github.com/gin-gonic/gin"
)

func AuthMiddleware(ssoClient ssov1.AuthClient) gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		token := jwt.ExtractToken(c.GetHeader("Authorization"))

		userID, err := jwt.ValidateToken(c.Request.Context(), ssoClient, token)
		if err != nil {
			c.JSON(401, gin.H{
				"error":   "unauthorized",
				"details": err.Error(),
			})
			c.Abort()
			return
		}

		c.Set("userID", userID)
		c.Next()
	})
}

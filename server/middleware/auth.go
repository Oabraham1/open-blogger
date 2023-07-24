package middleware

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func IsAuthenticated(ctx *gin.Context) {
	if ctx.GetHeader("X-API-KEY") != os.Getenv("API_KEY") {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	ctx.Next()
}
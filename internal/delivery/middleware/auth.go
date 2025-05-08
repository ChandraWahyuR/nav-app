package middleware

import (
	"net/http"
	"proyek1/internal/model"
	jwt "proyek1/utils"
	"strings"

	"github.com/gin-gonic/gin"
)

func NewAuth(jwt jwt.JWTInterface) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorize auth"}) // gin.H sama aja kaya map[string]string{"": ""}
			return
		}

		tokenParts := strings.Split(authHeader, " ") // pisah jadi dibuat ada 2 index 0 & 1. contoh "Bearer ......"
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "token tidak invalid, format token berbeda"})
			return
		}

		token := tokenParts[1]
		userData, err := jwt.VerifyToken(token)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"errors": err.Error()})
			return
		}
		ctx.Set("auth", userData)
		ctx.Next()
	}
}

// Bisa ambil data user dari token jwt
func GetUser(ctx *gin.Context) (*model.User, bool) {
	userData, exists := ctx.Get("auth")
	if !exists {
		return nil, false
	}
	user, ok := userData.(*model.User)
	if !ok {
		return nil, false
	}
	return user, true
}

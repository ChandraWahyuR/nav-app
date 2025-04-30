package utils

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func ConverResponse(err error) int {
	log.Printf("Received error: %v", err)
	switch err {
	case ErrGetData:
		return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}
}
func HandleEchoError(err error) (int, string) {
	if _, ok := err.(*gin.Error); ok {
		return http.StatusBadRequest, BadInput
	}
	return http.StatusBadRequest, BadInput
}

func UnauthorizedError(c *gin.Context) {
	c.JSON(http.StatusUnauthorized, gin.H{"message": Unauthorized})
}

func InternalServerError(c *gin.Context) {
	c.JSON(http.StatusInternalServerError, gin.H{"message": InternalServer})
}

func JWTErrorHandler(c *gin.Context, err error) {
	c.JSON(http.StatusUnauthorized, gin.H{"message": InternalServer})
}

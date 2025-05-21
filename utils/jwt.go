package utils

import (
	"fmt"
	"proyek1/internal/model"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWT struct {
	signKey string
}

type JWTInterface interface {
	GenerateToken(data *model.User) (string, error)
	VerifyToken(tokenString string) (*model.User, error)
}

func NewJWT(s string) JWTInterface {
	return &JWT{
		signKey: s,
	}
}

func (j *JWT) GenerateToken(data *model.User) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":    data.ID,
		"email": data.Email,
		"role":  data.Role,
		"exp":   time.Now().Add(time.Hour * 24 * 30).Unix(),
		"iat":   time.Now().Unix(), // tgl dibuat jwt ini
	})

	tokenString, err := token.SignedString([]byte(j.signKey))
	if err != nil {
		return "", fmt.Errorf(`terjadi kesalahan :%s`, err)
	}

	return tokenString, nil
}

func (j *JWT) VerifyToken(tokenString string) (*model.User, error) {
	tokenParse, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		return []byte(j.signKey), nil
	})
	if err != nil {
		return nil, fmt.Errorf(`terjadi kesalahan: %s`, err)
	}

	if !tokenParse.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	claims, ok := tokenParse.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("invalid claims")
	}

	return &model.User{
		ID:    claims["id"].(string),
		Email: claims["email"].(string),
		Role:  claims["role"].(string),
	}, nil
}

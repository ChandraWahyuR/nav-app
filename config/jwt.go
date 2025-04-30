package config

import (
	"os"
	"proyek1/utils"

	"github.com/sirupsen/logrus"
)

func NewJWT(log *logrus.Logger) utils.JWTInterface {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		log.Error("Gagal jwt secret")
	}
	return utils.NewJWT(secret)
}

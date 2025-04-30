package app

import (
	"database/sql"
	"proyek1/config"
	"proyek1/internal/delivery"
	"proyek1/internal/delivery/routes"
	"proyek1/internal/repository"
	"proyek1/internal/usecase"

	"proyek1/utils"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type BootstrapConfig struct {
	DB  *sql.DB
	App *gin.Engine
	Log *logrus.Logger
	JWT utils.JWTInterface
	Cfg *config.Config
}

func App(config *BootstrapConfig) {
	// Repository
	userRepository := repository.NewUserRepository(config.DB, config.Log)

	// UseCase
	userUsecase := usecase.NewUserUsecase(config.JWT, userRepository, config.Log, config.Cfg)

	// Delivery
	userHandler := delivery.NewUserHandler(config.JWT, userUsecase, config.Log)

	routeConfig := routes.RouteConfig{
		App:            config.App,
		UserController: userHandler,
		JWT:            config.JWT,
	}

	routeConfig.Setup()
}

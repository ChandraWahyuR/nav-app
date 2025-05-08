package app

import (
	"database/sql"
	"proyek1/config"
	"proyek1/internal/delivery"
	"proyek1/internal/delivery/routes"
	"proyek1/internal/repository"
	"proyek1/internal/usecase"
	"proyek1/utils/gmaps"
	"proyek1/utils/mailer"

	"proyek1/utils"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type BootstrapConfig struct {
	DB   *sql.DB
	App  *gin.Engine
	Log  *logrus.Logger
	JWT  utils.JWTInterface
	Cfg  *config.Config
	M    mailer.MailInterface
	Maps gmaps.GmapsInterface
}

func App(config *BootstrapConfig) {
	// Repository
	userRepository := repository.NewUserRepository(config.DB, config.Log)
	mapsRepository := repository.NewMapsRepository(config.DB, config.Log)

	// UseCase
	userUsecase := usecase.NewUserUsecase(config.JWT, userRepository, config.Log, config.Cfg, config.M)
	mapsUsecase := usecase.NewMapsUsercase(mapsRepository, config.Log, config.Maps)
	// Delivery
	userHandler := delivery.NewUserHandler(config.JWT, userUsecase, config.Log)
	mapsHandler := delivery.NewMapsHandler(config.JWT, config.Maps, mapsUsecase)

	routeConfig := routes.RouteConfig{
		App:            config.App,
		UserController: userHandler,
		MapsController: &mapsHandler,
		JWT:            config.JWT,
	}

	routeConfig.Setup()
}

package routes

import (
	"proyek1/internal/delivery"
	"proyek1/internal/delivery/middleware"
	"proyek1/utils"

	"github.com/gin-gonic/gin"
)

type RouteConfig struct {
	App            *gin.Engine
	UserController *delivery.UserHandler
	MapsController *delivery.MapsHandler
	JWT            utils.JWTInterface
}

func (c *RouteConfig) Setup() {
	c.SetupUserRoute()
	c.SetupMapsRoute()
}

func (c *RouteConfig) SetupUserRoute() {
	c.App.POST("/register", c.UserController.Register)
	c.App.POST("/login", c.UserController.Login)

	c.App.Use(middleware.NewAuth(c.JWT))
	c.App.GET("/profile", c.UserController.Profile)
	c.App.PUT("/profile", c.UserController.EditProfile)
}

func (c *RouteConfig) SetupMapsRoute() {

	c.App.Use(middleware.NewAuth(c.JWT))
	c.App.GET("/maps", c.MapsController.GmapsSearchbyObject)
	c.App.GET("/maps-list", c.MapsController.GmapsSearchbyList)
	c.App.GET("/place/:id", c.MapsController.GmapsSearchbyPlaceID)
	c.App.GET("/photo", c.MapsController.ProxyPhotoHandler)
}

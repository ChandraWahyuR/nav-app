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
	JWT            utils.JWTInterface
}

func (c *RouteConfig) Setup() {
	c.SetupUserRoute()
}

func (c *RouteConfig) SetupUserRoute() {
	c.App.POST("/register", c.UserController.Register)
	c.App.POST("/login", c.UserController.Login)

	c.App.Use(middleware.NewAuth(c.JWT))
	c.App.GET("/profile", c.UserController.Profile)
	c.App.PUT("/profile", c.UserController.EditProfile)
}

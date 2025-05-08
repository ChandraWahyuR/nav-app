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
	c.App.POST("/reg-admin", c.UserController.RegisterForAdmin)

	private := c.App.Group("/")
	private.Use(middleware.NewAuth(c.JWT))
	private.GET("/profile", c.UserController.Profile)
	private.PUT("/profile", c.UserController.EditProfile)
}

func (c *RouteConfig) SetupMapsRoute() {
	c.App.GET("/photo", c.MapsController.ProxyPhotoHandler) // app use global, nanti kena semua

	private := c.App.Group("/")
	private.Use(middleware.NewAuth(c.JWT))
	private.GET("/maps", c.MapsController.GmapsSearchbyObject)
	private.GET("/maps-list", c.MapsController.GmapsSearchbyList)
	private.GET("/place", c.MapsController.GetTempatPagination)
	private.GET("/place/:id", c.MapsController.GmapsSearchbyPlaceID)
	private.POST("/place/:id", c.MapsController.InsertData)

	private.POST("/route-maps/:id", c.MapsController.RouteDestination)
}

package delivery

import (
	"context"
	"fmt"
	"net/http"
	"proyek1/internal/delivery/middleware"
	"proyek1/internal/model"
	jwt "proyek1/utils"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type RegisterUsecaseInterface interface {
	Register(ctx context.Context, postData *model.User) error
	Login(ctx context.Context, postData *model.Login) (*model.Login, error)
	Profile(ctx context.Context, id string) (model.User, error)
	EditProfile(ctx context.Context, data *model.EditProfile, id string) error
}

type RegisterHandlerInterface interface {
	Register(c *gin.Context)
	Login(c *gin.Context)
	Profile(c *gin.Context)
	EditProfile(c *gin.Context)
}

type UserHandler struct {
	jwt jwt.JWTInterface
	uc  RegisterUsecaseInterface
	log *logrus.Logger
}

func NewUserHandler(j jwt.JWTInterface, uc RegisterUsecaseInterface, log *logrus.Logger) *UserHandler {
	return &UserHandler{
		jwt: j,
		uc:  uc,
		log: log,
	}
}

func (h *UserHandler) Register(c *gin.Context) {
	defer func() {
		if r := recover(); r != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Internal server error: %v", r)})
		}
	}()
	ctx := c.Request.Context()
	var data model.Register
	if err := c.Bind(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	modelData := model.User{
		Username:        data.Username,
		Email:           data.Email,
		Password:        data.Password,
		ConfirmPassword: data.ConfirmPassword,
	}
	err := h.uc.Register(ctx, &modelData)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Registrasi berhasil"})
}

func (h *UserHandler) Login(c *gin.Context) {
	var data model.Login
	if err := c.Bind(&data); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	ctx := c.Request.Context()
	result, err := h.uc.Login(ctx, &data)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res := result.Token
	c.JSON(http.StatusOK, gin.H{"message": "Berhasil login", "token": res})
}

func (h *UserHandler) Profile(c *gin.Context) {
	dataToken, ok := middleware.GetUser(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "unatuhorize"})
		return
	}

	ctx := c.Request.Context()
	userData, err := h.uc.Profile(ctx, dataToken.ID)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "error proses id"})
		return
	}

	resData := model.Profile{
		ID:           userData.ID,
		Username:     userData.Username,
		Email:        userData.Email,
		PhotoProfile: userData.PhotoProfile,
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Profile success",
		"data":    resData,
	})
}

func (h *UserHandler) EditProfile(c *gin.Context) {
	defer func() {
		if r := recover(); r != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Internal server error: %v", r)})
		}
	}()
	dataToken, ok := middleware.GetUser(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var data model.EditProfile
	if err := c.Bind(&data); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}
	dataUser := &model.EditProfile{
		ID:           data.ID,
		Username:     data.Username,
		PhotoProfile: data.Username,
		Password:     data.Password,
	}
	ctx := c.Request.Context()
	err := h.uc.EditProfile(ctx, dataUser, dataToken.ID)
	if err != nil {
		fmt.Println("err :", err)
		c.JSON(http.StatusBadRequest, gin.H{"message": "error proses data"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "data berhasil di edit"})
}

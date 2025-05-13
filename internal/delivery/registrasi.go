package delivery

import (
	"context"
	"fmt"
	"net/http"
	"proyek1/constant"
	"proyek1/internal/delivery/middleware"
	"proyek1/internal/model"
	crypto "proyek1/utils"
	jwt "proyek1/utils"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type RegisterUsecaseInterface interface {
	Register(ctx context.Context, postData *model.User) error
	Login(ctx context.Context, postData *model.Login) (*model.Login, error)
	Profile(ctx context.Context, id string) (model.User, error)
	EditProfile(ctx context.Context, data *model.EditProfile, id string) error

	RegisterForAdmin(ctx context.Context, req *model.User) error

	ForgotPassword(ctx context.Context, req *model.Otp) error
	OtpVerify(ctx context.Context, req *model.Otp) (*model.Otp, error)
	ResetPassword(ctx context.Context, req *model.User) error
	ActivateAcount(ctx context.Context, email string) error
}

type RegisterHandlerInterface interface {
	Register(c *gin.Context)
	Login(c *gin.Context)
	Profile(c *gin.Context)
	EditProfile(c *gin.Context)
	RegisterForAdmin(c *gin.Context)
	ForgotPassword(c *gin.Context)
	OtpVerify(c *gin.Context)
	ResetPassword(c *gin.Context)
	ActivateAcount(c *gin.Context)
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

	if crypto.IsUser(dataToken.Role) {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "role apa njir"})
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

	if crypto.IsUser(dataToken.Role) {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "role apa njir"})
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

func (h *UserHandler) RegisterForAdmin(c *gin.Context) {
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

	err := h.uc.RegisterForAdmin(ctx, &modelData)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Registrasi admin berhasil"})
}

func (h *UserHandler) ForgotPassword(c *gin.Context) {
	ctx := c.Request.Context()

	var data model.Otp
	if err := c.Bind(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	modelData := model.Otp{
		Email: data.Email,
	}

	err := h.uc.ForgotPassword(ctx, &modelData)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Otp berhasil dibuat, silahkan cek email"})
}

func (h *UserHandler) OtpVerify(c *gin.Context) {
	ctx := c.Request.Context()

	var data model.Otp
	if err := c.Bind(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	modelData := model.Otp{
		ID:        data.ID,
		Email:     data.Email,
		OtpNumber: data.OtpNumber,
	}

	res, err := h.uc.OtpVerify(ctx, &modelData)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "otp berhasil diverifikasi", "token": res.Token})
}

func (h *UserHandler) ResetPassword(c *gin.Context) {
	ctx := c.Request.Context()
	dataToken, ok := middleware.GetUser(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "unatuhorize"})
		return
	}

	if !crypto.IsForgot(dataToken.Role) {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "role apa njir"})
		return
	}

	var data model.User
	if err := c.Bind(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	modelData := model.User{
		Email:           dataToken.Email,
		Password:        data.Password,
		ConfirmPassword: data.ConfirmPassword,
	}

	err := h.uc.ResetPassword(ctx, &modelData)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "reset password berhasil"})
}

func (h *UserHandler) ActivateAcount(c *gin.Context) {
	token := c.Query("token")
	if token == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "token is required"})
		return
	}

	userData, err := h.jwt.VerifyToken(token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
		return
	}

	if !crypto.IsUser(userData.Role) {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "role tidak diizinkan"})
		return
	}

	ctx := c.Request.Context()
	err = h.uc.ActivateAcount(ctx, userData.Email)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.Redirect(http.StatusSeeOther, constant.VercelRoute) // 303
}

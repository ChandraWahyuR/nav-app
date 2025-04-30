package usecase

import (
	"context"
	"errors"
	"fmt"
	"os"
	"proyek1/config"
	"proyek1/internal/entity"
	"proyek1/internal/model"
	"proyek1/utils"
	jwt "proyek1/utils"
	"strings"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

var (
	DefaultPhoto = os.Getenv("DEFAULT_PP") // sementara
	RoleUser     = "user"
)

type RepositoryUserInterface interface {
	Register(ctx context.Context, postData *entity.User) error
	Login(ctx context.Context, postData *entity.User) (*entity.User, error)
	GetUserID(ctx context.Context, id string) (entity.User, error)
	ForgotPassword(ctx context.Context, data *entity.Otp) error // Nanti
	OtpVerify(ctx context.Context, data *entity.Otp) (*entity.Otp, error)
	ResetPassword(ctx context.Context, data *entity.User) error

	//
	IsEmailAvailable(ctx context.Context, email string) bool
	IsUsernameAvailable(ctx context.Context, username string) bool
}

type UsecaseUser struct {
	jwt      jwt.JWTInterface
	userRepo RepositoryUserInterface
	log      *logrus.Logger
	cfg      *config.Config
}

func NewUserUsecase(jwt jwt.JWTInterface, user RepositoryUserInterface, log *logrus.Logger, cfg *config.Config) *UsecaseUser {
	return &UsecaseUser{
		jwt:      jwt,
		userRepo: user,
		log:      log,
		cfg:      cfg,
	}
}

func (s *UsecaseUser) Register(ctx context.Context, userData *model.User) error {
	switch {
	case userData.Username == "":
		return errors.New("username tidak boleh kosong")
	case userData.Email == "":
		return errors.New("email tidak boleh kosong")
	case userData.Password == "":
		return errors.New("password tidak boleh kosong")
	case userData.Password != userData.ConfirmPassword:
		return errors.New("konfirmasi password tidak sama dengan password")
	}

	// Cek format password
	if !utils.ValidasiEmail(userData.Email) {
		return errors.New("format email tidak benar")
	}

	// Cek ketersediaan email
	if s.userRepo.IsEmailAvailable(ctx, userData.Email) {
		return errors.New("email sudah digunakan")
	}

	// Cek ketersediaan username
	if s.userRepo.IsUsernameAvailable(ctx, userData.Username) {
		return errors.New("username sudah digunakan")
	}

	// Cek validasi password
	pass, err := utils.ValidatePassword(userData.Password)
	if err != nil {
		return fmt.Errorf("password tidak valid: %w", err)
	}

	// Proses hashing password
	hashedPassword, err := utils.HashPassword(pass)
	if err != nil {
		return err
	}

	resData := entity.User{
		ID:           uuid.New().String(),
		Username:     userData.Username,
		Email:        userData.Email,
		Password:     hashedPassword,
		PhotoProfile: s.cfg.GeneralPhoto.DefaultPhoto,
		Role:         RoleUser,
		IsActive:     true, // sementara
	}

	err = s.userRepo.Register(ctx, &resData)
	if err != nil {
		return err
	}

	return nil
}

func (s *UsecaseUser) Login(ctx context.Context, postData *model.User) (*model.User, error) {
	switch {
	case postData.Email == "":
		return &model.User{}, errors.New("email tidak boleh kosong")
	case postData.Password == "":
		return &model.User{}, errors.New("password tidak boleh kosong")
	}
	// Cek format password
	if !utils.ValidasiEmail(postData.Email) {
		return &model.User{}, errors.New("format email tidak benar")
	}

	postData.Email = strings.ToLower(postData.Email)
	userEntity := entity.User{
		Email:    postData.Email,
		Password: postData.Password,
	}

	res, err := s.userRepo.Login(ctx, &userEntity)
	if err != nil {
		return &model.User{}, err
	}

	if postData.ID == "" {
		return &model.User{}, fmt.Errorf("id user tidak ditemukan")
	}

	tokenData := model.User{
		ID:       res.ID,
		Username: res.Username,
		Email:    res.Email,
		Role:     RoleUser,
	}
	token, err := s.jwt.GenerateToken(&tokenData)
	if err != nil {
		return &model.User{}, err
	}

	UserLoginData := &model.User{}
	UserLoginData.Token = token

	return UserLoginData, nil
}

func (s *UsecaseUser) Profile(ctx context.Context, id string) (model.User, error) {
	if id == "" {
		return model.User{}, errors.New("id kosong")
	}

	resData, err := s.userRepo.GetUserID(ctx, id)
	if err != nil {
		return model.User{}, err
	}

	entityData := model.User{
		ID:           resData.ID,
		Username:     resData.Username,
		Email:        resData.Email,
		PhotoProfile: resData.PhotoProfile,
	}

	return entityData, nil
}

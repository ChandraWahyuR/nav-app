package usecase

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"html/template"
	"proyek1/config"
	"proyek1/internal/entity"
	"proyek1/internal/model"
	"proyek1/utils"
	jwt "proyek1/utils"
	"proyek1/utils/mailer"
	"strings"
	"sync"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

var (
	RoleUser  = "users"
	RoleAdmin = "admin"
)

type RepositoryUserInterface interface {
	Register(ctx context.Context, postData *entity.User) error
	Login(ctx context.Context, postData *entity.User) (*entity.User, error)
	GetUserID(ctx context.Context, id string) (entity.User, error)
	ForgotPassword(ctx context.Context, data *entity.Otp) error // Nanti
	OtpVerify(ctx context.Context, data *entity.Otp) (*entity.Otp, error)
	ResetPassword(ctx context.Context, data *entity.User) error

	EditDataUser(ctx context.Context, data *entity.User, id string) error
	//
	IsDataAvailable(ctx context.Context, email, username string) bool
	RoleChecker(ctx context.Context, id string) string
}

type UsecaseUser struct {
	jwt      jwt.JWTInterface
	userRepo RepositoryUserInterface
	log      *logrus.Logger
	cfg      *config.Config
	m        mailer.MailInterface
}

func NewUserUsecase(jwt jwt.JWTInterface, user RepositoryUserInterface, log *logrus.Logger, cfg *config.Config, m mailer.MailInterface) *UsecaseUser {
	return &UsecaseUser{
		jwt:      jwt,
		userRepo: user,
		log:      log,
		cfg:      cfg,
		m:        m,
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
	if !s.userRepo.IsDataAvailable(ctx, userData.Email, userData.Username) {
		return errors.New("email atau username sudah digunakan")
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

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		fmt.Println(resData)
		err := s.sendVerificationEmail(resData, "adsfd")
		if err != nil {
			fmt.Println("Gagal mengirim email verifikasi:", err)
		}
	}()

	return nil
}

func (s *UsecaseUser) Login(ctx context.Context, postData *model.Login) (*model.Login, error) {
	switch {
	case postData.Email == "":
		return &model.Login{}, errors.New("email tidak boleh kosong")
	case postData.Password == "":
		return &model.Login{}, errors.New("password tidak boleh kosong")
	}
	// Cek format Email
	if !utils.ValidasiEmail(postData.Email) {
		return &model.Login{}, errors.New("format email tidak benar")
	}

	postData.Email = strings.ToLower(postData.Email)
	userEntity := entity.User{
		Email:    postData.Email,
		Password: postData.Password,
	}

	res, err := s.userRepo.Login(ctx, &userEntity)
	if err != nil {
		return &model.Login{}, err
	}

	userRole := s.userRepo.RoleChecker(ctx, res.ID)
	if userRole == "" {
		return nil, errors.New("Error role tidak dapat ditentukan")
	}

	tokenData := model.User{
		ID:    res.ID,
		Email: res.Email,
		Role:  userRole,
	}
	token, err := s.jwt.GenerateToken(&tokenData)
	if err != nil {
		return &model.Login{}, err
	}

	UserLoginData := &model.Login{}
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

func (s *UsecaseUser) EditProfile(ctx context.Context, data *model.EditProfile, id string) error {
	var dataUser entity.User

	oldData, err := s.userRepo.GetUserID(ctx, id)
	if err != nil {
		return errors.New("akun tidak ditemukan")
	}

	if id == "" {
		return errors.New("akun tidak ditemukan")
	}
	// Validasi username
	if data.Username != "" {
		if !s.userRepo.IsDataAvailable(ctx, "", data.Username) {
			return errors.New("username sudah digunakan")
		}
		dataUser.Username = data.Username
	} else {
		dataUser.Username = oldData.Username
	}

	// Password Update
	if data.Password != "" {
		// Cek format password
		pass, err := utils.ValidatePassword(data.Password)
		if err != nil {
			return fmt.Errorf("password tidak valid: %w", err)
		}
		// Hash password
		hashedPassword, err := utils.HashPassword(pass)
		if err != nil {
			return err
		}
		dataUser.Password = hashedPassword
	} else {
		dataUser.Password = oldData.Password
	}

	if data.PhotoProfile != "" {
		dataUser.PhotoProfile = data.PhotoProfile
	} else {
		dataUser.PhotoProfile = oldData.PhotoProfile
	}

	err = s.userRepo.EditDataUser(ctx, &dataUser, id)
	if err != nil {
		return err
	}

	return nil
}

func (s *UsecaseUser) sendVerificationEmail(user entity.User, link string) error {
	data := map[string]interface{}{
		"Name": user.Username,
		"Link": link,
	}

	t, err := template.ParseFiles("./static/body.email.html")
	if err != nil {
		panic(err)
	}

	var body bytes.Buffer
	if err := t.Execute(&body, data); err != nil {
		return fmt.Errorf("failed to execute email template: %w", err)
	}

	return s.m.SendMail(
		user.Email,
		"Verifikasi Akun Anda",
		body.String(),
		data,
	)
}

func (s *UsecaseUser) RegisterForAdmin(ctx context.Context, req *model.User) error {
	switch {
	case req.Username == "":
		return errors.New("username tidak boleh kosong")
	case req.Email == "":
		return errors.New("email tidak boleh kosong")
	case req.Password == "":
		return errors.New("password tidak boleh kosong")
	case req.Password != req.ConfirmPassword:
		return errors.New("konfirmasi password tidak sama dengan password")
	}

	// Cek format password
	if !utils.ValidasiEmail(req.Email) {
		return errors.New("format email tidak benar")
	}

	// Cek ketersediaan email
	if !s.userRepo.IsDataAvailable(ctx, req.Email, req.Username) {
		return errors.New("email atau username sudah digunakan")
	}

	// Cek validasi password
	pass, err := utils.ValidatePassword(req.Password)
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
		Username:     req.Username,
		Email:        req.Email,
		Password:     hashedPassword,
		PhotoProfile: s.cfg.GeneralPhoto.DefaultPhoto,
		Role:         RoleAdmin,
		IsActive:     true, // sementara
	}

	err = s.userRepo.Register(ctx, &resData)
	if err != nil {
		return err
	}

	return nil
}

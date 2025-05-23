package usecase

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"html/template"
	"math/rand"
	"proyek1/config"
	"proyek1/internal/entity"
	"proyek1/internal/model"
	"proyek1/utils"
	jwt "proyek1/utils"
	"proyek1/utils/mailer"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

var (
	RoleUser   = "users"
	RoleAdmin  = "admin"
	RoleForgot = "forgot"
	SubjectOTP = "OTP Reset Password"
)

type RepositoryUserInterface interface {
	Register(ctx context.Context, postData *entity.User) error
	Login(ctx context.Context, postData *entity.User) (*entity.User, error)
	GetUserID(ctx context.Context, id string) (entity.User, error)
	ForgotPassword(ctx context.Context, data *entity.Otp) error // Nanti
	OtpVerify(ctx context.Context, email string, otp int) (*entity.Otp, error)
	SoftDeleteOtpByID(ctx context.Context, id string) error
	ResetPassword(ctx context.Context, data *entity.User) error

	EditDataUser(ctx context.Context, data *entity.User, id string) error
	//
	IsDataAvailable(ctx context.Context, email, username string) bool
	RoleChecker(ctx context.Context, id string) string
	ActivateAcount(ctx context.Context, email string) error
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
		IsActive:     false,
	}

	err = s.userRepo.Register(ctx, &resData)
	if err != nil {
		return err
	}

	token, err := s.jwt.GenerateToken(&model.User{ID: resData.ID, Email: resData.Email, Role: resData.Role})
	if err != nil {
		return err
	}
	link := fmt.Sprintf(`%s/active?token=%s`, s.cfg.URL_Server, token)

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		fmt.Println(resData)
		err := s.sendVerificationEmail(resData, link)
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
		return errors.New("password is invalid")
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
		IsActive:     true,
	}

	err = s.userRepo.Register(ctx, &resData)
	if err != nil {
		return err
	}

	return nil
}

func (s *UsecaseUser) ForgotPassword(ctx context.Context, req *model.Otp) error {
	if req.Email == "" || !utils.ValidasiEmail(req.Email) {
		return errors.New("format email tidak benar")
	}

	var entityModel entity.Otp
	entityModel = entity.Otp{
		ID:    uuid.New().String(),
		Email: req.Email,
	}

	entityModel.OtpNumber, _ = strconv.Atoi(fmt.Sprintf("%05d", rand.Intn(100000)))
	entityModel.ValidUntil = time.Now().Add(5 * time.Minute)
	data := map[string]interface{}{
		"OTP":         entityModel.OtpNumber,
		"Valid_Until": entityModel.ValidUntil.Format("02 Jan 2006, 15:04:05 MST"),
	}

	template, err := template.ParseFiles("./static/body.otp.html")
	if err != nil {
		return fmt.Errorf("gagal load template email: %w", err)
	}
	var body bytes.Buffer
	if err = template.Execute(&body, data); err != nil {
		return fmt.Errorf("failed to execute email template: %w", err)
	}

	err = s.userRepo.ForgotPassword(ctx, &entityModel)
	if err != nil {
		return err
	}

	go func() {
		err := s.m.SendMail(entityModel.Email, SubjectOTP, body.String(), data)
		if err != nil {
			s.log.WithContext(ctx).Error("gagal kirim email OTP:", err)
		}
	}()

	return nil
}

func (s *UsecaseUser) OtpVerify(ctx context.Context, req *model.Otp) (*model.Otp, error) {
	if req.Email == "" || !utils.ValidasiEmail(req.Email) {
		return nil, errors.New("format email tidak benar")
	}

	data, err := s.userRepo.OtpVerify(ctx, req.Email, req.OtpNumber)
	if err != nil {
		return nil, err
	}

	if data.ValidUntil.Before(time.Now()) {
		return nil, errors.New("otp telah expired")
	}

	if req.OtpNumber != data.OtpNumber {
		return nil, errors.New("otp salah")
	}

	err = s.userRepo.SoftDeleteOtpByID(ctx, data.ID)
	if err != nil {
		return nil, fmt.Errorf("gagal menghapus otp: %v", err)
	}

	var modelUser model.User
	modelUser = model.User{
		ID:    req.ID,
		Email: req.Email,
		Role:  RoleForgot,
	}
	token, err := s.jwt.GenerateToken(&modelUser)
	if err != nil {
		return nil, err
	}
	tokenData := &model.Otp{}
	tokenData.Token = token
	return tokenData, nil
}

func (s *UsecaseUser) ResetPassword(ctx context.Context, req *model.User) error {
	if req.Password != req.ConfirmPassword {
		return errors.New("konfirmasi password tidak sama dengan password")
	}

	pass, err := utils.ValidatePassword(req.Password)
	if err != nil {
		return errors.New("password is invalid")
	}

	hashedPassword, err := utils.HashPassword(pass)
	if err != nil {
		return err
	}

	var entityModel entity.User
	entityModel = entity.User{
		ID:       req.ID,
		Email:    req.Email,
		Password: hashedPassword,
	}

	err = s.userRepo.ResetPassword(ctx, &entityModel)
	if err != nil {
		return err
	}

	return nil
}

func (s *UsecaseUser) ActivateAcount(ctx context.Context, email string) error {
	if email == "" {
		return errors.New("email tidak boleh kosong")
	}

	err := s.userRepo.ActivateAcount(ctx, email)
	if err != nil {
		return err
	}

	return nil
}

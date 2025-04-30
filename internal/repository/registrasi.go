package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"proyek1/internal/entity"
	"proyek1/utils"
	"time"

	"github.com/sirupsen/logrus"
)

type RegistrasiRepo struct {
	db  *sql.DB
	Log *logrus.Logger
}

func NewUserRepository(db *sql.DB, log *logrus.Logger) *RegistrasiRepo {
	return &RegistrasiRepo{
		db:  db,
		Log: log,
	}
}

// Tambahan context
func (r *RegistrasiRepo) Register(ctx context.Context, postData *entity.User) error {
	query := `INSERT INTO "users"(id, username, email, password, photo_profile, role, is_active)  VALUES($1, $2, $3, $4, $5, $6, $7)`
	_, err := r.db.ExecContext(ctx, query, postData.ID, postData.Email, postData.Username, postData.Password, postData.PhotoProfile, postData.Role, postData.IsActive)
	if err != nil {
		return fmt.Errorf(`gagal memasukkan data ke database : %s`, err)
	}

	return nil
}

func (r *RegistrasiRepo) Login(ctx context.Context, postData *entity.User) (*entity.User, error) {
	var userData entity.User
	query := `SELECT id, email, password FROM "users" WHERE email = $1`
	err := r.db.QueryRowContext(ctx, query, postData.Email).Scan(&userData.Email, &userData.Password)
	if err != nil {
		fmt.Errorf(`Gagal menginputkan data :%s`, err)
		if err == sql.ErrNoRows {
			fmt.Errorf(`Gagal menginputkan data :%s`, err)
			return nil, errors.New("email tidak ada")
		}
		return nil, err
	}

	if !utils.ValidasiEmail(postData.Email) {
		return nil, errors.New("error email tidak ada")
	}

	return &userData, nil
}

func (r *RegistrasiRepo) GetUserID(ctx context.Context, id string) (entity.User, error) {
	var data entity.User
	query := `SELECT * FROM "users" WHERE id = $1`
	err := r.db.QueryRowContext(ctx, query, id).Scan(&data)
	if err != nil {
		return entity.User{}, err
	}

	return data, nil
}

func (r *RegistrasiRepo) ForgotPassword(ctx context.Context, data *entity.Otp) error {
	var dummy string
	cekQuery := `SELECT email FROM "users" WHERE email = $1`
	err := r.db.QueryRowContext(ctx, cekQuery, data.Email).Scan(&dummy)
	if err == sql.ErrNoRows {
		r.Log.WithContext(ctx).Warn("email tidak ditemukan")
		return err
	}
	if err != nil {
		r.Log.WithContext(ctx).Error("gagal mengecek email:", err)
		return err
	}

	softDeleteQuery := `UPDATE otp SET deleted_at = NOW() WHERE email = $1 AND deleted_at IS NULL`
	_, err = r.db.ExecContext(ctx, softDeleteQuery, data.Email)
	if err != nil {
		r.Log.WithContext(ctx).Error("gagal soft delete OTP:", err)
		return err
	}

	insertQuery := `INSERT INTO "otp"(id, email, otp_number, valid_until) VALUES($1, $2, $3, $4)`
	_, err = r.db.ExecContext(ctx, insertQuery, data.ID, data.Email, data.OtpNumber, data.ValidUntil)
	if err != nil {
		r.Log.WithContext(ctx).Error("gagal menyimpan OTP:", err)
		return err
	}

	return nil
}

// return entity, data mau dipakai di usecase, kalau hanya cek atau create saja return error saja
func (r *RegistrasiRepo) OtpVerify(ctx context.Context, data *entity.Otp) (*entity.Otp, error) {
	var dataEntity entity.Otp
	query := `SELECT * FROM "otp" WHERE Email = $1 AND otp_number = $2 AND deleted_at IS Null`
	err := r.db.QueryRowContext(ctx, query, data.Email, data.OtpNumber).Scan(&dataEntity)
	if err != nil {
		return nil, err
	}

	queryUpdate := `UPDATE otp SET deleted_at = $1 WHERE id = $2`
	_, err = r.db.ExecContext(ctx, queryUpdate, time.Now(), dataEntity.ID)
	if err != nil {
		return nil, err
	}

	return &dataEntity, nil
}

func (r *RegistrasiRepo) ResetPassword(ctx context.Context, data *entity.User) error {
	query := `UPDATE users SET password = $1 WHERE email = $2 AND deleted_at IS NULL`
	_, err := r.db.ExecContext(ctx, query, data.Password, data.Email)
	if err != nil {
		return err
	}

	return nil
}

// Validasi
func (r *RegistrasiRepo) IsEmailAvailable(ctx context.Context, email string) bool {
	var data string // cuman cek aja pakai string langsung, kalau mau ambil data ambil dari entity
	query := `SELECT username FROM "users" WHERE email = $1`
	err := r.db.QueryRowContext(ctx, query, email).Scan(&data)
	if err != nil {
		return false
	}
	return true
}

func (r *RegistrasiRepo) IsUsernameAvailable(ctx context.Context, username string) bool {
	var data string
	query := `SELECT username FROM "users" WHERE username = $1`
	err := r.db.QueryRowContext(ctx, query, username).Scan(&data)
	if err != nil {
		return false
	}
	return true
}

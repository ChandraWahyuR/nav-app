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
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()
	query := `INSERT INTO "users"(id, username, email, password, photo_profile, role, is_active)  VALUES($1, $2, $3, $4, $5, $6, $7)`
	_, err := r.db.ExecContext(ctx, query,
		postData.ID,
		postData.Username,
		postData.Email,
		postData.Password,
		postData.PhotoProfile,
		postData.Role,
		postData.IsActive,
	)

	if err != nil {
		return utils.ParsePQError(err)
	}

	return nil
}

func (r *RegistrasiRepo) Login(ctx context.Context, postData *entity.User) (*entity.User, error) {
	var userData entity.User
	query := `SELECT id, email, password FROM "users" WHERE email = $1`
	err := r.db.QueryRowContext(ctx, query, postData.Email).Scan(
		&userData.ID, &userData.Email, &userData.Password)

	fmt.Println("Email dicari:", postData.Email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("email tidak ditemukan")
		}
		return nil, fmt.Errorf("gagal mengambil data login: %w", err)
	}

	if !utils.VerifyHashedPassword(postData.Password, userData.Password) {
		return nil, errors.New("password salah")
	}
	return &userData, nil
}

func (r *RegistrasiRepo) GetUserID(ctx context.Context, id string) (entity.User, error) {
	var data entity.User
	query := `SELECT id, email, username, password, role, photo_profile, is_active FROM "users" WHERE id = $1`
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&data.ID,
		&data.Email,
		&data.Username,
		&data.Password,
		&data.Role,
		&data.PhotoProfile,
		&data.IsActive,
	)

	if err != nil {
		return entity.User{}, err
	}

	return data, nil
}

func (r *RegistrasiRepo) ForgotPassword(ctx context.Context, data *entity.Otp) error {
	var dummy string
	cekQuery := `SELECT email FROM "users" WHERE email = $1 AND deleted_at IS NULL`
	err := r.db.QueryRowContext(ctx, cekQuery, data.Email).Scan(&dummy)
	if err == sql.ErrNoRows {
		r.Log.WithContext(ctx).Warn("email tidak ditemukan")
		return err
	}
	if err != nil {
		r.Log.WithContext(ctx).Error("gagal mengecek email:", err)
		return err
	}

	softDeleteQuery := `UPDATE otp SET deleted_at = NOW(), status = false WHERE email = $1 AND deleted_at IS NULL`
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
func (r *RegistrasiRepo) OtpVerify(ctx context.Context, email string, otp int) (*entity.Otp, error) {
	var dataEntity entity.Otp
	query := `SELECT id, email, otp_number, valid_until FROM "otp" WHERE Email = $1 AND otp_number = $2 AND deleted_at IS Null`
	err := r.db.QueryRowContext(ctx, query, email, otp).Scan(
		&dataEntity.ID,
		&dataEntity.Email,
		&dataEntity.OtpNumber,
		&dataEntity.ValidUntil,
	)
	if err != nil {
		return nil, utils.ParsePQError(err)
	}

	softDeleteQuery := `UPDATE otp SET deleted_at = NOW(), status = true WHERE email = $1 AND deleted_at IS NULL`
	_, err = r.db.ExecContext(ctx, softDeleteQuery, email)
	if err != nil {
		r.Log.WithContext(ctx).Error("gagal soft delete OTP:", err)
		return nil, utils.ParsePQError(err)
	}

	return &dataEntity, nil
}

func (r *RegistrasiRepo) SoftDeleteOtpByID(ctx context.Context, id string) error {
	query := `UPDATE otp SET deleted_at = $1 WHERE id = $2`
	_, err := r.db.ExecContext(ctx, query, time.Now(), id)
	return err
}

func (r *RegistrasiRepo) ResetPassword(ctx context.Context, data *entity.User) error {
	query := `UPDATE users SET password = $1 WHERE email = $2 AND deleted_at IS NULL`
	_, err := r.db.ExecContext(ctx, query, data.Password, data.Email)
	if err != nil {
		return utils.ParsePQError(err)
	}

	return nil
}

func (r *RegistrasiRepo) EditDataUser(ctx context.Context, data *entity.User, id string) error {
	query := `UPDATE users SET username = $1, password = $2, photo_profile = $3 WHERE id = $4 AND deleted_at IS NULL`
	result, err := r.db.ExecContext(ctx, query, data.Username, data.Password, data.PhotoProfile, id)
	if err != nil {
		return utils.ParsePQError(err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("user tidak ditemukan atau sudah dihapus")
	}

	return nil
}

// Validasi
func (r *RegistrasiRepo) IsDataAvailable(ctx context.Context, email, username string) bool {
	var data string // cuman cek aja pakai string langsung, kalau mau ambil data ambil dari entity
	query := `SELECT email, username FROM users WHERE email = $1 OR username = $2`
	err := r.db.QueryRowContext(ctx, query, email, username).Scan(&data)
	if err != nil {
		return true
	}
	return false
}

func (r *RegistrasiRepo) RoleChecker(ctx context.Context, id string) string {
	var data string
	query := `SELECT role FROM users WHERE id = $1`
	err := r.db.QueryRowContext(ctx, query, id).Scan(&data)
	if err != nil {
		return ""
	}
	return data
}

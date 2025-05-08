package utils

import (
	"errors"
	"fmt"
	"regexp"

	"golang.org/x/crypto/bcrypt"
)

// Validasi penulisan email
func ValidasiEmail(email string) bool {
	regexEmail := regexp.MustCompile(`^[a-zA-Z0-9._-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,4}$`)
	return regexEmail.MatchString(email)
}

// Password di hash
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// Verifikasi password dengan password yang dihash
func VerifyHashedPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// Ketentuan Password
func ValidatePassword(password string) (string, error) {
	if len(password) < 8 || len(password) > 12 {
		return "", errors.New("panjang password harus antara 8 sampai 12 karakter")
	}

	var (
		hasUpper  = regexp.MustCompile(`[A-Z]`).MatchString
		hasLower  = regexp.MustCompile(`[a-z]`).MatchString
		hasNumber = regexp.MustCompile(`[0-9]`).MatchString
		hasSymbol = regexp.MustCompile(`[@$!%*?&#]`).MatchString
	)

	switch {
	case !hasUpper(password):
		return "", errors.New("password harus mengandung huruf besar")
	case !hasLower(password):
		return "", errors.New("password harus mengandung huruf kecil")
	case !hasNumber(password):
		return "", errors.New("password harus mengandung angka")
	case !hasSymbol(password):
		return "", errors.New("password harus mengandung simbol spesial (@$!%*?&#)")
	}
	return password, nil
}

func FormatJam(jam string) string {
	if len(jam) < 4 {
		return fmt.Sprint("kurang")
	}

	angka1 := jam[:2]
	return fmt.Sprintf(`%s:%s`, angka1, jam[2:])
}

func TotalPageForPagination(totalDataDB, limit int) int {
	return (totalDataDB + limit - 1) / limit
}

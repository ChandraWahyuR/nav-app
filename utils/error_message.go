package utils

import "errors"

const (
	Unauthorized   = "Unauthorized"
	InternalServer = "Internal Server Error"
	BadInput       = "Format data not valid"
)

var (
	ErrGetData            = errors.New("gagal saat mengambil data")
	ErrEmailTaken         = errors.New("email sudah digunakan")
	ErrUsernameTaken      = errors.New("username sudah digunakan")
	ErrPlaceIDUniqueTaken = errors.New("id tempat sudah digunakan, id harus unique")
)

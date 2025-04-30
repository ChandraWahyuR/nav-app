package utils

import "errors"

const (
	Unauthorized   = "Unauthorized"
	InternalServer = "Internal Server Error"
	BadInput       = "Format data not valid"
)

var ErrGetData = errors.New("gagal saat mengambil data")

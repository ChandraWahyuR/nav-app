package model

import "time"

type Otp struct {
	ID         string    `json:"id"`
	Email      string    `json:"email"`
	OtpNumber  int       `json:"otp_number"`
	ValidUntil time.Time `json:"valid_until"`
}

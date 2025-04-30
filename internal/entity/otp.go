package entity

import "time"

type Otp struct {
	ID         string
	Email      string
	OtpNumber  int
	ValidUntil time.Time
}

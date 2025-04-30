package model

type User struct {
	ID              string `json:"id"`
	Username        string `json:"username"`
	Email           string `json:"email"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirm_password"`
	PhotoProfile    string `json:"photo_profile"`
	Role            string `json:"role"`
	Token           string `json:"token"`
	IsActive        bool   `json:"is_active"`
}

type Register struct {
	Username        string `json:"username"`
	Email           string `json:"email"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirm_password"`
}

type Profile struct {
	ID           string `json:"id"`
	Username     string `json:"username"`
	Email        string `json:"email"`
	PhotoProfile string `json:"photo_profile"`
}

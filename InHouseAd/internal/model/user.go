package model

type User struct {
	UserID      uint   `json:"id"`
	Username    string `json:"username" validate:"required"`
	Password    string `json:"password" validate:"required"`
	EMail       string `json:"email" validate:"required"`
	PhoneNumber string `json:"phone_number"`
}

type Register struct {
	Username    string `json:"username"`
	Password    string `json:"password"`
	EMail       string `json:"email"`
	PhoneNumber string `json:"phone-number"`
}

type Login struct {
	Username string
	Password string
}

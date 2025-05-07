package models

type User struct {
	ID uint64 `json:"id"`
	Email string `json:"email"`
	Password PasswordInfo
}

type PasswordInfo struct {
	Password string
}

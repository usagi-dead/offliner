package models

import "time"

type User struct {
	UserId         int64      `json:"user_id"`
	HashedPassword *string    `json:"hashed_password"`
	Role           string     `json:"role" default:"user"`
	Surname        *string    `json:"surname,omitempty"`
	Name           *string    `json:"name,omitempty"`
	Patronymic     *string    `json:"patronymic,omitempty"`
	DateOfBirth    *time.Time `json:"date_of_birth,omitempty"`
	PhoneNumber    *string    `json:"phone_number,omitempty"`
	Email          string     `json:"email"`
	AvatarUrl      *string    `json:"avatar_url,omitempty"`
	VerifiedEmail  bool       `json:"verified_email"`
	Gender         *string    `json:"gender,omitempty"`
}

package models

import "time"

type User struct {
	UserId         int64     `json:"user_id"`
	HashedPassword string    `json:"hashed_password"`
	Role           string    `json:"role"`
	Surname        string    `json:"surname"`
	Name           string    `json:"name"`
	Patronymic     string    `json:"patronymic"`
	DateOfBirth    time.Time `json:"date_of_birth"`
	PhoneNumber    string    `json:"phone_number"`
	Email          string    `json:"email"`
	Gender         string    `json:"gender"`
}

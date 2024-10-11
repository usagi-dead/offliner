package models

type User struct {
	UserId         int64
	HashedPassword string
	Role           string
	Surname        string
	Name           string
	Patronymic     string
	DateOfBirth    string
	PhoneNumber    string
	Email          string
	Gender         string
}

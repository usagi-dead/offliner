package user

import (
	"mime/multipart"
	"net/http"
	"time"
)

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

type UserHandler interface {
	SignUp(w http.ResponseWriter, r *http.Request)
	SignIn(w http.ResponseWriter, r *http.Request)
	Oauth(w http.ResponseWriter, r *http.Request)
	OauthCallback(w http.ResponseWriter, r *http.Request)
	SendConfirmedEmailCode(w http.ResponseWriter, r *http.Request)
	EmailConfirmed(w http.ResponseWriter, r *http.Request)
	RefreshToken(w http.ResponseWriter, r *http.Request)
	UpdateUserHandler(w http.ResponseWriter, r *http.Request)
	GetUserHandler(w http.ResponseWriter, r *http.Request)
	LogoutHandler(w http.ResponseWriter, r *http.Request)
	DeleteUserHandler(w http.ResponseWriter, r *http.Request)
}

type UserService interface {
	SignUp(email string, password string) error
	SignIn(email string, password string) (string, string, error)
	GetAuthURL(provider string) (string, error)
	Callback(provider, state, code string) (bool, string, string, error)
	SendEmailForConfirmed(email string) error
	EmailConfirmed(email string, code string) error
	RefreshToken(r *http.Request) (string, error)
	UpdateUser(avatar *multipart.File, user *User, resetAvatar bool) error
	GetUser(userId int64) (*User, error)
	DeleteUser(userId int64) error
}

type UserData interface {
	CreateUser(email string, hashedPassword string) (int64, error)
	CreateOauthUser(user *User) (int64, error)
	GetUserByEmail(Email string) (*User, error)
	GetUserById(userId int64) (*User, error)
	SaveStateCode(state string) error
	VerifyStateCode(state string) (bool, error)
	UpdateUser(user *User) error
	UploadAvatar(avatarSmall []byte, avatarLarge []byte, userId int64) (*string, error)
	ConfirmEmail(email string) error
	IsEmailConfirmed(email string) (bool, error)
	SaveEmailConfirmedCode(email string, code string) error
	GetEmailConfirmedCode(email string) (string, error)
	DeleteUser(userId int64) error
	DeleteAvatar(userId int64) error
}

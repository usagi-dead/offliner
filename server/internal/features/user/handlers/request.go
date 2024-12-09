package handlers

import "time"

type UserSignUpRequest struct {
	Email    string `json:"email" validate:"required,email" example:"jon.doe@gmail.com"`
	Password string `json:"password" validate:"required" example:"SuperPassword123"`
}

type UserSingInRequest struct {
	Email    string `json:"email" validate:"required,email" example:"jon.doe@gmail.com"`
	Password string `json:"password" validate:"required" example:"SuperPassword123"`
}

type SendConfirmedEmailCodeRequest struct {
	Email string `json:"email" validate:"required,email" example:"jon.doe@gmail.com"`
}

type EmailConfirmedRequest struct {
	Code  string `json:"code" validate:"required" example:"54JK64"`
	Email string `json:"email" validate:"required,email" example:"jon.doe@gmail.com"`
}

type UpdateUserRequest struct {
	Surname     *string    `json:"surname,omitempty" validate:"omitempty,min=1,max=50" example:"John"`
	Name        *string    `json:"name,omitempty" validate:"omitempty,min=1,max=50" example:"Doe"`
	Patronymic  *string    `json:"patronymic,omitempty" validate:"omitempty,min=1,max=50" example:"Smith"`
	DateOfBirth *time.Time `json:"date_of_birth,omitempty" validate:"omitempty,lte" example:"1985-04-12T00:00:00Z"`
	PhoneNumber *string    `json:"phone_number,omitempty" validate:"omitempty,e164" example:"+1234567890"`
	Gender      *string    `json:"gender,omitempty" validate:"omitempty,oneof=male female" example:"male"`
	ResetAvatar bool       `json:"reset_avatar" default:"false" example:"false"`
}

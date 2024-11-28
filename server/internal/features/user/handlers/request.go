package handlers

type UserSignUpRequest struct {
	Email    string `json:"email" validate:"required,email" example:"jon.doe@gmail.com"`
	Password string `json:"password" validate:"required" example:"SuperPassword123"`
}

type UserSingInRequest struct {
	Email    string `json:"email" validate:"required,email" example:"jon.doe@gmail.com"`
	Password string `json:"password" validate:"required" example:"SuperPassword123"`
}

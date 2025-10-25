package requests

type SignupRequest struct {
	Email    string `json:"email" binding:"required,email" validate:"required"`
	Password string `json:"password" binding:"required,min=6" validate:"required"`
}

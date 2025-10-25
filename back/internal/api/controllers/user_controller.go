package controllers

import (
	"chat/internal/api/dto/requests"
	"chat/internal/api/services"
	validatorHandler "chat/internal/api/validator"
	"context"
	"encoding/json"
	"net/http"

	"github.com/go-playground/validator/v10"
)

type AuthCtrl struct {
	userSrv   services.UserSrv
	validator *validator.Validate
}

func NewAuthCtrl() *AuthCtrl {
	return &AuthCtrl{
		userSrv:   *services.NewUserSrv(),
		validator: validatorHandler.GetValidator(),
	}
}

func (authCtrl *AuthCtrl) HandleSignupUser(w http.ResponseWriter, r *http.Request) {

	var req requests.SignupRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	err = authCtrl.validator.Struct(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ctx := context.Background()
	createdUser, err := authCtrl.userSrv.CreateUserSrv(ctx, req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(createdUser)
}

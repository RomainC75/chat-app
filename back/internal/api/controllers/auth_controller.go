package controllers

import (
	"chat/internal/api/dto/requests"
	repo_db "chat/internal/api/repos/db"
	validatorHandler "chat/internal/api/validator"
	user_management_app "chat/internal/modules/user-management/application"
	user_management_infra "chat/internal/modules/user-management/infra"
	"chat/utils"
	"context"
	"encoding/json"
	"net/http"

	"github.com/go-playground/validator/v10"
)

type VerifyResposne struct {
	Id    int32  `json:"id"`
	Email string `json:"email"`
}

type AuthCtrl struct {
	userSrv   *user_management_app.UserSrv
	validator *validator.Validate
}

func NewAuthCtrl() *AuthCtrl {
	return &AuthCtrl{
		userSrv: user_management_app.NewUserSrv(
			repo_db.NewUserRepo(),
			user_management_infra.NewInMemoryUUIDGenerator(),
			user_management_infra.NewInMemoryClock(),
			user_management_infra.NewInMemoryBcrypt(),
			user_management_infra.NewInMemoryJWT(),
		),
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

func (authCtrl *AuthCtrl) HandleLoginUser(w http.ResponseWriter, r *http.Request) {

	var req requests.LoginRequest
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
	token, err := authCtrl.userSrv.LogUserSrv(ctx, req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	utils.JsonResponse(w, token)
}

func (authCtrl *AuthCtrl) HandleVerify(w http.ResponseWriter, r *http.Request) {
	email := r.Context().Value("user_email")
	Id := r.Context().Value("user_id")

	id, _ := Id.(float64)
	id32 := int32(id)

	res := VerifyResposne{
		Id:    id32,
		Email: email.(string),
	}
	utils.JsonResponse(w, res)
}

package routes

import (
	"chat/internal/api/controllers"
	"chat/internal/api/middlewares"
	"net/http"

	"github.com/gorilla/mux"
)

func UserRoutes(mux *mux.Router) {
	authCtrl := controllers.NewAuthCtrl()

	userRouter := mux.PathPrefix("/user").Subrouter()

	userRouter.Use(middlewares.AuthMid)

	userRouter.Handle("/user/signup", http.HandlerFunc(authCtrl.HandleSignupUser)).Methods("POST")
	userRouter.Handle("/user/login", http.HandlerFunc(authCtrl.HandleLoginUser)).Methods("POST")
	userRouter.Handle("/user/verify", http.HandlerFunc(authCtrl.HandleVerify)).Methods("GET")

}

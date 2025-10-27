package routes

import (
	"chat/internal/api/controllers"
	"chat/internal/api/middlewares"
	"net/http"

	"github.com/gorilla/mux"
)

func AuthRoutes(mux *mux.Router) {
	authCtrl := controllers.NewAuthCtrl()

	userRouter := mux.PathPrefix("/auth").Subrouter()

	userRouter.Handle("/signup", http.HandlerFunc(authCtrl.HandleSignupUser)).Methods("POST")
	userRouter.Handle("/login", http.HandlerFunc(authCtrl.HandleLoginUser)).Methods("POST")
	userRouter.Handle("/verify", middlewares.AuthMid(http.HandlerFunc(authCtrl.HandleVerify))).Methods("GET")

}

package routes

import (
	"chat/internal/api/controllers"
	"chat/internal/api/middlewares"
	"net/http"

	"github.com/gorilla/mux"
)

func UserRoutes(mux *mux.Router) {
	authCtrl := controllers.NewAuthCtrl()

	mux.Use(middlewares.AuthMid)
	mux.Handle("/user/signup", http.HandlerFunc(authCtrl.HandleSignupUser)).Methods("POST")
	mux.Handle("/user/login", http.HandlerFunc(authCtrl.HandleLoginUser)).Methods("POST")
	mux.Handle("/user/verify", http.HandlerFunc(authCtrl.HandleVerify)).Methods("GET")

}

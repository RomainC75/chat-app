package routes

import (
	"chat/internal/api/controllers"
	"chat/internal/api/middlewares"
	"net/http"

	"github.com/gorilla/mux"
)

func UserRoutes(mux *mux.Router) {
	authCtrl := controllers.NewAuthCtrl()
	mux.Handle("/user/signup", middlewares.CORSMiddleware(http.HandlerFunc(authCtrl.HandleSignupUser))).Methods("POST")

}

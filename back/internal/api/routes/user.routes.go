package routes

import (
	"chat/internal/api/middlewares"
	"net/http"

	"github.com/gorilla/mux"
)

func UserRoutes(mux *mux.Router) {

	mux.Handle("/user/signup", middlewares.CORSMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello"))
	})))

}

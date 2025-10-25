package routes

import (
	"chat/internal/api/middlewares"
	"net/http"

	"github.com/gorilla/mux"
)

func HealthRoutes(mux *mux.Router) {

	mux.Handle("/health/test", middlewares.CORSMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("is healthy"))
	})))

}

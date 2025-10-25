package routes

import (
	"chat/internal/api/middlewares"

	"github.com/gorilla/mux"
)

func ConnectRoutes() *mux.Router {
	r := mux.NewRouter()

	r.Use(middlewares.CORSMiddleware)
	api := r.PathPrefix("/api").Subrouter()

	HealthRoutes(api)
	UserRoutes(api)

	return api
}

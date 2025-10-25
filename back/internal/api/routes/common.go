package routes

import (
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func ConnectRoutes() http.Handler {
	r := mux.NewRouter()

	corsHandler := handlers.CORS(
		handlers.AllowedOrigins([]string{"*"}),
		handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}),
		handlers.AllowedHeaders([]string{"Content-Type", "Authorization", "X-Requested-With"}),
		handlers.AllowCredentials(),
	)

	r.Use(corsHandler)

	api := r.PathPrefix("/api").Subrouter()

	ChatRoutes(api)
	HealthRoutes(api)
	UserRoutes(api)

	return corsHandler(r)
}

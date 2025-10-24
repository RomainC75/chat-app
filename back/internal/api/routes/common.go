package routes

import "net/http"

func ConnectRoutes() *http.ServeMux {
	mux := http.NewServeMux()
	mux.Handle("/health", HealthRoutes())
	// mux.Handle("/user", UserRoutes())

	return mux
}

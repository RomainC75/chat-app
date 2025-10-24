package routes

import "net/http"

func HealthRoutes() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("is healthy"))
	})
	return mux
}

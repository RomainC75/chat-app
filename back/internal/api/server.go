package api

import (
	"chat/internal/api/routes"
	"log/slog"
	"net/http"
)

func Serve() {
	srv := http.Server{
		Addr:    ":3000",
		Handler: routes.ConnectRoutes(),
	}

	err := srv.ListenAndServe()
	if err != nil {
		slog.Error("error running server")
	}
}

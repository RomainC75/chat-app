package api

import (
	"log/slog"
	"net/http"
)

func Serve() {
	srv := http.Server{}

	err := srv.ListenAndServe()
	if err != nil {
		slog.Error("error running server")
	}
}

package middlewares

import (
	"chat/utils/encrypt"
	"context"
	"net/http"
	"strings"
)

func AuthMid(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var token string

		auth_header := r.Header.Get("Authorization")

		w.Header()
		if !strings.HasPrefix(auth_header, "Bearer") {
			http.Error(w, "token missing", http.StatusBadRequest)
			return
		}
		token = strings.Split(auth_header, " ")[1]

		claim, err := encrypt.GetClaimsFromToken(token)
		if err != nil {
			http.Error(w, "unauthorized", http.StatusBadRequest)
			return
		}

		bgCtx := context.Background()
		wvCtx := context.WithValue(bgCtx, "user_email", claim["Email"])
		wvCtx = context.WithValue(wvCtx, "user_id", claim["ID"])

		next.ServeHTTP(w, r.WithContext(wvCtx))
	})
}

package middleware

import (
	"fmt"
	"net/http"

	"github.com/JakubDaleki/transfer-app/webapp/api/resource/auth"
	"github.com/golang-jwt/jwt/v4"
)

func AuthMiddleware(next func(http.ResponseWriter, *http.Request, string)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		bearer := r.Header.Get("bearer")
		if bearer == "" {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(fmt.Sprintf("{\"error\": \"%s\"}", "Bearer Token Not Provided")))
			return
		}

		claims := new(auth.Claims)
		_, err := jwt.ParseWithClaims(bearer, claims, func(token *jwt.Token) (interface{}, error) { return auth.SampleSecretKey, nil })
		if err == nil {
			next(w, r, claims.Username)
			return
		}

		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(fmt.Sprintf("{\"error\": \"%v\"}", err)))
	}
}

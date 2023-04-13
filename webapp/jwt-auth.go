package main

import (
	"encoding/json"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"log"
	"net/http"
	"time"
)

var sampleSecretKey = []byte("SecretYouShouldHide")

// Credentials Create a struct to read the username and password from the request body
type Credentials struct {
	Password string `json:"password"`
	Username string `json:"username"`
}

// Claims Create a struct that will be encoded to a JWT.
// We add jwt.RegisteredClaims as an embedded type, to provide fields like expiry time
type Claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

func authHandler(w http.ResponseWriter, r *http.Request) {
	var creds Credentials
	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	expectedPassword := users[creds.Username]

	if expectedPassword != creds.Password {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	expirationTime := time.Now().Add(5 * time.Hour)

	claims := Claims{
		Username: creds.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(sampleSecretKey)

	if err != nil {
		w.WriteHeader(http.StatusFailedDependency)
		log.Fatal(err)
	}

	http.SetCookie(w, &http.Cookie{
		Name:    "token",
		Value:   tokenString,
		Expires: expirationTime,
	})

	w.Write([]byte(fmt.Sprintf("{\"access_token\": \"%s\"}", tokenString)))
}

func authMiddleware(next func(http.ResponseWriter, *http.Request, string)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		bearer := r.Header.Get("bearer")
		if bearer == "" {
			w.Write([]byte(fmt.Sprintf("{\"error\": \"%s\"}", "Bearer Token Not Provided")))
			return
		}

		claims := new(Claims)
		_, err := jwt.ParseWithClaims(bearer, *claims, func(token *jwt.Token) (interface{}, error) { return sampleSecretKey, nil })
		if err != nil {
			next(w, r, claims.Username)
			return
		}

		w.Write([]byte(fmt.Sprintf("{\"access_token\": \"%s\"}", bearer)))
	}
}

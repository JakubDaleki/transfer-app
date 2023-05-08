package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/JakubDaleki/transfer-app/webapp/api/resource/auth"
	"github.com/JakubDaleki/transfer-app/webapp/utils/db"
	"github.com/golang-jwt/jwt/v4"
)

func AuthHandler(w http.ResponseWriter, r *http.Request, connector *db.Connector) {
	var creds auth.Credentials
	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	expectedPassword := connector.GetPassword(creds.Username)

	if expectedPassword != creds.Password {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(fmt.Sprintf("{\"error\": \"Wrong password or user does not exist.\"}")))
		return
	}

	expirationTime := time.Now().Add(5 * time.Hour)

	claims := auth.Claims{
		Username: creds.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(auth.SampleSecretKey)

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

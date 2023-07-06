package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
)

func (app *App) CreateJWT() (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["exp"] = time.Now().Add(time.Hour).Unix()
	tokenString, err := token.SignedString([]byte(app.Config.SecretKey))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func (app *App) ValidateJWT(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header["Token"] != nil {
			token, err := jwt.Parse(r.Header["Token"][0], func(token *jwt.Token) (interface{}, error) {
				_, ok := token.Method.(*jwt.SigningMethodHMAC)
				if !ok {
					w.WriteHeader(http.StatusUnauthorized)
					w.Write([]byte("not authorized"))
				}
				return []byte(app.Config.SecretKey), nil
			})

			if err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte("not authorized " + err.Error()))
			}

			if token.Valid {
				next.ServeHTTP(w, r)
			}
		} else {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("not authorized"))
		}
	})
}

func (app *App) GetJWT(w http.ResponseWriter, r *http.Request) {
	if r.Header["Api-Key"] != nil {
		if r.Header["Api-Key"][0] != app.Config.ApiKey {
			return
		} else {
			token, err := app.CreateJWT()
			if err != nil {
				return
			}
			fmt.Fprint(w, token)
		}
	}
}

package gameapi

import (
	"github.com/google/uuid"
	"net/http"
	"time"
)

func jsonContentTypeMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

const playerIdCookieName = "planning-player-id"

func playerIdCookieMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		_, err := request.Cookie(playerIdCookieName)
		if err != nil {
			expires := time.Now().Add(time.Hour * 24 * 30)
			cookie := http.Cookie{
				Name:    playerIdCookieName,
				Value:   uuid.New().String(),
				Expires: expires,
				MaxAge:  0,
				Path:    "/",
			}
			http.SetCookie(response, &cookie)
		}
		next.ServeHTTP(response, request)
	})
}

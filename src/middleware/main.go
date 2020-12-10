package middleware

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/nocubicles/skillbase.io/src/models"
	"github.com/nocubicles/skillbase.io/src/types"
	"github.com/nocubicles/skillbase.io/src/utils"
)

func CORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		next.ServeHTTP(w, r)
	})
}

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.RequestURI)
		next.ServeHTTP(w, r)
	})
}

func CheckCookie(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, _ := r.Cookie("skillbase")

		if cookie != nil {
			http.Redirect(w, r, "/app", http.StatusTemporaryRedirect)
			return
		}
		next(w, r)

	}
}

func CheckIsUsedLoggedIn(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("skillbase")
		if err != nil {
			currentPath := r.URL.Path
			if currentPath != "/" {
				http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
				return
			}

			next(w, r)
		}
		sessionID := cookie.Value

		db := utils.DbConnection()
		var session models.Session
		result := db.Where("session_id = ? AND expiration > ?", sessionID, time.Now()).First(&session)

		if result.RowsAffected > 0 {

			authContext := types.AuthContext{
				UserID:   uint(session.UserID),
				TenantID: uint(session.TenantID),
			}

			ctx := context.WithValue(r.Context(), "authContext", authContext)

			r := r.WithContext(ctx)
			if r.URL.Path == "/" {
				http.Redirect(w, r, "/app", http.StatusTemporaryRedirect)
			}
			next(w, r)
		} else {
			http.Redirect(w, r, "/", http.StatusTemporaryRedirect)

			return
		}
	}
}

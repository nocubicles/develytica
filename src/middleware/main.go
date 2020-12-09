package middleware

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/nocubicles/skillbase.io/src/models"
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

func CheckIsUsedLoggedIn(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("skillbase")
		if err != nil {
			utils.Render(w, "index.html", nil)

			return
		}
		sessionID := cookie.Value

		db := utils.DbConnection()
		var session models.Session
		result := db.Where("session_id = ? AND expiration > ?", sessionID, time.Now()).First(&session)

		if result.RowsAffected > 0 {

			ctx := context.WithValue(r.Context(), "userID", uint(session.UserID))
			r := r.WithContext(ctx)
			if r.URL.Path == "/" {
				http.Redirect(w, r, "/app", http.StatusTemporaryRedirect)
			}
			next(w, r)
		} else {
			utils.Render(w, "index.html", nil)

			return
		}
	}
}

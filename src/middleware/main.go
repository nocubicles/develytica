package middleware

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/nocubicles/develytica/src/models"
	"github.com/nocubicles/develytica/src/types"
	"github.com/nocubicles/develytica/src/utils"
)

func ForceHttps(next http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		http.Redirect(res, req, fmt.Sprintf("https://%s%s", req.Host, req.URL), http.StatusPermanentRedirect)
		return
	})
}

func CORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("x-forwarded-proto", "https")
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
		cookie, err := r.Cookie("develytica")
		if err != nil {
			http.Redirect(w, r, "/signin", http.StatusTemporaryRedirect)
			return
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
			next(w, r)
		} else {

			http.Redirect(w, r, "/signin", http.StatusTemporaryRedirect)
			return
		}
	}
}

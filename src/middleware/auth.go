package middleware

import (
	"context"
	"net/http"
	"time"

	"github.com/nocubicles/skillbase.io/src/models"
	"github.com/nocubicles/skillbase.io/src/utils"
)

func CheckIsUsedLoggedIn(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("skillbase.io")
		if err != nil {
			http.Redirect(w, r, "/logisisse", http.StatusTemporaryRedirect)

			return
		}
		sessionID := cookie.Value

		db := utils.DbConnection()
		var session models.Session
		result := db.Where("session_id = ? AND expiration > ?", sessionID, time.Now()).First(&session)

		if result.RowsAffected > 0 {
			ctx := context.WithValue(r.Context(), "userID", uint(session.UserID))
			r := r.WithContext(ctx)
			next(w, r)
		} else {
			http.Redirect(w, r, "/logisisse", http.StatusTemporaryRedirect)
			return
		}
	}
}

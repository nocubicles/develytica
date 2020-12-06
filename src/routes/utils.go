package routes

import (
	"net/http"
	"time"

	"github.com/gofrs/uuid"
	"github.com/nocubicles/skillbase.io/src/models"
	"github.com/nocubicles/skillbase.io/src/utils"
)

func setCookieForUser(w http.ResponseWriter, email string) error {
	expiration := time.Now().Add(14 * 24 * time.Hour)
	sessionID, err := uuid.NewV4()
	if err != nil {
		return err
	}

	cookie := http.Cookie{
		Name:     "skillbase.io",
		Value:    sessionID.String(),
		Expires:  expiration,
		HttpOnly: true,
	}

	http.SetCookie(w, &cookie)
	saveSession(email, sessionID, expiration)
	return nil
}

func saveSession(email string, sessionID uuid.UUID, expiration time.Time) {
	var user models.User
	db := utils.DbConnection()
	result := db.Where("Email = ?", email).First(&user)

	if result.RowsAffected > 0 {
		session := models.Session{
			Expiration: expiration,
			SessionID:  sessionID,
			UserID:     user.ID,
		}

		db.Create(&session)
	}
}

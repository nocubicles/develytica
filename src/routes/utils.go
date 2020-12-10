package routes

import (
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/gofrs/uuid"
	"github.com/nocubicles/skillbase.io/src/models"
	"github.com/nocubicles/skillbase.io/src/types"
	"github.com/nocubicles/skillbase.io/src/utils"
)

func setCookieForUser(w http.ResponseWriter, email string) error {
	expiration := time.Now().Add(14 * 24 * time.Hour)
	sessionID, err := uuid.NewV4()
	if err != nil {
		return err
	}

	cookie := http.Cookie{
		Name:     "skillbase",
		Value:    sessionID.String(),
		Expires:  expiration,
		HttpOnly: true,
		Path:     "/",
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
			TenantID:   user.TenantID,
		}

		db.Create(&session)
	}
}

func getAuthContextData(r *http.Request) (types.AuthContext, error) {
	authContext := types.AuthContext{}
	authContext, ok := r.Context().Value("authContext").(types.AuthContext)

	if !ok {
		log.Println("Getting auth context failed")
		return authContext, errors.New("Getting auth context failed")
	}
	return authContext, nil
}

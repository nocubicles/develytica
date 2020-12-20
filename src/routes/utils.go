package routes

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gofrs/uuid"
	"github.com/nocubicles/develytica/src/models"
	"github.com/nocubicles/develytica/src/types"
	"github.com/nocubicles/develytica/src/utils"
)

func setCookieForUser(w http.ResponseWriter, email string) error {
	expiration := time.Now().Add(14 * 24 * time.Hour)
	sessionID, err := uuid.NewV4()
	if err != nil {
		return err
	}

	cookie := http.Cookie{
		Name:     "develytica",
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

func convertStringToInt64(value string) int64 {

	u64, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		fmt.Println(err)
		return 0
	}

	return int64(u64)
}

func convertStringToUint(value string) uint {

	u64, err := strconv.ParseUint(value, 10, 64)
	if err != nil {
		fmt.Println(err)
		return 0
	}

	return uint(u64)
}

func getTeamMembers(tenantID uint, teamMembers *[]TeamMember, limit int) *[]TeamMember {
	db := utils.DbConnection()
	db.Raw(`select a.login, a.avatar_url, a.location, a.remote_id, count(ia.assignee_id) as issuescount
		from assignees a
		left join issue_assignees ia on ia.assignee_id = a.remote_id
		where a.tenant_id = ?
		group by a.login,a.avatar_url,a.location,a.remote_id
		order by issuescount desc
		limit ?

		`, tenantID, limit).
		Scan(&teamMembers)
	return teamMembers
}

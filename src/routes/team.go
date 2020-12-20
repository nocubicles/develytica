package routes

import (
	"fmt"
	"net/http"

	"github.com/nocubicles/develytica/src/models"
	"github.com/nocubicles/develytica/src/utils"
)

type TeamPageData struct {
	Authenticated    bool
	UserName         string
	TeamMembers      []TeamMember
	ValidationErrors map[string]string
}

func TeamHandler(w http.ResponseWriter, r *http.Request) {
	authContext, err := getAuthContextData(r)
	if err != nil {
		fmt.Println(err)
	}
	user := models.User{}
	db := utils.DbConnection()
	result := db.First(&user, authContext.UserID)
	teamMembers := []TeamMember{}
	data := TeamPageData{
		Authenticated:    false,
		UserName:         "",
		TeamMembers:      teamMembers,
		ValidationErrors: map[string]string{},
	}

	if result.RowsAffected > 0 {
		data.Authenticated = true
		data.UserName = user.Email
	}

	if r.Method == http.MethodGet {

		data.TeamMembers = *getTeamMembers(user.TenantID, &teamMembers, 100)
		utils.Render(w, "team.gohtml", data)

		return
	}

}

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

		skills, ok := r.URL.Query()["skills"]

		if !ok || len(skills[0]) < 1 {
			data.TeamMembers = *getAllTeamMembers(user.TenantID, &teamMembers, 100)
		} else if len(skills[0]) > 0 && checkSkillExist(user.TenantID, skills[0]) {
			data.TeamMembers = *getTeamMembersBySkillName(user.TenantID, &teamMembers, skills[0], 100)
		}
		utils.Render(w, "team.gohtml", data)

		return
	}

}

func checkSkillExist(tenantID uint, skillName string) bool {
	db := utils.DbConnection()
	issueLabels := models.IssueLabel{}
	result := db.Where("tenant_id = ? AND name = ?", tenantID, skillName).Find(&issueLabels).Limit(1)

	if result.RowsAffected > 0 {
		return true
	}

	return false
}

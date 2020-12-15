package routes

import (
	"fmt"
	"net/http"

	"github.com/nocubicles/skillbase.io/src/models"
	"github.com/nocubicles/skillbase.io/src/utils"
)

type TeamMember struct {
	Login       string
	AvatarURL   string
	Location    string
	IssuesCount int64 `gorm:"column:issuescount"`
}

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
		db.Raw(`select a.login,a.avatar_url, a.location, count(ia.assignee_id) as issuescount
		from assignees a
		left join issue_assignees ia on ia.assignee_id = a.remote_id
		where a.tenant_id = ?
		group by a.login,a.avatar_url,a.location`, user.TenantID).
			Scan(&teamMembers)
		data.TeamMembers = teamMembers
		utils.Render(w, "team.html", data)

		return
	}

}

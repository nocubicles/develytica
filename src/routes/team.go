package routes

import (
	"fmt"
	"net/http"

	"github.com/nocubicles/develytica/src/models"
	"github.com/nocubicles/develytica/src/utils"
)

type TeamMember struct {
	Login       string
	AvatarURL   string
	Location    string
	RemoteID    int64
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
		db.Raw(`select a.login, a.avatar_url, a.location, a.remote_id, count(ia.assignee_id) as issuescount
		from assignees a
		left join issue_assignees ia on ia.assignee_id = a.remote_id
		where a.tenant_id = ?
		group by a.login,a.avatar_url,a.location,a.remote_id
		order by issuescount desc
		`, user.TenantID).
			Scan(&teamMembers)
		data.TeamMembers = teamMembers
		utils.Render(w, "team.gohtml", data)

		return
	}

}

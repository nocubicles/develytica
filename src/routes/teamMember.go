package routes

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/nocubicles/skillbase.io/src/models"
	"github.com/nocubicles/skillbase.io/src/utils"
)

type UserSkill struct {
	SkillName string `gorm:"column:skillname"`
	DoneCount int    `gorm:"column:donecount"`
}

type TeamMemberData struct {
	Login       string
	AvatarURL   string
	Location    string
	RemoteID    int64
	IssuesCount int64 `gorm:"column:issuescount"`
	UserSkills  []UserSkill
}

type TeamMemberPageData struct {
	Authenticated    bool
	UserName         string
	TeamMemberData   TeamMemberData
	ValidationErrors map[string]string
}

func TeamMemberHandler(w http.ResponseWriter, r *http.Request) {
	authContext, err := getAuthContextData(r)
	if err != nil {
		fmt.Println(err)
	}
	user := models.User{}
	db := utils.DbConnection()
	result := db.First(&user, authContext.UserID)
	teamMember := TeamMemberData{}
	data := TeamMemberPageData{
		Authenticated:    false,
		UserName:         "",
		TeamMemberData:   teamMember,
		ValidationErrors: map[string]string{},
	}

	if result.RowsAffected > 0 {
		data.Authenticated = true
		data.UserName = user.Email
	}

	vars := mux.Vars(r)
	teamMemberID := convertStringToInt64(vars["teamMember"])

	if r.Method == http.MethodGet {
		db.Raw(`select a.login, a.avatar_url, a.location, a.remote_id, count(ia.assignee_id) as issuescount
		from assignees a
		left join issue_assignees ia on ia.assignee_id = a.remote_id
		where a.tenant_id = ? AND a.remote_id = ?
		group by a.login,a.avatar_url,a.location,a.remote_id
		order by issuescount desc
		`, user.TenantID, teamMemberID).
			Scan(&teamMember)
		data.TeamMemberData = teamMember

		userSkills := []UserSkill{}
		db.Raw(`
			SELECT labels.name as skillname, 
			issue_assignees.assignee_id as assigneeID, 
			count(issue_labels.label_id) as donecount
			FROM issue_assignees
			LEFT JOIN issue_labels on issue_labels.issue_id = issue_assignees.issue_id
			LEFT JOIN labels on labels.label_id = issue_labels.label_id
			WHERE issue_assignees.tenant_id = ? AND issue_assignees.assignee_id = ? AND labels.tracked = true
			group by skillname, assigneeID
			order by donecount desc
		`, user.TenantID, teamMemberID).
			Scan(&userSkills)

		data.TeamMemberData.UserSkills = userSkills
		utils.Render(w, "teamMember.html", data)

		return
	}

}

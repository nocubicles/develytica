package routes

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/nocubicles/develytica/src/models"
	"github.com/nocubicles/develytica/src/utils"
)

type UserSkill struct {
	SkillName       string    `gorm:"column:skillname"`
	DoneCount       int       `gorm:"column:donecount"`
	LastUsed        time.Time `gorm:"column:lastused"`
	LastUsedDaysAgo int
}

type TeamMemberData struct {
	Login       string
	AvatarURL   string
	Location    string
	RemoteID    int64
	IssuesCount int64 `gorm:"column:issuescount"`
}

type TeamMemberPageData struct {
	Authenticated    bool
	UserName         string
	TeamMemberData   TeamMemberData
	ValidationErrors map[string]string
	UserSkills       []UserSkill
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
		FROM assignees a
		LEFT JOIN issue_assignees ia on ia.assignee_id = a.remote_id
		WHERE a.tenant_id = ? AND a.remote_id = ?
		GROUP BY a.login,a.avatar_url,a.location,a.remote_id
		ORDER BY issuescount desc
		`, user.TenantID, teamMemberID).
			Scan(&teamMember)
		data.TeamMemberData = teamMember

		userSkills := []UserSkill{}
		db.Raw(`
			SELECT 
			issue_labels.name as skillname,
			count(issue_assignees.issue_id) as donecount,
			issues.closed_at as lastused
			FROM issue_assignees
			LEFT JOIN issue_labels ON issue_labels.issue_id = issue_assignees.issue_id
			LEFT JOIN label_trackings ON label_trackings.name = issue_labels.name
			LEFT JOIN issues ON issues.remote_id = issue_assignees.issue_id
			WHERE issue_assignees.tenant_id = ?
			AND issue_assignees.assignee_id = ?
			AND label_trackings.is_tracked = true
			GROUP BY skillname, lastused
			ORDER BY donecount desc
		`, user.TenantID, teamMemberID).
			Scan(&userSkills)

		for i := range userSkills {
			userSkills[i].LastUsedDaysAgo = daysBetween(time.Now(), userSkills[i].LastUsed)
		}

		data.UserSkills = userSkills
		utils.Render(w, "teamMember.gohtml", data)

		return
	}

}

func daysBetween(a, b time.Time) int {
	if a.After(b) {
		a, b = b, a
	}

	days := -a.YearDay()
	for year := a.Year(); year < b.Year(); year++ {
		days += time.Date(year, time.December, 31, 0, 0, 0, 0, time.UTC).YearDay()
	}
	days += b.YearDay()

	return days
}

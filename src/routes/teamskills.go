package routes

import (
	"fmt"
	"net/http"
	"time"

	"github.com/nocubicles/develytica/src/models"
	"github.com/nocubicles/develytica/src/utils"
)

type Skill struct {
	SkillName       string    `gorm:"column:skillname"`
	DoneCount       int       `gorm:"column:donecount"`
	LastUsed        time.Time `gorm:"column:lastused"`
	LastUsedDaysAgo int
}

type TeamSkillPageData struct {
	Authenticated    bool
	UserName         string
	ValidationErrors map[string]string
	Skills           []Skill
}

func TeamSkillsHandler(w http.ResponseWriter, r *http.Request) {
	authContext, err := getAuthContextData(r)
	if err != nil {
		fmt.Println(err)
	}
	user := models.User{}
	db := utils.DbConnection()
	result := db.First(&user, authContext.UserID)
	data := TeamSkillPageData{
		Authenticated:    false,
		UserName:         "",
		ValidationErrors: map[string]string{},
	}
	if result.RowsAffected > 0 {
		data.Authenticated = true
		data.UserName = user.Email
	}

	if r.Method == http.MethodGet {
		Skills := []Skill{}

		db.Raw(`
			SELECT 
			issue_labels.name as skillname,
			count(issue_assignees.issue_id) as donecount
			FROM issue_assignees
			LEFT JOIN issue_labels ON issue_labels.issue_id = issue_assignees.issue_id
			LEFT JOIN label_trackings ON label_trackings.name = issue_labels.name
			WHERE issue_assignees.tenant_id = ?
			AND label_trackings.is_tracked = true
			GROUP BY skillname
			ORDER BY donecount desc
		`, user.TenantID).
			Scan(&Skills)

		for i := range Skills {
			Skills[i].LastUsed = getIssueClosedAtByLabelName(user.TenantID, Skills[i].SkillName)
			Skills[i].LastUsedDaysAgo = daysBetween(time.Now(), Skills[i].LastUsed)
		}

		data.Skills = Skills
		utils.Render(w, "teamskills.gohtml", data)

		return
	}
}

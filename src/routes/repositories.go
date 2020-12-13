package routes

import (
	"fmt"
	"net/http"

	"github.com/nocubicles/skillbase.io/src/models"
	"github.com/nocubicles/skillbase.io/src/utils"
)

type RepoData struct {
	OrgName         string
	Name            string
	IsTracked       bool
	OpenIssuesCount int
}

type RepoPageData struct {
	Authenticated    bool
	UserName         string
	ReposData        []RepoData
	ReposNotFound    bool
	ValidationErrors map[string]string
}

func Repository(w http.ResponseWriter, r *http.Request) {
	authContext, err := getAuthContextData(r)
	if err != nil {
		fmt.Println(err)
	}
	user := models.User{}
	db := utils.DbConnection()
	result := db.First(&user, authContext.UserID)

	data := RepoPageData{
		Authenticated:    false,
		UserName:         "",
		ReposData:        []RepoData{},
		ReposNotFound:    true,
		ValidationErrors: map[string]string{},
	}

	if result.RowsAffected > 0 {
		data.Authenticated = true
		data.UserName = user.Email
	}

	if r.Method == http.MethodPost {

	}

	if r.Method == http.MethodGet {
		result := []RepoData{}

		db.Model(&models.Organization{}).
			Select("organizations.login, repos.name, repos.open_issues_count, repo_trackings.id").
			Joins(`
				LEFT JOIN repos on organizations.remote_id = repos.remote_org_id 
				LEFT JOIN repo_trackings ON organizations.user_id = repo_trackings.user_id 
				AND organizations.tenant_id = repo_trackings.tenant_id 
				AND repos.remote_id = repo_trackings.repo_id`).
			Where("organizations.user_id = ? AND organizations.tenant_id = ?", user.ID, user.TenantID).
			Scan(&result)
		data.ReposData = result

		utils.Render(w, "repositories.html", data)
		return
	}

}

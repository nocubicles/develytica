package routes

import (
	"fmt"
	"net/http"

	"github.com/nocubicles/develytica/src/models"
	"github.com/nocubicles/develytica/src/services"
	"github.com/nocubicles/develytica/src/utils"
)

type RepoData struct {
	Login           string
	Name            string
	OpenIssuesCount int
	IsTracked       bool
	RemoteID        int64
}

type RepoPageData struct {
	Authenticated    bool
	UserName         string
	ReposData        []RepoData
	ReposNotFound    bool
	ValidationErrors map[string]string
}

func RepoHandler(w http.ResponseWriter, r *http.Request) {
	authContext, err := getAuthContextData(r)
	if err != nil {
		fmt.Println(err)
		return
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

	if r.Method == http.MethodPut {
		if err := r.ParseForm(); err != nil {
			fmt.Println(err)
			return
		}

		for key, values := range r.PostForm {

			if key == "repoTracked" && len(values) > 0 {
				repoTracking := models.RepoTracking{}
				reposTrackings := []models.RepoTracking{}

				db.Where("tenant_id = ?", user.TenantID).Delete(&repoTracking)

				for i := range values {
					repoTracking.RepoID = convertStringToInt64(values[i])
					repoTracking.TenantID = user.TenantID
					repoTracking.IsTracked = true
					reposTrackings = append(reposTrackings, repoTracking)
				}
				db.Create(&reposTrackings)
				go services.DoImmidiateFullSyncByTenantID(user.TenantID)
			}
		}
		http.Redirect(w, r, "/repositories", http.StatusTemporaryRedirect)
	}

	if r.Method == http.MethodGet {

		data.ReposData = getReposData(user.ID, user.TenantID)

		utils.Render(w, "repositories.gohtml", data)
	}
}

func getReposData(userID uint, tenantID uint) (reposData []RepoData) {
	db := utils.DbConnection()

	result := []RepoData{}

	db.Raw(`SELECT organizations.login, repos.name, repos.open_issues_count, repo_trackings.is_tracked, repos.remote_id FROM "organizations" 
		LEFT JOIN repos on organizations.remote_id = repos.remote_org_id 
		LEFT JOIN repo_trackings ON organizations.tenant_id = repo_trackings.tenant_id 
		AND organizations.tenant_id = repo_trackings.tenant_id 
		AND repos.remote_id = repo_trackings.repo_id
		WHERE organizations.tenant_id = ?
		ORDER BY organizations.login desc`, tenantID).
		Scan(&result)
	return result
}

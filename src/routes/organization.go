package routes

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/nocubicles/skillbase.io/src/models"
	"github.com/nocubicles/skillbase.io/src/services"
	"github.com/nocubicles/skillbase.io/src/utils"
)

type OrgPageData struct {
	Authenticated    bool
	UserName         string
	Organizations    []models.Organization
	ValidationErrors map[string]string
}

func Organization(w http.ResponseWriter, r *http.Request) {
	authContext, err := getAuthContextData(r)
	if err != nil {
		fmt.Println(err)
	}
	user := models.User{}
	db := utils.DbConnection()
	result := db.First(&user, authContext.UserID)
	data := OrgPageData{
		Authenticated:    false,
		UserName:         "",
		Organizations:    []models.Organization{},
		ValidationErrors: map[string]string{},
	}
	organizations := []models.Organization{}

	if result.RowsAffected > 0 {
		data.Authenticated = true
		data.UserName = user.Email
	}

	if r.Method == http.MethodPost {
		githubOrgName := r.PostFormValue("newOrg")

		if len(githubOrgName) < 1 || len(githubOrgName) > 200 {
			data.ValidationErrors = map[string]string{
				"wrongOrg": "Please enter valid org name",
			}
			utils.Render(w, "organizations.html", data)
			return
		}

		githubClient, ctx := utils.GetGithubClientByUserAndTenant(user.ID, user.TenantID)
		githubOrg, _, err := githubClient.Organizations.Get(ctx, githubOrgName)

		if err != nil {
			errorMessage := err.Error()

			if strings.Contains(errorMessage, "Not Found") {
				fmt.Println(err)
				data.ValidationErrors = map[string]string{
					"orgNotFound": "Organization not found",
				}
				utils.Render(w, "organizations.html", data)
				return
			}
			data.ValidationErrors = map[string]string{
				"githubApiError": "Github API Error",
			}
			utils.Render(w, "organizations.html", data)
			return

		}

		if githubOrg.GetPublicRepos() == 0 {
			data.ValidationErrors = map[string]string{
				"noPublicRepos": "This organization has no public repos. Therefore no point in syncing it",
			}
			utils.Render(w, "organizations.html", data)
			return
		}

		if len(githubOrg.GetLogin()) > 0 {
			services.SyncGithubOrganization(githubOrg, user.ID, user.TenantID, true)
			orgResult := db.Where("user_id = ? AND tenant_ID = ?", user.ID, user.TenantID).Find(&organizations)

			if orgResult.RowsAffected > 0 {
				data.Organizations = organizations
			}
			utils.Render(w, "organizations.html", data)
			return
		}

	}

	if r.Method == http.MethodGet {

		orgResult := db.Where("user_id = ? AND tenant_ID = ?", user.ID, user.TenantID).Find(&organizations)

		if orgResult.RowsAffected > 0 {
			data.Organizations = organizations
		}
		utils.Render(w, "organizations.html", data)
		return
	}

}

package routes

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/nocubicles/develytica/src/models"
	"github.com/nocubicles/develytica/src/services"
	"github.com/nocubicles/develytica/src/utils"
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
			utils.Render(w, "organizations.gohtml", data)
			return
		}

		if utils.IsOrgLimitOver(user.TenantID) {
			data.ValidationErrors = map[string]string{
				"limit": "Organization limit exceeded",
			}
			utils.Render(w, "organizations.gohtml", data)
			return
		}

		githubClient, ctx := utils.GetGithubClientByTenant(user.TenantID)
		githubOrg, _, err := githubClient.Organizations.Get(ctx, githubOrgName)

		if err != nil {
			errorMessage := err.Error()

			if strings.Contains(errorMessage, "Not Found") {
				fmt.Println(err)
				data.ValidationErrors = map[string]string{
					"orgNotFound": "Organization not found",
				}
				utils.Render(w, "organizations.gohtml", data)
				return
			}
			data.ValidationErrors = map[string]string{
				"githubApiError": "Github API Error",
			}
			utils.Render(w, "organizations.gohtml", data)
			return

		}

		if githubOrg.GetPublicRepos() == 0 {
			data.ValidationErrors = map[string]string{
				"noPublicRepos": "This organization has no public repos. Therefore no point in syncing it",
			}
			utils.Render(w, "organizations.gohtml", data)
			return
		}

		if len(githubOrg.GetLogin()) > 0 {
			services.SyncGithubOrganization(db, githubOrg, user.TenantID, true)
			orgResult := db.Where("tenant_ID = ?", user.TenantID).Find(&organizations)

			if orgResult.RowsAffected > 0 {
				data.Organizations = organizations
			}
			utils.Render(w, "organizations.gohtml", data)
			return
		}

	}

	if r.Method == http.MethodGet {

		orgResult := db.Where("tenant_ID = ?", user.TenantID).Find(&organizations)

		if orgResult.RowsAffected > 0 {
			data.Organizations = organizations
		}
		utils.Render(w, "organizations.gohtml", data)
		return
	}

}

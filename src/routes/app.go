package routes

import (
	"fmt"
	"net/http"

	"github.com/nocubicles/develytica/src/models"
	"github.com/nocubicles/develytica/src/utils"
)

type HomePageData struct {
	Authenticated      bool
	UserName           string
	OrganizationsCount int64
	ReposCount         int64
	LabelsCount        int64
	UsersCount         int64
	TeamMembers        []TeamMember
}

func RenderApp(w http.ResponseWriter, r *http.Request) {
	authContext, err := getAuthContextData(r)

	if err != nil {
		fmt.Println(err)
	}

	user := models.User{}
	db := utils.DbConnection()
	result := db.First(&user, authContext.UserID)
	data := HomePageData{
		Authenticated: false,
		UserName:      "",
	}
	teamMembers := []TeamMember{}
	if result.RowsAffected > 0 {
		data.Authenticated = true
		data.UserName = user.Email
		data.TeamMembers = teamMembers
	}
	var OrgCount int64
	db.Model(&models.Organization{}).Count(&OrgCount)
	var ReposCount int64
	db.Model(&models.Repo{}).Count(&ReposCount)
	var LabelsCount int64
	db.Model(&models.Label{}).Count(&LabelsCount)
	var UsersCount int64
	db.Model(&models.Assignee{}).Count(&UsersCount)

	data.LabelsCount = LabelsCount
	data.ReposCount = ReposCount
	data.OrganizationsCount = OrgCount
	data.UsersCount = UsersCount
	data.TeamMembers = *getAllTeamMembers(user.TenantID, &teamMembers, 10)

	utils.Render(w, "app.gohtml", data)
}

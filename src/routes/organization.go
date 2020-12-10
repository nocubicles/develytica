package routes

import (
	"fmt"
	"net/http"

	"github.com/nocubicles/skillbase.io/src/models"
	"github.com/nocubicles/skillbase.io/src/utils"
)

type OrgPageData struct {
	Authenticated bool
	UserName      string
	Organizations []models.Organization
}

func ReceiveNewOrganization(w http.ResponseWriter, r *http.Request) {
	authContext, err := getAuthContextData(r)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(authContext.UserID)
	fmt.Println(authContext.TenantID)

}

func RenderOrganizationPage(w http.ResponseWriter, r *http.Request) {
	data := OrgPageData{
		Authenticated: false,
		UserName:      "",
		Organizations: []models.Organization{},
	}
	authContext, err := getAuthContextData(r)
	if err != nil {
		fmt.Println(err)
	}
	user := models.User{}
	db := utils.DbConnection()
	result := db.First(&user, authContext.UserID)

	if result.RowsAffected > 0 {
		data.Authenticated = true
		data.UserName = user.Email
	}
	organizations := []models.Organization{}
	orgResult := db.Where("user_id = ? AND tenant_ID = ?", user.ID, user.TenantID).Find(&organizations)

	if orgResult.RowsAffected > 0 {
		data.Organizations = organizations
	}
	utils.Render(w, "organizations.html", data)
}

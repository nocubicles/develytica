package routes

import (
	"fmt"
	"net/http"

	"github.com/nocubicles/skillbase.io/src/models"
	"github.com/nocubicles/skillbase.io/src/services"
	"github.com/nocubicles/skillbase.io/src/utils"
)

func Sync(w http.ResponseWriter, r *http.Request) {
	authContext, err := getAuthContextData(r)

	if err != nil {
		fmt.Println(err)
	}

	user := models.User{}
	db := utils.DbConnection()
	result := db.First(&user, authContext.UserID)

	if result.RowsAffected > 0 {
		services.DoImmidiateFullSyncByUserTenantID(user.ID, user.TenantID)
	}
}

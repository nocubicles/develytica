package routes

import (
	"fmt"
	"net/http"

	"github.com/nocubicles/skillbase.io/src/models"
	"github.com/nocubicles/skillbase.io/src/services"
	"github.com/nocubicles/skillbase.io/src/utils"
)

type LabelData struct {
	Name        string
	Description string
	Color       string
	Tracked     bool
	RemoteID    int64
}

type LabelPageData struct {
	Authenticated    bool
	UserName         string
	LabelsData       []LabelData
	LabelsNotFound   bool
	ValidationErrors map[string]string
}

func LabelHandler(w http.ResponseWriter, r *http.Request) {
	authContext, err := getAuthContextData(r)
	if err != nil {
		fmt.Println(err)
		return
	}
	user := models.User{}
	db := utils.DbConnection()
	result := db.First(&user, authContext.UserID)

	data := LabelPageData{
		Authenticated:    false,
		UserName:         "",
		LabelsData:       []LabelData{},
		LabelsNotFound:   true,
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

			if key == "labelTracked" && len(values) > 0 {
				label := models.Label{}
				valuesInt64 := []int64{}
				db.Model(label).Where("user_id = ? AND tenant_id = ?", user.ID, user.TenantID).Update("tracked", false)

				for i := range values {
					valueInt64 := convertStringToInt64(values[i])
					valuesInt64 = append(valuesInt64, valueInt64)
				}

				db.Table("labels").Where("remote_id IN ?", valuesInt64).Updates(map[string]interface{}{"tracked": true})

				go services.DoImmidiateFullSyncByUserTenantID(user.ID, user.TenantID)
			}
		}
		http.Redirect(w, r, "/labels", http.StatusTemporaryRedirect)
	}

	if r.Method == http.MethodGet {
		result := []LabelData{}

		db.Model(&models.Label{}).
			Select("name, description, color, tracked, remote_id").
			Where("user_id = ? AND tenant_id = ?", user.ID, user.TenantID).
			Order("name desc").
			Scan(&result)
		data.LabelsData = result

		utils.Render(w, "labels.html", data)
		return
	}

}

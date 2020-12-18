package routes

import (
	"fmt"
	"net/http"

	"github.com/nocubicles/develytica/src/models"
	"github.com/nocubicles/develytica/src/services"
	"github.com/nocubicles/develytica/src/utils"
)

type LabelData struct {
	Name      string
	IsTracked bool
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
				labelTracking := models.LabelTracking{}
				labelsTrackings := []models.LabelTracking{}
				db.Where("tenant_id = ?", user.TenantID).Delete(&labelTracking)

				for i := range values {
					labelTracking.Name = values[i]
					labelTracking.TenantID = user.TenantID
					labelTracking.IsTracked = true
					labelsTrackings = append(labelsTrackings, labelTracking)
				}
				db.Create(&labelsTrackings)
				go services.DoImmidiateFullSyncByTenantID(user.TenantID)
			}
		}
		http.Redirect(w, r, "/labels", http.StatusTemporaryRedirect)
	}

	if r.Method == http.MethodGet {
		result := []LabelData{}

		db.Raw(`SELECT labels.name, label_trackings.is_tracked FROM labels 
		LEFT JOIN label_trackings ON label_trackings.tenant_id = labels.tenant_id 
		AND label_trackings.name = labels.name 
		WHERE labels.tenant_id = ?
		ORDER BY labels.name desc`, user.TenantID).
			Scan(&result)
		data.LabelsData = result

		utils.Render(w, "labels.gohtml", data)
		return
	}

}

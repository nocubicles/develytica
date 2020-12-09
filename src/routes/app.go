package routes

import (
	"net/http"

	"github.com/nocubicles/skillbase.io/src/models"
	"github.com/nocubicles/skillbase.io/src/utils"
)

type HomePageData struct {
	Authenticated bool
	UserName      string
}

func RenderApp(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID")
	if userID != nil {
		userID = userID.(uint)
	}
	user := models.User{}
	db := utils.DbConnection()
	result := db.First(&user, userID)
	data := HomePageData{
		Authenticated: false,
		UserName:      "",
	}
	if result.RowsAffected > 0 {
		data.Authenticated = true
		data.UserName = user.Email
	}

	utils.Render(w, "app.html", data)
}

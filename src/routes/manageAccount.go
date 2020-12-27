package routes

import (
	"net/http"

	"github.com/nocubicles/develytica/src/utils"
)

type ManageAccountPageData struct {
	Authenticated    bool
	UserName         string
	ValidationErrors map[string]string
}

func RenderManageAccount(w http.ResponseWriter, r *http.Request) {
	manageAccountPageData := ManageAccountPageData{
		Authenticated: true,
	}
	utils.Render(w, "manageAccount.gohtml", manageAccountPageData)
}

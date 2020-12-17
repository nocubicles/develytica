package routes

import (
	"net/http"

	"github.com/nocubicles/skillbase.io/src/utils"
)

func RenderSignIn(w http.ResponseWriter, r *http.Request) {

	utils.Render(w, "signin.html", nil)

}

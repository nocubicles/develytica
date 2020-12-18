package routes

import (
	"net/http"

	"github.com/nocubicles/develytica/src/utils"
)

func RenderSignIn(w http.ResponseWriter, r *http.Request) {

	utils.Render(w, "signin.gohtml", nil)

}

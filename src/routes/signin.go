package routes

import (
	"fmt"
	"log"
	"net/http"
	"regexp"

	"github.com/gofrs/uuid"
	"github.com/nocubicles/skillbase.io/src/models"
	"github.com/nocubicles/skillbase.io/src/utils"
)

var emailRegex = regexp.MustCompile(".+@.+\\..+")

type SignInData struct {
	Email  string
	Errors map[string]string
}

func receiveToken(w http.ResponseWriter, r *http.Request, signInToken string) {
	db := utils.DbConnection()
	var user models.User

	result := db.Where("sign_in_token = ?", signInToken).First(&user)

	if result.RowsAffected > 0 {
		err := setCookieForUser(w, user.Email)
		if err != nil {
			fmt.Println(err)
		}
		http.Redirect(w, r, "/kuulutus", http.StatusTemporaryRedirect)
	}
}

func RenderSignIn(w http.ResponseWriter, r *http.Request) {
	signInToken := r.URL.Query().Get("signintoken")
	if signInToken != "" {
		receiveToken(w, r, signInToken)
	} else {
		defaultData := SignInData{
			Email: "",
		}
		utils.Render(w, "signin.html", defaultData)
	}
}

func SendSignInEmail(w http.ResponseWriter, r *http.Request) {
	signInData := &SignInData{
		Email: r.PostFormValue("email"),
	}

	if signInData.validate() == false {
		utils.Render(w, "signin.html", signInData)
		return
	}

	if err := signInData.send(); err != nil {
		log.Println(err)
		http.Error(w, "Vabandame. Emaili ei saa saata", http.StatusInternalServerError)
		return
	}
	utils.Render(w, "confirmation.html", nil)

}

func (signIndata *SignInData) validate() bool {
	signIndata.Errors = make(map[string]string)

	match := emailRegex.Match([]byte(signIndata.Email))

	if match == false {
		signIndata.Errors["Email"] = "Palun sisestage korrektne email."
	}

	return len(signIndata.Errors) == 0
}

func updateUserSignInToken(email string) error {
	var user models.User
	db := utils.DbConnection()
	newSignInToken, err := uuid.NewV4()

	if err != nil {
		return err
	}

	db.Model(&user).Where("Email = ?", email).Update("SignInToken", newSignInToken)

	return nil
}

func (signInData *SignInData) send() error {
	var user models.User
	db := utils.DbConnection()

	result := db.Where("Email = ?", signInData.Email).First(&user)

	if result.RowsAffected > 0 {
		err := updateUserSignInToken(user.Email)

		db.Where("Email = ?", signInData.Email).First(&user)
		if err != nil {
			return err
		}
		err = utils.SendSignInEmail(user.Email, user.SignInToken.String())
		if err != nil {
			return err
		}
	} else {
		user := models.User{
			Email: signInData.Email,
		}
		result := db.Create(&user)

		if result.RowsAffected < 0 {
			err := utils.SendSignInEmail(user.Email, user.SignInToken.String())
			if err != nil {
				return err
			}
		}
	}
	return nil
}

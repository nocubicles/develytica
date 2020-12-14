package routes

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/nocubicles/skillbase.io/src/models"
	"github.com/nocubicles/skillbase.io/src/services"
	"github.com/nocubicles/skillbase.io/src/utils"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
)

var githubOauthConfig = &oauth2.Config{}

func init() {
	err := godotenv.Load(".env")

	if err != nil {
		panic("cannot load .env file")
	}

	githubOauthConfig.RedirectURL = "http://localhost:3000/auth/github/callback"
	githubOauthConfig.ClientID = os.Getenv("GITHUB_OAUTH_CLIENT_ID")
	githubOauthConfig.ClientSecret = os.Getenv("GITHUB_OAUTH_CLIENT_SECRET")
	githubOauthConfig.Scopes = []string{"repo", "user", "read:org"}
	githubOauthConfig.Endpoint = github.Endpoint
}

func GithubOauthLogin(w http.ResponseWriter, r *http.Request) {
	oauthState := generateStateOauthCookie(w)
	u := githubOauthConfig.AuthCodeURL(oauthState)
	http.Redirect(w, r, u, http.StatusTemporaryRedirect)
}

func GithubOauthCallback(w http.ResponseWriter, r *http.Request) {
	oauthState, err := r.Cookie("oauthstate")

	if err != nil {
		fmt.Println(err)
	}

	responseStateValue := r.URL.Query().Get("state")

	if responseStateValue != oauthState.Value {
		fmt.Println("Invalid github oauth state")
		return
	}

	code := r.URL.Query().Get("code")

	user, err := setupUserFromGithub(code)
	if err != nil {
		fmt.Println(err)
		return
	}

	err = setCookieForUser(w, user.Email)
	if err != nil {
		fmt.Println(err)
	}

	http.Redirect(w, r, "/app", http.StatusTemporaryRedirect)
}

func setupUserFromGithub(code string) (models.User, error) {
	ctx := context.Background()
	provider := "github"
	token, err := githubOauthConfig.Exchange(ctx, code)
	db := utils.DbConnection()
	user := models.User{}
	tenant := models.Tenant{}
	userClaim := models.UserClaim{}
	if err != nil {
		return user, fmt.Errorf("code exchange wrong: %s", err.Error())
	}

	githubClient, ctx := utils.GetGithubClientByToken(token.AccessToken)

	githubUser, _, err := githubClient.Users.Get(ctx, "")

	if err != nil {
		fmt.Println(err.Error())
	}

	result := db.Where("email = ?", *githubUser.Email).Find(&user)

	if result.RowsAffected > 0 {
		//existing user, new login, update token in claims

		result = db.Model(&userClaim).Where("user_id = ? AND provider = ? AND tenant_id = ?", user.ID, provider, user.TenantID).Update("access_token", token.AccessToken)
		if result.RowsAffected > 0 {
			services.UpdateSyncJobs(user.TenantID)
			go services.DoImmidiateFullSyncByTenantID(user.TenantID)
		}
	} else {
		//new user, never seen before
		db.Create(&tenant)
		user.Email = *githubUser.Email
		user.TenantID = tenant.ID
		db.Create(&user)
		userClaim.AccessToken = token.AccessToken
		userClaim.UserID = user.ID
		userClaim.Provider = provider
		userClaim.TenantID = user.TenantID

		db.Create(&userClaim)
		services.CreateSyncJobs(user.TenantID)
		go services.DoImmidiateFullSyncByTenantID(user.TenantID)
	}
	return user, nil
}

func generateStateOauthCookie(w http.ResponseWriter) string {
	var expiration = time.Now().Add(20 * time.Minute)

	b := make([]byte, 16)
	rand.Read(b)
	state := base64.URLEncoding.EncodeToString(b)
	cookie := http.Cookie{Name: "oauthstate", Value: state, Expires: expiration, HttpOnly: true}
	http.SetCookie(w, &cookie)

	return state
}

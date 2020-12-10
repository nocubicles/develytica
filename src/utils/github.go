package utils

import (
	"context"

	"github.com/google/go-github/v33/github"
	"github.com/nocubicles/skillbase.io/src/models"
	"golang.org/x/oauth2"
)

//GetGithubClientByUserAndTenant returns client to use for accessing github API
func GetGithubClientByUserAndTenant(userID uint, tenantID uint) (*github.Client, context.Context) {
	userClaim := models.UserClaim{}
	db.Where("tenant_id = ? AND user_id = ?", tenantID, userID).Find(&userClaim)

	ctx := context.Background()
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: userClaim.AccessToken})

	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)

	return client, ctx
}

//GetGithubClientByToken returns client to use for accessing github API
func GetGithubClientByToken(accessToken string) (*github.Client, context.Context) {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: accessToken})

	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)

	return client, ctx
}

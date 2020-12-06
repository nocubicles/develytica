package utils

import (
	"context"

	"github.com/google/go-github/v33/github"
	"golang.org/x/oauth2"
)

//GetGithubClient returns client to use for accessing github API
func GetGithubClient(accessToken string) *github.Client {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: accessToken})

	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)

	return client
}

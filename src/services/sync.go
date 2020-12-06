package services

import (
	"context"
	"fmt"
	"time"

	"github.com/google/go-github/v33/github"
	"github.com/nocubicles/skillbase.io/src/models"
	"github.com/nocubicles/skillbase.io/src/utils"
	"gorm.io/gorm/clause"
)

var neededSyncs = []string{
	"organization",
	"repo",
	"issue",
}

//UpdateSyncJobs ensures that user has syncs enabled
func UpdateSyncJobs(userID uint, tenantID uint) {
	db := utils.DbConnection()

	for _, neededSync := range neededSyncs {
		sync := models.Sync{}
		sync.TenantID = tenantID
		sync.Name = neededSync
		sync.UserID = userID
		db.FirstOrCreate(&sync, models.Sync{Name: neededSync})

	}
	DoSync()
}

//CreateSyncJobs creates needed syncs for the user
func CreateSyncJobs(userID uint, tenantID uint) {

	db := utils.DbConnection()

	for _, neededSync := range neededSyncs {
		sync := models.Sync{}
		sync.TenantID = tenantID
		sync.Name = neededSync
		sync.UserID = userID
		db.Create(&sync)
	}
}

func DoSync() {
	syncs := []models.Sync{}
	syncInterval := 4 * time.Hour
	syncDateNeeded := time.Now().Add(-syncInterval)
	db := utils.DbConnection()

	result := db.Where("last_run < ?", syncDateNeeded).Find(&syncs)

	//all syncs that needs new sync
	if result.RowsAffected > 0 {
		for _, sync := range syncs {
			SyncGithubData(sync.UserID, sync.TenantID, sync.Name)
		}
	}
}

func checkIfRateLimitErr(err error) {
	if _, ok := err.(*github.RateLimitError); ok {
		fmt.Println("github request hit rate limit")
	}
}

func checkIfAcceptedError(err error) {
	if _, ok := err.(*github.AcceptedError); ok {
		fmt.Println("scheduled on github side")
	}
}

//SyncGithubData syncs user/tenant data from github
func SyncGithubData(userID uint, tenantID uint, syncName string) {
	start := time.Now()

	db := utils.DbConnection()
	ctx := context.Background()
	userClaim := models.UserClaim{}
	db.Where("tenant_id = ? AND user_id = ?", tenantID, userID).Find(&userClaim)

	githubClient := utils.GetGithubClient(userClaim.AccessToken)

	if syncName == "organization" {
		var allOrgs []*github.Organization
		options := &github.ListOptions{
			PerPage: 100,
		}

		for {
			orgs, response, err := githubClient.Organizations.List(ctx, "", options)
			if err != nil {
				checkIfRateLimitErr(err)
				checkIfAcceptedError(err)
				fmt.Println(err)
			}

			allOrgs = append(allOrgs, orgs...)
			if response.NextPage == 0 {
				break
			}
			options.Page = response.NextPage
		}

		for i := range allOrgs {
			org := models.GithubOrganization{}
			githubOrg := allOrgs[i]
			org.AvatarURL = githubOrg.GetAvatarURL()
			org.Collaborators = githubOrg.GetCollaborators()
			org.Company = githubOrg.GetCompany()
			org.Login = githubOrg.GetLogin()
			org.Name = githubOrg.GetName()
			org.GithubID = githubOrg.GetID()
			org.Type = githubOrg.GetType()
			org.Followers = githubOrg.GetFollowers()
			org.UserID = userID
			org.TenantID = tenantID

			db.Clauses(clause.OnConflict{
				UpdateAll: true,
			}).Create(&org)
		}

	} else if syncName == "repo" {

	} else if syncName == "issue" {

	}

	// orgIssues, _, _ := githubClient.Issues.ListByOrg(ctx, orgName, nil)

	// for i := range orgIssues {
	// 	labels := *&orgIssues[i].Labels

	// 	for x := range labels {
	// 		label := *labels[x]

	// 		fmt.Println(label)
	// 	}
	// }

	fmt.Printf("Sync duration: %v", time.Since(start).Milliseconds())

}

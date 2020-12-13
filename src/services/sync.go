package services

import (
	"fmt"
	"time"

	"github.com/google/go-github/v33/github"
	"github.com/nocubicles/skillbase.io/src/models"
	"github.com/nocubicles/skillbase.io/src/utils"
	"gorm.io/gorm/clause"
)

type neededSync struct {
	Name     string
	Priority int
}

var neededSyncs = []neededSync{
	{
		Name:     "organization",
		Priority: 1,
	},
	{
		Name:     "repo",
		Priority: 2,
	},
	{
		Name:     "issue",
		Priority: 3,
	},
}

//UpdateSyncJobs ensures that user has syncs enabled
func UpdateSyncJobs(userID uint, tenantID uint) {
	db := utils.DbConnection()

	for _, neededSync := range neededSyncs {
		sync := models.Sync{}
		sync.TenantID = tenantID
		sync.Name = neededSync.Name
		sync.Priority = neededSync.Priority
		sync.UserID = userID
		db.FirstOrCreate(&sync, models.Sync{Name: neededSync.Name})

	}
}

//CreateSyncJobs creates needed syncs for the user
func CreateSyncJobs(userID uint, tenantID uint) {

	db := utils.DbConnection()

	for _, neededSync := range neededSyncs {
		sync := models.Sync{}
		sync.TenantID = tenantID
		sync.Name = neededSync.Name
		sync.Priority = neededSync.Priority
		sync.UserID = userID
		db.Create(&sync)
	}
}

func DoImmidiateFullSyncByUserTenantID(userID uint, tenantID uint) {
	syncs := []models.Sync{}
	db := utils.DbConnection()

	result := db.Order("priority asc").Where("user_id = ? AND tenant_id = ?", userID, tenantID).Find(&syncs)

	if result.RowsAffected > 0 {
		for _, sync := range syncs {
			SyncGithubData(sync.UserID, sync.TenantID, sync.Name, sync.ID)
		}
	}
}

func DoFullSyncAllUsersPeriodic() {
	syncs := []models.Sync{}
	syncInterval := 4 * time.Hour
	syncDateNeeded := time.Now().Add(-syncInterval)
	db := utils.DbConnection()

	result := db.Order("priority asc").Where("last_run < ?", syncDateNeeded).Find(&syncs)

	//all syncs that needs new sync
	if result.RowsAffected > 0 {
		for _, sync := range syncs {
			SyncGithubData(sync.UserID, sync.TenantID, sync.Name, sync.ID)
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
func SyncGithubData(userID uint, tenantID uint, syncName string, syncID uint) {
	start := time.Now()
	perPage := 100
	db := utils.DbConnection()

	githubClient, ctx := utils.GetGithubClientByUserAndTenant(userID, tenantID)

	if syncName == "organization" {
		syncHistory := initiateSyncHistory(userID, tenantID, syncID)

		var allOrgs []*github.Organization
		options := &github.ListOptions{
			PerPage: perPage,
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
			SyncGithubOrganization(allOrgs[i], userID, tenantID, false)

		}
		finishSyncHistory(syncHistory)

	} else if syncName == "repo" {
		syncHistory := initiateSyncHistory(userID, tenantID, syncID)
		var allUserOrgs []models.Organization
		options := &github.RepositoryListByOrgOptions{ListOptions: github.ListOptions{PerPage: perPage}}
		db.Where("user_id = ? AND tenant_id = ? ", userID, tenantID).Find(&allUserOrgs)

		for _, userOrg := range allUserOrgs {
			var allRepos []*github.Repository

			for {
				repos, response, err := githubClient.Repositories.ListByOrg(ctx, userOrg.Login, options)

				if err != nil {
					checkIfRateLimitErr(err)
					checkIfAcceptedError(err)
					fmt.Println(err)
				}

				allRepos = append(allRepos, repos...)

				if response.NextPage == 0 {
					break
				}
				options.Page = response.NextPage
			}
			for i := range allRepos {
				repo := models.Repo{}
				githubRepo := allRepos[i]
				repo.Archived = githubRepo.GetArchived()
				repo.DefaultBranch = githubRepo.GetDefaultBranch()
				repo.Description = githubRepo.GetDescription()
				repo.Disabled = githubRepo.GetDisabled()
				repo.FullName = githubRepo.GetFullName()
				repo.RemoteID = githubRepo.GetID()
				repo.HTMLURL = githubRepo.GetHTMLURL()
				repo.HasIssues = githubRepo.GetHasIssues()
				repo.HasProjects = githubRepo.GetHasProjects()
				repo.Homepage = githubRepo.GetHomepage()
				repo.MasterBranch = githubRepo.GetMasterBranch()
				repo.Name = githubRepo.GetName()
				repo.OpenIssuesCount = githubRepo.GetOpenIssuesCount()
				repo.RemoteOrgID = userOrg.RemoteID
				repo.Private = githubRepo.GetPrivate()
				repo.PushedAt = githubRepo.GetPushedAt().Time
				repo.Size = githubRepo.GetSize()
				repo.StargazersCount = githubRepo.GetStargazersCount()
				repo.SubscribersCount = githubRepo.GetSubscribersCount()
				repo.WatchersCount = githubRepo.GetWatchersCount()
				repo.UpdatedAt = githubRepo.GetUpdatedAt().Time

				db.Clauses(clause.OnConflict{
					UpdateAll: true,
				}).Create(&repo)
			}
		}

		finishSyncHistory(syncHistory)
	} else if syncName == "issue" {
		syncHistory := initiateSyncHistory(userID, tenantID, syncID)
		options := &github.IssueListByRepoOptions{ListOptions: github.ListOptions{PerPage: perPage}, State: "all"}
		userOrgs := []models.Organization{}
		repoTracking := models.RepoTracking{}
		db.Where("user_id = ? AND tenant_id = ?", userID, tenantID).Find(&userOrgs)

		for _, userOrg := range userOrgs {
			userRepos := []models.Repo{}

			db.Where("remote_org_id = ?", userOrg.RemoteID).Find(&userRepos)

			for _, userRepo := range userRepos {
				if !shouldSyncRepo(userID, tenantID, userRepo.RemoteID) {
					continue
				}

				repoTrackingResult := db.Where("user_id = ? AND tenant_id = ? AND repo_ID = ?", userID, tenantID, userRepo.RemoteID).Find(&repoTracking)
				if repoTrackingResult.RowsAffected > 0 {
					for {
						issues, response, err := githubClient.Issues.ListByRepo(ctx, userOrg.Login, userRepo.Name, options)

						if err != nil {
							checkIfRateLimitErr(err)
							checkIfAcceptedError(err)
							fmt.Println(err)
						}

						SyncGithubIssues(issues, userID, tenantID, userRepo.RemoteID)

						if response.NextPage == 0 {
							break
						}

						options.Page = response.NextPage
					}
				}
			}

		}

		finishSyncHistory(syncHistory)
	}

	fmt.Printf("Sync duration: %v", time.Since(start).Milliseconds())

}

func shouldSyncRepo(userID uint, tenantID uint, repoID int64) bool {
	db := utils.DbConnection()
	repoTracking := models.RepoTracking{}
	result := db.Where("user_id = ? AND tenant_ID = ? and repo_ID = ?", userID, tenantID, repoID).First(&repoTracking)
	if result.RowsAffected > 0 {
		return true
	}
	return false
}

func syncLabelsFromIssue(userID uint, tenantID uint, issueID int64, RemoteIssueLabels []*github.Label) {
	db := utils.DbConnection()

	for i := range RemoteIssueLabels {
		remoteIssueLabel := RemoteIssueLabels[i]
		label := models.Label{}
		label.Color = remoteIssueLabel.GetColor()
		label.Description = remoteIssueLabel.GetDescription()
		label.RemoteID = remoteIssueLabel.GetID()
		label.Name = remoteIssueLabel.GetName()
		label.URL = remoteIssueLabel.GetURL()
		label.UserID = userID
		label.TenantID = tenantID
		db.Clauses(clause.OnConflict{
			UpdateAll: true,
		}).Create(&label)

		issueLabel := models.IssueLabel{}
		issueLabel.IssueID = issueID
		issueLabel.TenantID = tenantID
		issueLabel.UserID = userID
		db.FirstOrCreate(&issueLabel)
	}

}

func syncUsersFromIssue(userID uint, tenantID uint, issueID int64, RemoteIssueUsers []*github.User) {
	db := utils.DbConnection()

	for i := range RemoteIssueUsers {
		RemoteIssueUser := RemoteIssueUsers[i]
		assignee := models.Assignee{}
		assignee.RemoteID = RemoteIssueUser.GetID()
		assignee.AvatarURL = RemoteIssueUser.GetAvatarURL()
		assignee.Location = RemoteIssueUser.GetLocation()
		assignee.Login = RemoteIssueUser.GetLogin()
		assignee.Name = RemoteIssueUser.GetName()
		assignee.RemoteIssueID = issueID
		assignee.UserID = userID
		assignee.TenantID = tenantID
		db.Clauses(clause.OnConflict{
			UpdateAll: true,
		}).Create(&assignee)

		issueUser := models.IssueAssignee{}
		issueUser.AssigneeID = RemoteIssueUser.GetID()
		issueUser.IssueID = issueID
		issueUser.TenantID = tenantID
		issueUser.UserID = userID
		db.FirstOrCreate(&issueUser)
	}

}

func initiateSyncHistory(userID uint, tenantID uint, syncID uint) (syncHistoryID *models.SyncHistory) {
	db := utils.DbConnection()
	syncHistory := models.SyncHistory{}
	syncHistory.UserID = userID
	syncHistory.TenantID = tenantID
	syncHistory.Success = false
	syncHistory.SyncStart = time.Now()
	syncHistory.SyncID = syncID
	db.Create(&syncHistory)
	return &syncHistory
}

func finishSyncHistory(syncHistory *models.SyncHistory) {
	db := utils.DbConnection()
	syncHistory.Success = true
	syncHistory.SyncEnd = time.Now()
	db.Clauses(clause.OnConflict{UpdateAll: true}).Create(&syncHistory)
}

func SyncGithubIssues(issues []*github.Issue, userID uint, tenantID uint, repoID int64) {
	db := utils.DbConnection()
	for i := range issues {
		issue := models.Issue{}
		githubIssue := issues[i]
		if len(githubIssue.Labels) == 0 {
			continue
		}
		issue.RemoteRepoID = repoID
		issue.AssigneeID = githubIssue.GetAssignee().GetID()
		issue.AuthorAssociation = githubIssue.GetAuthorAssociation()
		issue.ClosedAt = githubIssue.GetClosedAt()
		issue.CreatedAt = githubIssue.GetCreatedAt()
		issue.ClosedByID = githubIssue.GetClosedBy().GetID()
		issue.RemoteID = githubIssue.GetID()
		issue.Locked = githubIssue.GetLocked()
		issue.Number = githubIssue.GetNumber()
		issue.RemoteUserID = githubIssue.GetUser().GetID()
		issue.State = githubIssue.GetState()
		issue.Title = githubIssue.GetTitle()

		db.Clauses(clause.OnConflict{
			UpdateAll: true,
		}).Create(&issue)
		syncLabelsFromIssue(userID, tenantID, issue.RemoteID, githubIssue.Labels)
		syncUsersFromIssue(userID, tenantID, issue.RemoteID, githubIssue.Assignees)

	}
}

//SyncGithubOrganization syncs github organization to db
func SyncGithubOrganization(githubOrg *github.Organization, userID uint, tenantID uint, manuallyAdded bool) {
	db := utils.DbConnection()
	org := models.Organization{}
	org.AvatarURL = githubOrg.GetAvatarURL()
	org.Collaborators = githubOrg.GetCollaborators()
	org.Company = githubOrg.GetCompany()
	org.Login = githubOrg.GetLogin()
	org.Name = githubOrg.GetName()
	org.RemoteID = githubOrg.GetID()
	org.Type = githubOrg.GetType()
	org.Followers = githubOrg.GetFollowers()
	org.UserID = userID
	org.TenantID = tenantID
	org.ManuallyAdded = manuallyAdded

	db.Clauses(clause.OnConflict{
		UpdateAll: true,
	}).Create(&org)
}

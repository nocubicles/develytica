package services

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/google/go-github/v33/github"
	"github.com/nocubicles/develytica/src/models"
	"github.com/nocubicles/develytica/src/utils"
	"gorm.io/gorm"
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
	{
		Name:     "label",
		Priority: 4,
	},
}

//UpdateSyncJobs ensures that tenant has syncs enabled
func UpdateSyncJobs(tenantID uint) {
	db := utils.DbConnection()

	for _, neededSync := range neededSyncs {
		sync := models.Sync{}
		sync.TenantID = tenantID
		sync.Name = neededSync.Name
		sync.Priority = neededSync.Priority
		db.FirstOrCreate(&sync, models.Sync{Name: neededSync.Name})

	}
}

//CreateSyncJobs creates needed syncs for the user
func CreateSyncJobs(tenantID uint) {

	db := utils.DbConnection()

	for _, neededSync := range neededSyncs {
		sync := models.Sync{}
		sync.TenantID = tenantID
		sync.Name = neededSync.Name
		sync.Priority = neededSync.Priority
		sync.InProgress = false
		db.Create(&sync)
	}
}

func ScanAndDoSyncs() {
	for {
		time.Sleep(5 * time.Second)
		go DoFullSyncAllUsersPeriodic()
	}
}

// DoImmidiateFullSyncByTenantID will start full sync cycle
func DoImmidiateFullSyncByTenantID(tenantID uint, db *gorm.DB) {
	syncs := []models.Sync{}

	result := db.Order("priority asc").Where("tenant_id = ?", tenantID).Find(&syncs)

	if result.RowsAffected > 0 {
		for _, sync := range syncs {
			if isSyncInProgress(sync) {
				continue
			}
			syncGithubData(sync.TenantID, sync.Name, sync.ID)
		}
	}
}

func DoFullSyncAllUsersPeriodic() {
	syncs := []models.Sync{}
	syncEveryHour, err := strconv.ParseInt(os.Getenv("SYNC_EVERY_HOUR"), 10, 64)
	var DEFAULT_SYNC_EVERY_HOUR = int64(2)

	if err != nil {
		syncEveryHour = DEFAULT_SYNC_EVERY_HOUR
	}

	syncInterval := time.Duration(syncEveryHour) * time.Hour
	syncDateNeeded := time.Now().Add(-syncInterval)
	db := utils.DbConnection()

	result := db.Order("priority asc").Where("last_run < ?", syncDateNeeded).Find(&syncs)

	//all syncs that needs new sync
	if result.RowsAffected > 0 {
		for _, sync := range syncs {
			if isSyncInProgress(sync) {
				continue
			}
			syncGithubData(sync.TenantID, sync.Name, sync.ID)
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

//syncGithubData syncs tenant data from github
func syncGithubData(tenantID uint, syncName string, syncID uint) {
	start := time.Now()
	perPage := 100
	db := utils.DbConnection()
	updateSyncInProgress(syncID, true)
	defer updateSyncInProgress(syncID, false)

	githubClient, ctx := utils.GetGithubClientByTenant(tenantID)

	if syncName == "organization" {
		syncHistory := initiateSyncHistory(tenantID, syncID)

		var allOrgs []*github.Organization
		options := &github.ListOptions{
			PerPage: perPage,
		}

		for {
			orgs, response, err := githubClient.Organizations.List(ctx, "", options)
			if err != nil {
				checkIfRateLimitErr(err)
				checkIfAcceptedError(err)
				utils.Logger.Warnw(err.Error(), "tenantID", tenantID)
				break
			}

			allOrgs = append(allOrgs, orgs...)
			if response.NextPage == 0 {
				break
			}
			options.Page = response.NextPage
		}

		for i := range allOrgs {
			SyncGithubOrganization(db, allOrgs[i], tenantID, false)
		}
		finishSyncHistory(syncHistory)

	} else if syncName == "repo" {
		syncHistory := initiateSyncHistory(tenantID, syncID)
		var allTenantOrgs []models.Organization
		options := &github.RepositoryListByOrgOptions{ListOptions: github.ListOptions{PerPage: perPage}}
		db.Where("tenant_id = ? ", tenantID).Find(&allTenantOrgs)

		for _, userOrg := range allTenantOrgs {
			var allRepos []*github.Repository

			for {
				repos, response, err := githubClient.Repositories.ListByOrg(ctx, userOrg.Login, options)

				if err != nil {
					checkIfRateLimitErr(err)
					checkIfAcceptedError(err)
					utils.Logger.Warnw(err.Error(), "tenantID", tenantID)
					break
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
		syncHistory := initiateSyncHistory(tenantID, syncID)
		options := &github.IssueListByRepoOptions{ListOptions: github.ListOptions{PerPage: perPage}, State: "closed"}
		tenantOrgs := []models.Organization{}
		db.Where("tenant_id = ?", tenantID).Find(&tenantOrgs)

		for _, tenantOrg := range tenantOrgs {
			tenantRepos := []models.Repo{}

			db.Select("name, remote_id").Where("remote_org_id = ?", tenantOrg.RemoteID).Find(&tenantRepos)

			for _, tenantRepo := range tenantRepos {
				if !shouldSyncRepo(tenantID, tenantRepo.RemoteID) {
					continue
				}
				repoTracking := models.RepoTracking{}
				repoTrackingResult := db.Where("tenant_id = ? AND repo_ID = ?", tenantID, tenantRepo.RemoteID).Find(&repoTracking)
				if repoTrackingResult.RowsAffected > 0 {
					for {
						issues, response, err := githubClient.Issues.ListByRepo(ctx, tenantOrg.Login, tenantRepo.Name, options)

						if err != nil {
							checkIfRateLimitErr(err)
							checkIfAcceptedError(err)
							utils.Logger.Warnw(err.Error(), "tenantID", tenantID)
							break
						}

						SyncGithubIssues(issues, tenantID, tenantRepo.RemoteID)

						if response.NextPage == 0 {
							break
						}

						options.Page = response.NextPage
					}
				}
			}

		}

		finishSyncHistory(syncHistory)
	} else if syncName == "label" {
		syncHistory := initiateSyncHistory(tenantID, syncID)
		options := &github.ListOptions{
			PerPage: perPage,
		}
		tenantOrgs := []models.Organization{}
		db.Where("tenant_id = ?", tenantID).Find(&tenantOrgs)

		for _, tenantOrg := range tenantOrgs {
			tenantRepos := []models.Repo{}
			db.Select("name, remote_id").Where("remote_org_id = ?", tenantOrg.RemoteID).Find(&tenantRepos)

			for _, tenantRepo := range tenantRepos {
				if !shouldSyncRepo(tenantID, tenantRepo.RemoteID) {
					continue
				}
				labels, response, err := githubClient.Repositories.ListLabels(ctx, tenantOrg.Login, tenantRepo.Name, options)

				if err != nil {
					checkIfRateLimitErr(err)
					checkIfAcceptedError(err)
					utils.Logger.Warnw(err.Error(), "tenantID", tenantID)
					break
				}

				syncLabels(tenantID, labels)

				if response.NextPage == 0 {
					break
				}

				options.Page = response.NextPage
			}
		}
		finishSyncHistory(syncHistory)
	}

	fmt.Printf("Sync duration: %v", time.Since(start).Milliseconds())

}

func shouldSyncRepo(tenantID uint, repoID int64) bool {
	db := utils.DbConnection()
	repoTracking := models.RepoTracking{}
	result := db.Where("tenant_ID = ? and repo_ID = ?", tenantID, repoID).First(&repoTracking)
	if result.RowsAffected > 0 {
		return true
	}
	return false
}

func syncLabels(tenantID uint, Labels []*github.Label) {
	db := utils.DbConnection()

	for i := range Labels {
		remoteIssueLabel := Labels[i]
		label := models.Label{}
		label.Name = remoteIssueLabel.GetName()
		label.TenantID = tenantID
		db.FirstOrCreate(&label)

	}

}

func syncLabelsFromIssue(tenantID uint, issueID int64, RemoteIssueLabels []*github.Label) {
	db := utils.DbConnection()

	for i := range RemoteIssueLabels {
		label := RemoteIssueLabels[i]
		issueLabel := models.IssueLabel{}
		issueLabel.Name = label.GetName()
		issueLabel.IssueID = issueID
		issueLabel.TenantID = tenantID
		db.FirstOrCreate(&issueLabel)
	}

}

func syncUsersFromIssue(tenantID uint, issueID int64, RemoteIssueUsers []*github.User) {
	db := utils.DbConnection()

	for i := range RemoteIssueUsers {
		RemoteIssueUser := RemoteIssueUsers[i]
		assignee := models.Assignee{}
		assignee.RemoteID = RemoteIssueUser.GetID()
		assignee.AvatarURL = RemoteIssueUser.GetAvatarURL()
		assignee.Location = RemoteIssueUser.GetLocation()
		assignee.Login = RemoteIssueUser.GetLogin()
		assignee.Name = RemoteIssueUser.GetName()
		assignee.TenantID = tenantID
		db.Clauses(clause.OnConflict{
			UpdateAll: true,
		}).Create(&assignee)

		issueAssignee := models.IssueAssignee{}
		issueAssignee.AssigneeID = RemoteIssueUser.GetID()
		issueAssignee.IssueID = issueID
		issueAssignee.TenantID = tenantID
		db.FirstOrCreate(&issueAssignee)
	}

}

func initiateSyncHistory(tenantID uint, syncID uint) (syncHistoryID *models.SyncHistory) {
	db := utils.DbConnection()
	syncHistory := models.SyncHistory{}
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

func isSyncInProgress(sync models.Sync) bool {
	if sync.InProgress {
		return true
	}

	return false
}

func updateSyncInProgress(syncID uint, status bool) {
	db := utils.DbConnection()
	sync := models.Sync{}
	db.Where("id = ?", syncID).First(&sync)
	sync.InProgress = status
	sync.LastRunSuccess = !sync.LastRunSuccess
	sync.LastRun = time.Now()
	db.Save(&sync)
}

func SyncGithubIssues(issues []*github.Issue, tenantID uint, repoID int64) {
	db := utils.DbConnection()
	for i := range issues {
		issue := models.Issue{}
		githubIssue := issues[i]
		if len(githubIssue.Labels) == 0 {
			continue
		}
		if len(githubIssue.Assignees) == 0 {
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
		syncLabelsFromIssue(tenantID, issue.RemoteID, githubIssue.Labels)
		syncUsersFromIssue(tenantID, issue.RemoteID, githubIssue.Assignees)
	}
}

//SyncGithubOrganization syncs github organization to db
func SyncGithubOrganization(db *gorm.DB, githubOrg *github.Organization, tenantID uint, manuallyAdded bool) {
	org := models.Organization{}
	org.AvatarURL = githubOrg.GetAvatarURL()
	org.Collaborators = githubOrg.GetCollaborators()
	org.Company = githubOrg.GetCompany()
	org.Login = githubOrg.GetLogin()
	org.Name = githubOrg.GetName()
	org.RemoteID = githubOrg.GetID()
	org.Type = githubOrg.GetType()
	org.Followers = githubOrg.GetFollowers()
	org.TenantID = tenantID
	org.ManuallyAdded = manuallyAdded

	db.Clauses(clause.OnConflict{
		UpdateAll: true,
	}).Create(&org)
}

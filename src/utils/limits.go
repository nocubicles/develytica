package utils

import (
	"github.com/nocubicles/develytica/src/models"
)

func GetAvailableRepoLimit(tenantID uint) int {
	db := DbConnection()

	tenantLimits := models.TenantLimit{}
	tenantLimitResultsresult := db.First(&tenantLimits, tenantID)

	if tenantLimitResultsresult.RowsAffected == 0 {
		tenantLimits.Repos = 10
	}

	return tenantLimits.Repos

}

func IsRepoLimitOver(tenantID uint) bool {
	db := DbConnection()

	tenantLimits := models.TenantLimit{}
	tenantLimitResultsresult := db.First(&tenantLimits, tenantID)
	repoTrackings := []models.RepoTracking{}

	if tenantLimitResultsresult.RowsAffected == 0 {
		tenantLimits.Repos = 10
	}

	repoResults := db.Where("tenant_id = ?", tenantID).Find(&repoTrackings)

	if repoResults.RowsAffected <= int64(tenantLimits.Org) {
		return false
	}
	return true
}

func IsOrgLimitOver(tenantID uint) bool {
	db := DbConnection()

	tenantLimits := models.TenantLimit{}
	tenantLimitResultsresult := db.First(&tenantLimits, tenantID)
	orgs := []models.Organization{}

	if tenantLimitResultsresult.RowsAffected == 0 {
		tenantLimits.Org = 3
	}

	orgResults := db.Where("tenant_id = ?", tenantID).Find(&orgs)

	if orgResults.RowsAffected <= int64(tenantLimits.Org) {
		return false
	}
	return true
}

func CreateTenantLimits(tenantID uint) {
	db := DbConnection()
	tenantLimits := models.TenantLimit{}

	tenantLimits.TenantID = tenantID
	tenantLimits.Org = 3
	tenantLimits.Repos = 10
	db.Create(&tenantLimits)
	return
}

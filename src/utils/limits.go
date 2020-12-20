package utils

import (
	"github.com/nocubicles/develytica/src/models"
)

func IsRepoLimitOver(tenantID uint) bool {
	db := DbConnection()

	tenantLimits := models.TenantLimit{}
	tenantLimitResultsresult := db.First(&tenantLimits, tenantID)
	repos := models.Repo{}

	if tenantLimitResultsresult.RowsAffected == 0 {
		tenantLimits.Repos = 10
	}

	repoResults := db.Where("tenant_id = ?", tenantID).Find(&repos)

	if repoResults.RowsAffected <= int64(tenantLimits.Org) {
		return false
	}
	return true
}

func IsOrgLimitOver(tenantID uint) bool {
	db := DbConnection()

	tenantLimits := models.TenantLimit{}
	tenantLimitResultsresult := db.First(&tenantLimits, tenantID)
	orgs := models.Organization{}

	if tenantLimitResultsresult.RowsAffected == 0 {
		tenantLimits.Repos = 10
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
	db.Create(&tenantLimits)
	return
}

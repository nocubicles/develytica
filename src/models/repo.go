package models

import (
	"time"
)

type Repo struct {
	UserID           uint `gorm:"primaryKey; not null"`
	TenantID         uint `gorm:"primaryKey; not null"`
	RemoteOrgID      int64
	RemoteID         int64 `gorm:"primaryKey; not null"`
	Name             string
	FullName         string
	Description      string
	Homepage         string
	DefaultBranch    string
	MasterBranch     string
	PushedAt         time.Time
	UpdatedAt        time.Time
	HTMLURL          string
	OpenIssuesCount  int
	StargazersCount  int
	SubscribersCount int
	WatchersCount    int
	Size             int
	Disabled         bool
	Archived         bool
	Private          bool
	HasIssues        bool
	HasProjects      bool
	TeamID           int64
}

package models

import (
	"github.com/jinzhu/gorm"
)

type GithubRepo struct {
	gorm.Model
	UserID        uint `gorm:"uniqueIndex:githubOrg"`
	TenantID      uint `gorm:"uniqueIndex:githubOrg"`
	Name          string
	Collaborators int
	Type          string
	Followers     int
	Location      string
	Company       string
	AvatarURL     string
	GithubID      int64 `gorm:"uniqueIndex:githubOrg"`
	Login         string
}

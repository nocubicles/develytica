package models

import (
	"github.com/jinzhu/gorm"
)

type GithubOrganization struct {
	gorm.Model
	UserID        uint `gorm:"index"`
	TenantID      uint
	Name          string
	Collaborators int
	Type          string
	Followers     int
	Location      string
	Company       string
	AvatarURL     string
	GithubID      int64
	Login         string
}

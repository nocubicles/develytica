package models

import (
	"github.com/jinzhu/gorm"
)

type GithubRepo struct {
	gorm.Model
	ID            uint `gorm:"primaryKey; autoIncrement:false; not null"`
	UserID        uint `gorm:"primaryKey; autoIncrement:false; not null"`
	TenantID      uint `gorm:"primaryKey; autoIncrement:false; not null"`
	Name          string
	Collaborators int
	Type          string
	Followers     int
	Location      string
	Company       string
	AvatarURL     string
	GithubID      int64 `gorm:"primaryKey; autoIncrement:false; not null"`
	Login         string
}

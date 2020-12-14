package models

import "github.com/jinzhu/gorm"

type RepoTracking struct {
	gorm.Model
	TenantID  uint  `gorm:"primaryKey; not null"`
	RepoID    int64 `gorm:"primaryKey; not null"`
	IsTracked bool
}

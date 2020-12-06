package models

import (
	"github.com/jinzhu/gorm"
)

type SyncHistory struct {
	gorm.Model
	UserID   uint `gorm:"index"`
	Name     string
	SyncID   uint `gorm:"index"`
	Success  bool
	TenantID uint
}

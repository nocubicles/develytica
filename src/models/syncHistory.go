package models

import (
	"time"

	"github.com/jinzhu/gorm"
)

type SyncHistory struct {
	gorm.Model
	UserID    uint `gorm:"index"`
	SyncID    uint `gorm:"index"`
	Success   bool
	TenantID  uint `gorm:"index"`
	SyncStart time.Time
	SyncEnd   time.Time
}

package models

import (
	"time"

	"github.com/jinzhu/gorm"
)

type Sync struct {
	gorm.Model
	TenantID       uint
	UserID         uint
	Name           string
	LastRunSuccess bool `gorm:"default:false"`
	LastRun        time.Time
}

package models

import (
	"time"

	"github.com/jinzhu/gorm"
)

type Sync struct {
	gorm.Model
	TenantID       uint   `gorm:"primaryKey"`
	Name           string `gorm:"primaryKey"`
	LastRunSuccess bool   `gorm:"default:false"`
	LastRun        time.Time
	Priority       int `gorm:"index"`
}

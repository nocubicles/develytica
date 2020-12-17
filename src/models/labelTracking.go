package models

import "github.com/jinzhu/gorm"

type LabelTracking struct {
	gorm.Model
	TenantID  uint `gorm:"primaryKey; not null"`
	IsTracked bool
	Name      string `gorm:"primaryKey; not null"`
}

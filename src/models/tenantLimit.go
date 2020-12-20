package models

import (
	"github.com/jinzhu/gorm"
)

type TenantLimit struct {
	gorm.Model
	TenantID uint `gorm:"primaryKey"`
	Org      int
	Repos    int
}

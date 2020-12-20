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

func (tenant *TenantLimit) BeforeCreate(tx *gorm.DB) (err error) {
	tenant.Org = 3
	tenant.Repos = 10
	return
}

package models

import (
	"github.com/jinzhu/gorm"
)

type UserClaim struct {
	gorm.Model
	UserID      uint   `gorm:"uniqueIndex:userClaim"`
	Provider    string `gorm:"uniqueIndex:userClaim"`
	AccessToken string
	TenantID    uint
}

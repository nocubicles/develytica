package models

import (
	"github.com/jinzhu/gorm"
)

type Tenant struct {
	gorm.Model
	Name          string
	Users         []User
	Syncs         []Sync
	Claims        []UserClaim
	SyncHistories []SyncHistory
	StripeID      string
}

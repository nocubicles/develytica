package models

import (
	"time"

	"github.com/gofrs/uuid"
	"github.com/jinzhu/gorm"
)

type Session struct {
	gorm.Model
	UserID uint
	SessionID uuid.UUID
	Expiration time.Time
}
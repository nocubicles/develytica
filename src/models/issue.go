package models

import "time"

type Issue struct {
	RemoteID          int64 `gorm:"primaryKey; not null"`
	RemoteRepoID      int64 `gorm:"primaryKey; not null"`
	Number            int
	State             string
	Locked            bool
	Title             string
	AuthorAssociation string
	RemoteUserID      int64
	AssigneeID        int64
	ClosedAt          time.Time
	CreatedAt         time.Time
	UpdatedAt         time.Time
	ClosedByID        int64
}

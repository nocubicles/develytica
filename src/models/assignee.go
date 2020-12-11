package models

type Assignee struct {
	RemoteID      int64 `gorm:"primaryKey; not null"`
	RemoteIssueID int64 `gorm:"primaryKey; not null"`
	Login         string
	Name          string
	AvatarURL     string
	Location      string
}

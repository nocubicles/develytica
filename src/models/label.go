package models

type Label struct {
	RemoteID      int64 `gorm:"primaryKey; not null"`
	RemoteIssueID int64 `gorm:"primaryKey; not null"`
	URL           string
	Name          string
	Color         string
	Description   string
	Tracked       bool `gorm:"index"`
}

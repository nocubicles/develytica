package models

type IssueAssignee struct {
	UserID     uint  `gorm:"primaryKey; not null"`
	TenantID   uint  `gorm:"primaryKey; not null"`
	AssigneeID int64 `gorm:"primaryKey; not null"`
	IssueID    int64 `gorm:"primaryKey; not null"`
}

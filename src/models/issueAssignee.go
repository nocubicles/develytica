package models

type IssueAssignee struct {
	TenantID   uint  `gorm:"primaryKey; not null"`
	AssigneeID int64 `gorm:"primaryKey; not null"`
	IssueID    int64 `gorm:"primaryKey; not null"`
}

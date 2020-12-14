package models

type IssueLabel struct {
	TenantID uint  `gorm:"primaryKey; not null"`
	LabelID  int64 `gorm:"primaryKey; not null"`
	IssueID  int64 `gorm:"primaryKey; not null"`
}

package models

type Assignee struct {
	UserID    uint  `gorm:"primaryKey; not null"`
	TenantID  uint  `gorm:"primaryKey; not null"`
	RemoteID  int64 `gorm:"primaryKey; not null"`
	Login     string
	Name      string
	AvatarURL string
	Location  string
}

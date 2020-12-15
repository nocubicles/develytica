package models

type Assignee struct {
	RemoteID  int64 `gorm:"primaryKey; not null"`
	Login     string
	Name      string
	AvatarURL string
	Location  string
	TenantID  uint `gorm:"primaryKey; not null"`
}

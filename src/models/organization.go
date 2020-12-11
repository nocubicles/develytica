package models

type Organization struct {
	UserID        uint `gorm:"primaryKey; not null"`
	TenantID      uint `gorm:"primaryKey; not null"`
	Name          string
	Collaborators int
	Type          string
	Followers     int
	Location      string
	Company       string
	AvatarURL     string
	RemoteID      int64 `gorm:"primaryKey; not null"`
	Login         string
	ManuallyAdded bool `gorm:"primaryKey"`
}

package models

type Label struct {
	UserID      uint  `gorm:"primaryKey; not null"`
	TenantID    uint  `gorm:"primaryKey; not null"`
	RemoteID    int64 `gorm:"primaryKey; not null"`
	URL         string
	Name        string
	Color       string
	Description string
	Tracked     bool `gorm:"index"`
}

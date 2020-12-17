package models

type Label struct {
	TenantID uint   `gorm:"primaryKey; not null"`
	Name     string `gorm:"primaryKey; not null"`
	Tracked  bool   `gorm:"index"`
	LabelID  int64  `gorm:"primaeryKey;not null"`
}

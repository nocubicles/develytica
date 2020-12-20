package routes

type TeamMember struct {
	Login       string
	AvatarURL   string
	Location    string
	RemoteID    int64
	IssuesCount int64 `gorm:"column:issuescount"`
}

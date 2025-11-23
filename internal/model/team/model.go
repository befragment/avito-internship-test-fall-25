package model

import "time"

type TeamMember struct {
	ID        string
	TeamID    string
	UserID    string
	CreatedAt time.Time
}

type Team struct {
	TeamID    string
	Name      string
	CreatedAt time.Time
	Members   []TeamMember
}

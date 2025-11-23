package model

import "time"

type User struct {
	UserID    string    `json:"user_id" db:"id"`
	Username  string    `json:"username" db:"username"`
	IsActive  bool      `json:"is_active" db:"is_active"`
	TeamName  string    `json:"team_name" db:"team_name"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

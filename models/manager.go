package models

//Manager is also a user
type Manager struct {
	FullName string `json:"name"`
	ID       int    `json:"id"`
	TeamID   int    `json:"team_id"`
}

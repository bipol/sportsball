package models

//RequestTransaction is what is sent to create a transaction
type RequestTransaction struct {
	Player   int64 `json:"player_id"`
}

//Transaction outlines a player transfer from one team to the next
type Transaction struct {
	FromTeam int64 `json:"from_team_id"`
	ToTeam   int64 `json:"to_team_id"`
	Player   int64 `json:"player_id"`
	ID       int64 `json:"id,omitempty"`
}

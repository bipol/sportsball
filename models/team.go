package models

//CreateTeamBody represents a group of players
type CreateTeamBody struct {
	Name    string   `json:"name"`
	Players []Player `json:"players"`
	Manager Manager  `json:"manager"`
}

//Todo: UUID for who can manage the team, accept transactions?

// Team represents a group of players
type Team struct {
	Name    string         `json:"name"`
	Players map[int]Player `json:"players"`
	Manager int64          `json:"manager"`
	ID      int64          `json:"id"`
}

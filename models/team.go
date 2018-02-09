package models

import (
	"database/sql"
)

// Team represents a group of players
type Team struct {
	Name    string         `json:"name"`
	Players map[int]Player `json:"players"`
	Manager int            `json:"manager"`
	ID      int            `json:"id"`
}

//AssignPlayer assigns a player to a team, and returns the team
func (team Team) AssignPlayer(player Player) (Team, error) {
	player.Team = team
	team.Players[player.ID] = player
	return team, nil
}

//AssignManager assigns a manager to a team, and returns the team
func (team Team) AssignManager(manager Manager) (Team, error) {
	manager.TeamID = team.ID
	team.Manager = manager.ID
	return team, nil
}

//Commit will save the team to the db
func (team Team) Commit(*sql.DB) (Team, error) {
	return team, nil
}

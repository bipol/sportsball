package models

import (
	"crypto/sha256"
	"hash"
	"sql"
)

//Transaction outlines a player transfer from one team to the next
type Transaction struct {
	FromTeam int       `json:"from_team_id"`
	ToTeam   int       `json:"to_team_id"`
	Player   int       `json:"player_id"`
	ID       hash.Hash `json:"id"`
}

//New instantiates a transaction
func (Transaction) New(fromTeamID int, toTeamID int, playerID int) (Transaction, error) {
	return Transaction{fromTeamID, toTeamID, playerID, sha256.New()}, nil
}

//Commit will save the transaction to the db
func (transaction Transaction) Commit(*sql.DB) (Transaction, error) {
	return transaction, nil
}

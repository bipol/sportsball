package storage

import (
	"database/sql"
	"fmt"
	"github.com/bipol/sportsball/config"
	"github.com/bipol/sportsball/models"
	_ "github.com/go-sql-driver/mysql"
)

//DatabaseContext contains the connection to the database
type DatabaseContext struct {
	Connection *sql.DB
}

//New instantiates a new database connection
func New(conf config.Config) (*DatabaseContext, error) {
	context := &DatabaseContext{}

	db, err := sql.Open("mysql", conf.DatabaseURL)

	if err != nil {
		return nil, err
	}

	err = db.Ping()

	if err != nil {
		return nil, err
	}

	context.Connection = db

	return context, nil
}

func createPlayerStatement(tx *sql.Tx) (*sql.Stmt, error) {
	statement := "INSERT INTO player (full_name, team_id, field_position) VALUES (?, ?, ?)"

	prepared, err := tx.Prepare(statement)

	return prepared, err
}

func createManagerStatement(tx *sql.Tx) (*sql.Stmt, error) {
	statement := "INSERT INTO manager (full_name, team_id) VALUES (?, ?)"

	prepared, err := tx.Prepare(statement)

	return prepared, err
}

func updateManagerStatement(tx *sql.Tx) (*sql.Stmt, error) {
	statement := "UPDATE team SET manager=? WHERE id=?"

	prepared, err := tx.Prepare(statement)

	return prepared, err
}

func (c *DatabaseContext) createTeamReturnID(tx *sql.Tx, name string) (int64, error) {
	res, err := tx.Exec("INSERT INTO team (name) VALUES (?)", name)

	if err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()

	return id, err
}

//CreateTeam will create a team in the database and return the pk
func (c *DatabaseContext) CreateTeam(team *models.CreateTeamBody) error {
	players := team.Players
	manager := team.Manager

	if team.Name == "" {
		return fmt.Errorf("Missing team name")
	}

	tx, err := c.Connection.Begin()
	defer tx.Rollback()

	if err != nil {
		return fmt.Errorf("Error creating transaction: %s", err)
	}

	id, err := c.createTeamReturnID(tx, team.Name)

	if err != nil {
		return fmt.Errorf("CreateTeamReturnID error: %s", err)
	}

	playerStatement, err := createPlayerStatement(tx)
	defer playerStatement.Close()

	if err != nil {
		return fmt.Errorf("createPlayerStatement error: %s", err)
	}

	for _, player := range players {
		_, err = playerStatement.Exec(player.FullName, id, player.Position)

		if err != nil {
			return fmt.Errorf("playerStatement error: %s", err)
		}
	}

	// todo: does this need to be prepared?
	managerStatement, err := createManagerStatement(tx)
	defer managerStatement.Close()

	if err != nil {
		return fmt.Errorf("Error creating manager statement: %s", err)
	}

	if manager.FullName == "" {
		return fmt.Errorf("Missing manager name")
	}

	managerRes, err := managerStatement.Exec(manager.FullName, id)

	if err != nil {
		return fmt.Errorf("Error executing manager statement: %s", err)
	}

	managerID, err := managerRes.LastInsertId()

	if err != nil {
		return fmt.Errorf("Could not retrieve manager id: %s", err)
	}

	updateManagerStatement, err := updateManagerStatement(tx)
	defer updateManagerStatement.Close()

	if err != nil {
		return fmt.Errorf("Error creating update manager statement: %s", err)
	}


	_, err = managerStatement.Exec(managerID, id)

	if err != nil {
		return fmt.Errorf("Error executing update manager statement: %s", err)
	}


	err = tx.Commit()

	if err != nil {
		return fmt.Errorf("Error commiting transaction: %s", err)
	}

	return err
}

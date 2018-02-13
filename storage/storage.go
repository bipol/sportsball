package storage

import (
	"database/sql"
	"fmt"
	"github.com/bipol/sportsball/config"
	"github.com/bipol/sportsball/models"
	_ "github.com/go-sql-driver/mysql"
	"github.com/satori/go.uuid"
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

//TODO: Implement
func (db *DatabaseContext) DoesTransactionAlreadyExist(transaction *models.Transaction) (error) {
}

//GetManagerByAccessKey will retrieve a manager from the database
func (db *DatabaseContext) CreateTransaction(transaction *models.Transaction) (error) {
	statement := "INSERT INTO transaction (from_team_id, to_team_id, player_id) VALUES (?, ?, ?)"

	_, err := db.Connection.Exec(statement, transaction.FromTeam, transaction.ToTeam, transaction.Player)

	return err
}

//GetManagerByAccessKey will retrieve a manager from the database
func (db *DatabaseContext) GetManagerByAccessKey(key string) (*models.Manager, error) {
	statement := fmt.Sprintf("SELECT id, full_name, team_id FROM manager WHERE access_key='%s'", key)

	row, err := db.Connection.Query(statement)

	if err != nil {
		return nil, err
	}

	defer row.Close()

	manager := models.Manager{}

	row.Next()
	err = row.Scan(&manager.ID, &manager.FullName, &manager.TeamID)

	return &manager, err
}

//GetPlayer returns a player from the database
func (db *DatabaseContext) GetPlayer(id int64) (*models.Player, error) {
	statement := fmt.Sprintf("SELECT id, full_name, team_id, field_position FROM player WHERE id=%d", id)

	row, err := db.Connection.Query(statement)

	if err != nil {
		return nil, err
	}

	defer row.Close()

	player := models.Player{}

	row.Next()
	err = row.Scan(&player.ID, &player.FullName, &player.Team, &player.Position)

	return &player, err
}


func createPlayerStatement(tx *sql.Tx) (*sql.Stmt, error) {
	statement := "INSERT INTO player (full_name, team_id, field_position) VALUES (?, ?, ?)"

	prepared, err := tx.Prepare(statement)

	return prepared, err
}

func createManagerStatement(tx *sql.Tx) (*sql.Stmt, error) {
	statement := "INSERT INTO manager (full_name, team_id, access_key) VALUES (?, ?, ?)"

	prepared, err := tx.Prepare(statement)

	return prepared, err
}

func associatePlayerWithTeamStatement(tx *sql.Tx) (*sql.Stmt, error) {
	statement := "INSERT INTO team_player (team, player) VALUES (?, ?)"

	prepared, err := tx.Prepare(statement)

	return prepared, err
}

func updateManagerStatement(tx *sql.Tx) (*sql.Stmt, error) {
	statement := "UPDATE team SET manager=? WHERE id=?"

	prepared, err := tx.Prepare(statement)

	return prepared, err
}

func (c *DatabaseContext) createTeam(tx *sql.Tx, name string) (int64, error) {
	res, err := tx.Exec("INSERT INTO team (name) VALUES (?)", name)

	if err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()

	return id, err
}

//CreateTeam will create a team in the database and return the access key for the user to use
func (c *DatabaseContext) CreateTeam(team *models.CreateTeamBody) (string, error) {
	players := team.Players
	manager := team.Manager

	if team.Name == "" {
		return "", fmt.Errorf("Missing team name")
	}

	tx, err := c.Connection.Begin()
	defer tx.Rollback()

	if err != nil {
		return "", fmt.Errorf("Error creating transaction: %s", err)
	}

	id, err := c.createTeam(tx, team.Name)

	if err != nil {
		return "", fmt.Errorf("CreateTeamReturnID error: %s", err)
	}

	playerStatement, err := createPlayerStatement(tx)
	defer playerStatement.Close()

	if err != nil {
		return "", fmt.Errorf("createPlayerStatement error: %s", err)
	}

	associatePlayerStatement, err := associatePlayerWithTeamStatement(tx)
	defer associatePlayerStatement.Close()

	if err != nil {
		return "", fmt.Errorf("associatePlayerStatement error: %s", err)
	}

	for _, player := range players {
		playerRes, err := playerStatement.Exec(player.FullName, id, player.Position)

		if err != nil {
			return "", fmt.Errorf("playerStatement error: %s", err)
		}

		playerID, err := playerRes.LastInsertId()

		if err != nil {
			return "", fmt.Errorf("playerID error: %s", err)
		}

		_, err = associatePlayerStatement.Exec(id, playerID)

		if err != nil {
			return "", fmt.Errorf("playerID error: %s", err)
		}
	}

	// todo: does this need to be prepared?
	managerStatement, err := createManagerStatement(tx)
	defer managerStatement.Close()

	if err != nil {
		return "", fmt.Errorf("Error creating manager statement: %s", err)
	}

	if manager.FullName == "" {
		return "", fmt.Errorf("Missing manager name")
	}

	accessKey := uuid.NewV4().String()

	managerRes, err := managerStatement.Exec(manager.FullName, id, accessKey)

	if err != nil {
		return "", fmt.Errorf("Error executing manager statement: %s", err)
	}

	managerID, err := managerRes.LastInsertId()

	if err != nil {
		return "", fmt.Errorf("Could not retrieve manager id: %s", err)
	}

	updateManagerStatement, err := updateManagerStatement(tx)
	defer updateManagerStatement.Close()

	if err != nil {
		return "", fmt.Errorf("Error creating update manager statement: %s", err)
	}

	_, err = updateManagerStatement.Exec(managerID, id)

	if err != nil {
		return "", fmt.Errorf("Error executing update manager statement: %s", err)
	}

	err = tx.Commit()

	if err != nil {
		return "", fmt.Errorf("Error commiting transaction: %s", err)
	}

	return accessKey, err
}

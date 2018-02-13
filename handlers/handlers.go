package handlers

import (
	"encoding/json"
	"github.com/bipol/sportsball/context"
	"github.com/bipol/sportsball/models"
	"goji.io"
	"goji.io/pat"
	"net/http"
	"time"
)

//APIMux generates a submux for all the api endpoints
func APIMux(appContext *context.AppCtx) *goji.Mux {
	mux := goji.SubMux()
	mux.HandleFunc(pat.Post("/team"), func(w http.ResponseWriter, r *http.Request) {
		CreateTeam(appContext, w, r)
	})

	mux.HandleFunc(pat.Post("/transaction"), func(w http.ResponseWriter, r *http.Request) {
		CreateTransaction(appContext, w, r)
	})

	return mux
}

//CreateTeam will build a team from JSON
func CreateTeam(appContext *context.AppCtx, w http.ResponseWriter, r *http.Request) {
	startNanos := time.Now().UnixNano()

	jsonDecoder := json.NewDecoder(r.Body)
	team := &models.CreateTeamBody{}
	err := jsonDecoder.Decode(team)

	if err != nil {
		appContext.Logger.Errorf("CreateTeam error: %s", err)
		appContext.Stats.Incr("api.create_team.400", 1)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	defer r.Body.Close()

	accessKey, err := appContext.Database.CreateTeam(team)
	if err != nil {
		appContext.Logger.Errorf("CreateTeam error: %s", err)
		appContext.Stats.Incr("api.create_team.500", 1)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	jsonBody, err := json.Marshal(models.AccessKeyResponse{accessKey})
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonBody)

	finishMilis := (time.Now().UnixNano() - startNanos) / 1000000
	appContext.Stats.Timing("api.create_team.response_time", finishMilis)
}

//CreateManager will build a manager from JSON
func CreateManager(appContext *context.AppCtx, w http.ResponseWriter, r *http.Request) {
	startNanos := time.Now().UnixNano()
	finishMilis := (time.Now().UnixNano() - startNanos) / 1000000
	appContext.Stats.Timing("api.create_manager.response_time", finishMilis)
}

//CreatePlayer will build a player from JSON
func CreatePlayer(appContext *context.AppCtx, w http.ResponseWriter, r *http.Request) {
	startNanos := time.Now().UnixNano()
	finishMilis := (time.Now().UnixNano() - startNanos) / 1000000
	appContext.Stats.Timing("api.create_player.response_time", finishMilis)

}

//CreateTransaction will create a transaction for a player, created by a manager
// identified with their access key
func CreateTransaction(appContext *context.AppCtx, w http.ResponseWriter, r *http.Request) {
	startNanos := time.Now().UnixNano()
	accessKey := r.Header.Get("X-Access-Key")

	if accessKey == "" {
		appContext.Logger.Errorf("CreateTransaction Error: Missing Access Key")
		appContext.Stats.Incr("api.create_team.400", 1)
		http.Error(w, "Missing X-Access-Key header", http.StatusBadRequest)
		return
	}

	manager, err := appContext.Database.GetManagerByAccessKey(accessKey)

	if err != nil {
		appContext.Logger.Errorf("CreateTransaction error: %s", err)
		appContext.Stats.Incr("api.create_transaction.404", 1)
		http.Error(w, "Can't find a manager with that access key", http.StatusNotFound)
		return
	}

	jsonDecoder := json.NewDecoder(r.Body)
	requestTransaction := &models.RequestTransaction{}
	err = jsonDecoder.Decode(requestTransaction)

	if err != nil {
		appContext.Logger.Errorf("CreateTransaction error: %s", err)
		appContext.Stats.Incr("api.create_transaction.400", 1)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	defer r.Body.Close()

	player, err := appContext.Database.GetPlayer(requestTransaction.Player)

	if err != nil {
		appContext.Logger.Errorf("CreateTransaction error: %s", err)
		appContext.Stats.Incr("api.create_transaction.404", 1)
		http.Error(w, "Can't find player under that ID", http.StatusNotFound)
		return
	}

	if player.Team == manager.TeamID {
		appContext.Logger.Errorf("CreateTransaction error: %s", err)
		appContext.Stats.Incr("api.create_transaction.422", 1)
		http.Error(w, "Player already assigned to manager", http.StatusUnprocessableEntity)
		return
	}

	transaction := &models.Transaction{
		FromTeam: player.Team,
		ToTeam: manager.TeamID,
		Player: player.ID,
	}

	err = appContext.Database.CreateTransaction(transaction)

	if err != nil {
		appContext.Logger.Errorf("CreateTransaction error: %s", err)
		appContext.Stats.Incr("api.create_transaction.500", 1)
		http.Error(w, "Error creating transaction", http.StatusInternalServerError)
		return
	}

	finishMilis := (time.Now().UnixNano() - startNanos) / 1000000
	appContext.Stats.Timing("api.create_player.response_time", finishMilis)
}

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
		appContext.Stats.Incr("api.create_team.500", 1)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}

	accessKey, err := appContext.Database.CreateTeam(team)
	if err != nil {
		appContext.Logger.Errorf("CreateTeam error: %s", err)
		appContext.Stats.Incr("api.create_team.500", 1)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
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

//CreateTransaction will create a transaction for
func CreateTransaction(appContext *context.AppCtx, w http.ResponseWriter, r *http.Request) {
	startNanos := time.Now().UnixNano()
	accessKey := r.Header.Get("X-Access-Key")

	if accessKey == "" {
		appContext.Logger.Errorf("CreateTransaction Error: Missing Access Key")
		appContext.Stats.Incr("api.create_team.400", 1)
		http.Error(w, "Missing X-Access-Key header", http.StatusBadRequest)
	}

	finishMilis := (time.Now().UnixNano() - startNanos) / 1000000
	appContext.Stats.Timing("api.create_player.response_time", finishMilis)
}

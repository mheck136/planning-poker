package gameapi

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/mheck136/planning-poker/gamecommands"
	"github.com/mheck136/planning-poker/gameregistry"
	"net/http"
)

func New(registry *gameregistry.GameRegistry) *GameApi {
	r := mux.NewRouter()
	gameApi := &GameApi{
		mux:          r,
		gameRegistry: registry,
	}

	r.Use(playerIdCookieMiddleware)
	r.Use(jsonContentTypeMiddleware)
	gameRouter := r.PathPrefix("/game/{gameId}").Subrouter()

	gameRouter.HandleFunc("/join", gameApi.gameHandler(gameApi.joinHandler)).Methods("POST")

	return gameApi
}

type GameApi struct {
	mux          *mux.Router
	gameRegistry *gameregistry.GameRegistry
}

func (g *GameApi) joinHandler(response http.ResponseWriter, request *http.Request, ctx gameContext) {
	type joinRequest struct {
		Name string `json:"name"`
	}
	var req joinRequest
	err := json.NewDecoder(request.Body).Decode(&req)
	if err != nil {
		response.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(response).Encode(map[string]string{"error": err.Error()})
		return
	}
	aggregate := g.gameRegistry.GetGameAggregateProxy(ctx.gameId)
	err = aggregate.SendJoinCommand(gamecommands.JoinCommand{
		PlayerId: ctx.playerId,
		Name:     req.Name,
	})
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(response).Encode(map[string]string{"error": err.Error()})
		return
	}
	_ = json.NewEncoder(response).Encode(map[string]string{"name": req.Name})
}

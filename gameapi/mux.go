package gameapi

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/mheck136/planning-poker/gamecommands"
	"github.com/mheck136/planning-poker/gameregistry"
	"github.com/rs/zerolog/log"
	"net/http"
)

func New(registry *gameregistry.GameRegistry, notificationService NotificationService) *GameApi {
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

type NotificationService interface {
	Register(gameId, playerId uuid.UUID, conn *websocket.Conn)
	SendJsonNotification(gameId uuid.UUID, message interface{})
}

type GameApi struct {
	mux                 *mux.Router
	gameRegistry        *gameregistry.GameRegistry
	notificationService NotificationService
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
	aggregate.SendJoinCommand(gamecommands.JoinCommand{
		PlayerId: ctx.playerId,
		Name:     req.Name,
	})
	log.Info().Str("playerId", ctx.playerId.String()).Str("playerId", ctx.playerId.String()).Msg("player joined game")
	g.notificationService.SendJsonNotification(ctx.gameId, map[string]string{"event": "PLAYER_JOINED", "playerId": ctx.playerId.String(), "name": req.Name})
	_ = json.NewEncoder(response).Encode(map[string]string{"name": req.Name})
}

func (a *GameApi) subscribeHandler(response http.ResponseWriter, request *http.Request, ctx gameContext) {
	conn, err := upgrader.Upgrade(response, request, nil)
	if err != nil {
		response.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(response).Encode(map[string]string{"error": err.Error()})
		return
	}
	a.notificationService.Register(ctx.gameId, ctx.playerId, conn)
	log.Info().Str("playerId", ctx.playerId.String()).Str("playerId", ctx.playerId.String()).Msg("new subscription started")
}

package api

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/mheck136/planning-poker/aggregate"
	"github.com/mheck136/planning-poker/notifier"
	"github.com/mheck136/planning-poker/registry"
	"github.com/rs/zerolog/log"
	"net/http"
)

func New(registry *registry.Registry, notifier *notifier.Notifier) *GameApi {
	r := mux.NewRouter()
	gameApi := &GameApi{
		mux:          r,
		gameRegistry: registry,
		notifier:     notifier,
	}

	r.Use(loggingMiddleware)
	r.Use(playerIdCookieMiddleware)
	r.Use(jsonContentTypeMiddleware)

	r.HandleFunc("/games/{gameId}", gameApi.createGameHandler).Methods("POST")
	r.HandleFunc("/games/{gameId}/actions/join", actionHandler(gameApi.joinActionHandler)).Methods("POST")
	r.HandleFunc("/games/{gameId}/actions/start", actionHandler(gameApi.startActionHandler)).Methods("POST")
	r.HandleFunc("/games/{gameId}/actions/vote", actionHandler(gameApi.voteActionHandler)).Methods("POST")
	r.HandleFunc("/games/{gameId}/actions/reveal", actionHandler(gameApi.revealActionHandler)).Methods("POST")
	r.HandleFunc("/games/{gameId}/actions/finish", actionHandler(gameApi.finishActionHandler)).Methods("POST")

	r.HandleFunc("/games/{gameId}/updates", gameApi.websocketHandler).Methods("GET")

	return gameApi
}

type NotificationService interface {
	Register(gameId, playerId uuid.UUID, conn *websocket.Conn)
	SendJsonNotification(gameId uuid.UUID, message interface{})
}

type GameApi struct {
	mux          *mux.Router
	gameRegistry *registry.Registry
	notifier     *notifier.Notifier
}

type CreateGameRequest struct {
	GameTitle string `json:"gameTitle"`
}

func (a *GameApi) createGameHandler(response http.ResponseWriter, request *http.Request) {
	gameId, err := extractGameId(request)
	if err != nil {
		response.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(response).Encode(map[string]string{"error": "missing or invalid game id"})
		return
	}
	req := &CreateGameRequest{}
	if ok := mustDecodeBody(response, request, req); ok {
		cmd := aggregate.NewCreateGameCommand(req.GameTitle)
		gameAggregate := a.gameRegistry.GetAggregateRoot(gameId)
		gameAggregate.HandleCommand(cmd)
		response.WriteHeader(http.StatusAccepted)
	}
}

func actionHandler(handler func(response http.ResponseWriter, request *http.Request, gameId, playerId uuid.UUID)) func(response http.ResponseWriter, request *http.Request) {
	return func(response http.ResponseWriter, request *http.Request) {
		gameId, err := extractGameId(request)
		if err != nil {
			response.WriteHeader(http.StatusBadRequest)
			_ = json.NewEncoder(response).Encode(map[string]string{"error": "missing or invalid game id"})
			return
		}
		playerId, err := extractPlayerId(request)
		if err != nil {
			response.WriteHeader(http.StatusBadRequest)
			_ = json.NewEncoder(response).Encode(map[string]string{"error": "missing or invalid player id"})
			return
		}
		handler(response, request, gameId, playerId)
	}
}

type JoinCommandRequest struct {
	Name string `json:"name"`
}

func (a *GameApi) joinActionHandler(response http.ResponseWriter, request *http.Request, gameId, playerId uuid.UUID) {
	req := &JoinCommandRequest{}
	if ok := mustDecodeBody(response, request, req); ok {
		cmd := aggregate.NewJoinCommand(playerId, req.Name)
		gameAggregate := a.gameRegistry.GetAggregateRoot(gameId)
		gameAggregate.HandleCommand(cmd)
		response.WriteHeader(http.StatusAccepted)
	}
}

type StartRoundCommandRequest struct {
	RoundName string `json:"roundName"`
}

func (a *GameApi) startActionHandler(response http.ResponseWriter, request *http.Request, gameId, _ uuid.UUID) {
	req := &StartRoundCommandRequest{}
	if ok := mustDecodeBody(response, request, req); ok {
		cmd := aggregate.NewStartRoundCommand(req.RoundName)
		gameAggregate := a.gameRegistry.GetAggregateRoot(gameId)
		gameAggregate.HandleCommand(cmd)
		response.WriteHeader(http.StatusAccepted)
	}
}

type CastVoteCommandRequest struct {
	Vote string `json:"vote"`
}

func (a *GameApi) voteActionHandler(response http.ResponseWriter, request *http.Request, gameId, playerId uuid.UUID) {
	req := &CastVoteCommandRequest{}
	if ok := mustDecodeBody(response, request, req); ok {
		cmd := aggregate.NewCastVoteCommand(playerId, req.Vote)
		gameAggregate := a.gameRegistry.GetAggregateRoot(gameId)
		gameAggregate.HandleCommand(cmd)
		response.WriteHeader(http.StatusAccepted)
	}
}

func (a *GameApi) revealActionHandler(response http.ResponseWriter, _ *http.Request, gameId, _ uuid.UUID) {
	cmd := aggregate.NewRevealCardsCommand()
	gameAggregate := a.gameRegistry.GetAggregateRoot(gameId)
	gameAggregate.HandleCommand(cmd)
	response.WriteHeader(http.StatusAccepted)
}

type FinishRoundCommandRequest struct {
	Result string `json:"result"`
}

func (a *GameApi) finishActionHandler(response http.ResponseWriter, request *http.Request, gameId, _ uuid.UUID) {
	req := &FinishRoundCommandRequest{}
	if ok := mustDecodeBody(response, request, req); ok {
		cmd := aggregate.NewFinishRoundCommand(req.Result)
		gameAggregate := a.gameRegistry.GetAggregateRoot(gameId)
		gameAggregate.HandleCommand(cmd)
		response.WriteHeader(http.StatusAccepted)
	}
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func (a *GameApi) websocketHandler(response http.ResponseWriter, request *http.Request) {
	log.Info().Msg("websocket request")
	gameId, err := extractGameId(request)
	if err != nil {
		response.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(response).Encode(map[string]string{"error": "missing or invalid game id"})
		return
	}
	playerId, err := extractPlayerId(request)
	if err != nil {
		response.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(response).Encode(map[string]string{"error": "missing or invalid player id"})
		return
	}
	log.Info().Msg("upgrading")
	conn, err := upgrader.Upgrade(response, request, nil)
	if err != nil {
		log.Error().Err(err).Msg("error while upgrading to websocket connection")
		return
	}
	a.notifier.HandleNewConnection(conn, playerId, gameId)
}

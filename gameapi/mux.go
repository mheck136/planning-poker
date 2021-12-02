package gameapi

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/mheck136/planning-poker/game"
	"github.com/mheck136/planning-poker/notifier"
	"github.com/rs/zerolog/log"
	"net/http"
)

func New(registry *game.Registry, notifier *notifier.Notifier) *GameApi {
	r := mux.NewRouter()
	gameApi := &GameApi{
		mux:          r,
		gameRegistry: registry,
		notifier:     notifier,
	}

	r.Use(loggingMiddleware)
	r.Use(playerIdCookieMiddleware)
	r.Use(jsonContentTypeMiddleware)

	r.HandleFunc("/games/{gameId}/{action}", gameApi.commandHandler).Methods("POST")
	r.HandleFunc("/game-updates/{gameId}", gameApi.websocketHandler).Methods("GET")

	return gameApi
}

type NotificationService interface {
	Register(gameId, playerId uuid.UUID, conn *websocket.Conn)
	SendJsonNotification(gameId uuid.UUID, message interface{})
}

type GameApi struct {
	mux          *mux.Router
	gameRegistry *game.Registry
	notifier     *notifier.Notifier
}

func (a *GameApi) commandHandler(response http.ResponseWriter, request *http.Request) {
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
	action, ok := mux.Vars(request)["action"]
	if !ok {
		response.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(response).Encode(map[string]string{"error": "missing or invalid action"})
		return
	}
	var cr CommandRequest
	var cmd game.Command
	switch action {
	case "join":
		cr = &JoinCommandRequest{}
	case "start":
		cr = &StartRoundCommandRequest{}
	case "vote":
		cr = &CastVoteCommandRequest{}
	case "reveal-cards":
		cmd = game.RevealCardsCommand{}
	case "finish-round":
		cr = &FinishRoundCommandRequest{}
	}
	if cr != nil {
		err = json.NewDecoder(request.Body).Decode(cr)
		if err != nil {
			response.WriteHeader(http.StatusBadRequest)
			_ = json.NewEncoder(response).Encode(map[string]string{"error": "invalid body request", "message": err.Error()})
			return
		}
		cmd = cr.toCommand(playerId)
	}
	gameAggregate := a.gameRegistry.GetAggregateRoot(gameId)
	gameAggregate.HandleCommand(cmd)
	response.WriteHeader(http.StatusAccepted)
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

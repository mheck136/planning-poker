package api

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"net/http"
)

func (a *GameApi) ListenAndServe(address string) error {
	return http.ListenAndServe(address, a.mux)
}

type gameContext struct {
	gameId   uuid.UUID
	playerId uuid.UUID
}

func (a *GameApi) gameHandler(handler func(response http.ResponseWriter, request *http.Request, ctx gameContext)) func(http.ResponseWriter, *http.Request) {
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
		handler(response, request, gameContext{
			gameId:   gameId,
			playerId: playerId,
		})
	}
}

func extractGameId(request *http.Request) (uuid.UUID, error) {
	vars := mux.Vars(request)
	gameId, ok := vars["gameId"]
	if !ok {
		return uuid.UUID{}, fmt.Errorf("gameId not set")
	}
	return uuid.Parse(gameId)
}

func extractPlayerId(request *http.Request) (uuid.UUID, error) {
	cookie, err := request.Cookie(playerIdCookieName)
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("player id cookie not present")
	}
	return uuid.Parse(cookie.Value)
}

func mustDecodeBody(response http.ResponseWriter, request *http.Request, target interface{}) bool {
	err := json.NewDecoder(request.Body).Decode(target)
	if err != nil {
		_ = json.NewEncoder(response).Encode(map[string]string{"error": "invalid body request", "message": err.Error()})
		return false
	}
	return true
}

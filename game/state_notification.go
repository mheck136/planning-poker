package game

import "github.com/google/uuid"

type State struct {
	GameId          uuid.UUID            `json:"gameId"`
	Players         map[uuid.UUID]string `json:"players"`
	ActiveRoundName string               `json:"activeRoundName,omitempty"`
	GameState       string               `json:"gameState"`
	Votes           map[uuid.UUID]string `json:"votes"`
}

package game

import "github.com/google/uuid"

type Snapshot struct {
	GameId          uuid.UUID            `json:"id"`
	Players         map[uuid.UUID]string `json:"players"`
	ActiveRoundName string               `json:"activeRoundName,omitempty"`
	GameState       string               `json:"gameState"`
	Votes           map[uuid.UUID]string `json:"votes"`
}

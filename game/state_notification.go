package game

import "github.com/google/uuid"

type State struct {
	GameId          uuid.UUID
	Players         map[uuid.UUID]string
	ActiveRoundName string
	GameState       string
	Votes           map[uuid.UUID]string
}

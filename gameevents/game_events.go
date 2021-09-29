package gameevents

import "github.com/google/uuid"

type GameEvent interface {
	isGameEvent()
}

type gameEventBase struct{}

func (e gameEventBase) isGameEvent() {}

type JoinedEvent struct {
	gameEventBase
	PlayerId uuid.UUID
	Name     string
}

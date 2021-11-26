package gamecommands

import "github.com/google/uuid"

type GameCommand interface {
	isGameCommand()
}

type gameCommandBase struct{}

func (e gameCommandBase) isGameCommand() {}

type JoinCommand struct {
	gameCommandBase
	PlayerId uuid.UUID
	Name     string
}

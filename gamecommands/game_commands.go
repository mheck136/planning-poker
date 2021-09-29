package gamecommands

import "github.com/google/uuid"

//JoinCommand
//no reply, no side effect
type JoinCommand struct {
	PlayerId uuid.UUID
	Name     string
}

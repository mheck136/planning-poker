package gamedomain

import (
	"github.com/google/uuid"
	"github.com/mheck136/planning-poker/gamecommands"
	"github.com/mheck136/planning-poker/gameevents"
)

func NewGameAggregate(id uuid.UUID) *GameAggregate {
	return &GameAggregate{
		table: NewTable(id),
	}
}

type GameAggregate struct {
	table *Table
}

func (a *GameAggregate) HandleJoinCommand(c gamecommands.JoinCommand) []gameevents.GameEvent {
	if a.table.HasPlayerJoined(c.PlayerId) {
		return nil
	}
	return []gameevents.GameEvent{
		gameevents.JoinedEvent{
			PlayerId: c.PlayerId,
			Name:     c.Name,
		},
	}
}

func (a *GameAggregate) HandleJoinedEvent(e gameevents.JoinedEvent) {
	a.table.Join(e.PlayerId, e.Name)
}

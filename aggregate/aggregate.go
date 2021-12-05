package aggregate

import (
	"github.com/google/uuid"
	"github.com/mheck136/planning-poker/game"
)

type Notifier interface {
	Notify(state game.Snapshot)
}

type Finisher interface {
	Finish(aggregate *Aggregate, err error)
}

type Command interface {
	CanExecute(aggregate *Aggregate) bool
	ExecuteCommand(aggregate *Aggregate) []Event
}

type Event interface {
	ExecuteEvent(aggregate *Aggregate)
	EventName() EventName
}

func New(id uuid.UUID, notifier Notifier) *Aggregate {
	return &Aggregate{
		board:    game.NewBoard(id),
		notifier: notifier,
	}
}

type Aggregate struct {
	board    *game.Board
	notifier Notifier
}

func (a *Aggregate) SendSnapshot() {
	a.notifier.Notify(a.board.ToSnapshot())
}

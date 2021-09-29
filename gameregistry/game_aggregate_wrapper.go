package gameregistry

import (
	"github.com/google/uuid"
	"github.com/mheck136/planning-poker/gamecommands"
	"github.com/mheck136/planning-poker/gamedomain"
	"github.com/mheck136/planning-poker/gameevents"
)

func NewGameAggregateProxy(id uuid.UUID) GameAggregateProxy {
	w := &GameAggregateWrapper{
		in:        make(chan interface{}),
		aggregate: gamedomain.NewGameAggregate(id),
	}
	go w.run()
	return w
}

type GameAggregateWrapper struct {
	in        chan interface{}
	aggregate *gamedomain.GameAggregate
}

func (w *GameAggregateWrapper) run() {
	for c := range w.in {
		switch command := c.(type) {
		case joinCommandWrapper:
			events := w.aggregate.HandleJoinCommand(command.joinCommand)
			// TODO: persist events
			err := w.handleEvents(events)
			command.reply <- err
			close(command.reply)
		}
	}
}

func (w *GameAggregateWrapper) handleEvents(events []gameevents.GameEvent) error {
	for _, e := range events {
		switch event := e.(type) {
		case gameevents.JoinedEvent:
			w.aggregate.HandleJoinedEvent(event)
		}
	}
	return nil
}

type joinCommandWrapper struct {
	joinCommand gamecommands.JoinCommand
	reply       chan error
}

func (w *GameAggregateWrapper) SendJoinCommand(c gamecommands.JoinCommand) error {
	reply := make(chan error, 1)
	w.in <- joinCommandWrapper{
		joinCommand: c,
		reply:       reply,
	}
	return <-reply
}

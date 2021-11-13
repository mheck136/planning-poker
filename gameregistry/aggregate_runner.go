package gameregistry

import (
	"github.com/google/uuid"
	"github.com/mheck136/planning-poker/gamecommands"
	"github.com/mheck136/planning-poker/gamedomain"
	"github.com/mheck136/planning-poker/gameevents"
)

func NewAggregateRunner(id uuid.UUID) GameAggregateProxy {
	w := &AggregateRunner{
		in:        make(chan interface{}),
		aggregate: gamedomain.NewGameAggregate(id),
	}
	go w.run()
	return w
}

type AggregateRunner struct {
	in        chan interface{}
	aggregate *gamedomain.GameAggregate
}

func (r *AggregateRunner) run() {
	for c := range r.in {
		switch command := c.(type) {
		case joinCommandWrapper:
			events := r.aggregate.HandleJoinCommand(command.joinCommand)
			err := r.handleEvents(events)
			command.reply <- err
			close(command.reply)
		}
	}
}

func (r *AggregateRunner) handleEvents(events []gameevents.GameEvent) error {
	// TODO: persist events
	for _, e := range events {
		switch event := e.(type) {
		case gameevents.JoinedEvent:
			r.aggregate.HandleJoinedEvent(event)
		}
	}
	return nil
}

type joinCommandWrapper struct {
	joinCommand gamecommands.JoinCommand
	reply       chan error
}

func (r *AggregateRunner) SendJoinCommand(c gamecommands.JoinCommand) error {
	reply := make(chan error, 1)
	r.in <- joinCommandWrapper{
		joinCommand: c,
		reply:       reply,
	}
	return <-reply
}

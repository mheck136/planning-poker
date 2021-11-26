package gameregistry

import (
	"github.com/google/uuid"
	"github.com/mheck136/planning-poker/gamecommands"
	"github.com/mheck136/planning-poker/gamedomain"
	"github.com/mheck136/planning-poker/gameevents"
)

type Aggregate interface {
	HandleJoinCommand(command gamecommands.JoinCommand) []gameevents.GameEvent
	HandleJoinedEvent(event gameevents.JoinedEvent)
	GetTableSnapshot() gamedomain.Table
}

type StatePublisher interface {
	PublishState(table gamedomain.Table)
}

type AggregateRunnerCfg struct {
	AggregateFactory func(uuid.UUID) Aggregate
	Publisher        StatePublisher
}

func NewAggregateRunnerFactory(cfg AggregateRunnerCfg) func(uuid.UUID) GameAggregateProxy {
	config := cfg
	return func(id uuid.UUID) GameAggregateProxy {
		w := &AggregateRunner{
			in:        make(chan gamecommands.GameCommand),
			aggregate: config.AggregateFactory(id),
			publisher: config.Publisher,
		}
		go w.run()
		return w
	}

}

type AggregateRunner struct {
	in        chan gamecommands.GameCommand
	aggregate Aggregate
	publisher StatePublisher
}

func (r *AggregateRunner) run() {
	for c := range r.in {
		var err error
		switch command := c.(type) {
		case gamecommands.JoinCommand:
			events := r.aggregate.HandleJoinCommand(command)
			err = r.handleEvents(events)
		}
		if err != nil {
			// TODO handle error
		} else {
			snapshot := r.aggregate.GetTableSnapshot()
			r.publisher.PublishState(snapshot)
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

func (r *AggregateRunner) SendJoinCommand(c gamecommands.JoinCommand) {
	r.in <- c
}

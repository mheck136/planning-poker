package registry

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/mheck136/planning-poker/aggregate"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func newRoot(id uuid.UUID, eventStore EventStore, aggregateFactory func(id uuid.UUID) *aggregate.Aggregate) *Root {
	logger := log.With().Str("component", "aggregate.Root").Str("id", id.String()).Logger()
	root := &Root{
		id:         id,
		version:    0,
		commands:   make(chan aggregate.Command),
		eventStore: eventStore,
		log:        logger,
		aggregate:  aggregateFactory(id),
	}
	go root.run()
	return root
}

type Root struct {
	id         uuid.UUID
	version    int
	commands   chan aggregate.Command
	eventStore EventStore
	log        zerolog.Logger
	aggregate  *aggregate.Aggregate
}

func (r *Root) HandleCommand(c aggregate.Command) {
	r.commands <- c
}

func (r *Root) run() {
	for command := range r.commands {
		var err error
		if command.CanExecute(r.aggregate) {
			events := command.ExecuteCommand(r.aggregate)
			if len(events) > 0 {
				err = r.persistEvents(events)
				if err != nil {
					log.Error().Err(err).Msg("error while persisting events")
				} else {
					r.handleEvents(events)
					r.version += len(events)
				}
			}
		}
		if finisher, ok := command.(aggregate.Finisher); ok {
			finisher.Finish(r.aggregate, err)
		}
	}
}

func (r *Root) persistEvents(events []aggregate.Event) error {
	pEvents := make([]PersistedEvent, len(events))
	for i, event := range events {
		eventName := event.EventName()
		marshalledData, err := json.Marshal(event)
		if err != nil {
			return err
		}
		pEvents[i] = PersistedEvent{
			AggregateId: r.id,
			SequenceNo:  r.version + i,
			EventName:   eventName,
			Data:        marshalledData,
		}
	}
	return r.eventStore.StoreEvents(pEvents)
}

func (r *Root) handleEvents(events []aggregate.Event) {
	for _, event := range events {
		event.ExecuteEvent(r.aggregate)
	}
}

type EventStore interface {
	StoreEvents([]PersistedEvent) error
	RetrieveAllEvents(aggregateId uuid.UUID) ([]PersistedEvent, error)
}

type PersistedEvent struct {
	AggregateId uuid.UUID
	SequenceNo  int
	EventName   aggregate.EventName
	Data        []byte
}

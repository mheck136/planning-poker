package game

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type PersistedEvent struct {
	AggregateId uuid.UUID
	SequenceNo  int
	EventName   string
	Data        []byte
}

type EventStore interface {
	StoreEvents([]PersistedEvent) error
	RetrieveAllEvents(aggregateId uuid.UUID) ([]PersistedEvent, error)
}

type Command interface {
	ExecuteCommand(aggregateRoot *AggregateRoot) []Event
}

type Event interface {
	ExecuteEvent(aggregateRoot *AggregateRoot)
}

func newAggregateRoot(id uuid.UUID, eventStore EventStore) *AggregateRoot {
	logger := log.With().Str("component", "game.AggregateRoot").Str("id", id.String()).Logger()
	root := &AggregateRoot{
		id:         id,
		version:    0,
		commands:   make(chan Command),
		eventStore: eventStore,
		log:        logger,
	}
	go root.run()
	return root
}

type AggregateRoot struct {
	id         uuid.UUID
	version    int
	commands   chan Command
	eventStore EventStore
	log        zerolog.Logger
}

func (a *AggregateRoot) HandleCommand(c Command) {
	a.commands <- c
}

func (a *AggregateRoot) run() {
	for command := range a.commands {
		events := command.ExecuteCommand(a)
		if len(events) > 0 {
			err := a.persistEvents(events)
			if err != nil {
				log.Error().Err(err).Msg("error while persisting events")
			} else {
				a.handleEvents(events)
				a.version += len(events)
				a.publishState()
			}
		}
	}
}

func (a *AggregateRoot) persistEvents(events []Event) error {
	pEvents := make([]PersistedEvent, len(events))
	for i, event := range events {
		eventName := ""
		switch event.(type) {
		case PlayerJoinedEvent:
			eventName = playerJoinedEventName
		}
		marshalledData, err := json.Marshal(event)
		if err != nil {
			return err
		}
		pEvents[i] = PersistedEvent{
			AggregateId: a.id,
			SequenceNo:  a.version + i,
			EventName:   eventName,
			Data:        marshalledData,
		}
	}
	return a.eventStore.StoreEvents(pEvents)
}

func (a *AggregateRoot) handleEvents(events []Event) {
	for _, event := range events {
		event.ExecuteEvent(a)
	}
}

func (a *AggregateRoot) publishState() {

}

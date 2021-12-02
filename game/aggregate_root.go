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
	EventName   EventName
	Data        []byte
}

type EventStore interface {
	StoreEvents([]PersistedEvent) error
	RetrieveAllEvents(aggregateId uuid.UUID) ([]PersistedEvent, error)
}

type Notifier interface {
	Notify(state State)
}

type Command interface {
	CanExecute(aggregateRoot *AggregateRoot) bool
	ExecuteCommand(aggregateRoot *AggregateRoot) []Event
}

type Event interface {
	ExecuteEvent(aggregateRoot *AggregateRoot)
	EventName() EventName
}

func newAggregateRoot(id uuid.UUID, eventStore EventStore, notifier Notifier) *AggregateRoot {
	logger := log.With().Str("component", "game.AggregateRoot").Str("id", id.String()).Logger()
	root := &AggregateRoot{
		id:         id,
		version:    0,
		commands:   make(chan Command),
		eventStore: eventStore,
		notifier:   notifier,
		log:        logger,
		board:      newBoard(),
	}
	go root.run()
	return root
}

type AggregateRoot struct {
	id         uuid.UUID
	version    int
	commands   chan Command
	eventStore EventStore
	notifier   Notifier
	log        zerolog.Logger
	board      *board
}

func (a *AggregateRoot) HandleCommand(c Command) {
	a.commands <- c
}

func (a *AggregateRoot) run() {
	for command := range a.commands {
		if command.CanExecute(a) {
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
}

func (a *AggregateRoot) persistEvents(events []Event) error {
	pEvents := make([]PersistedEvent, len(events))
	for i, event := range events {
		eventName := event.EventName()
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
	a.notifier.Notify(a.toState())
}

func (a *AggregateRoot) toState() State {
	state := State{
		Players:         make(map[uuid.UUID]string),
		GameId:          a.id,
		GameState:       a.board.state.String(),
		ActiveRoundName: a.board.activeRoundName,
	}
	for _, p := range a.board.players {
		state.Players[p.id] = p.name
	}
	votes := make(map[uuid.UUID]string, len(a.board.players))
	if a.board.isOpenForVotes() || a.board.isDeciding() {
		for playerId, vote := range a.board.votes {
			if a.board.isOpenForVotes() {
				votes[playerId] = ""
			} else {
				votes[playerId] = vote
			}
		}
	}
	state.Votes = votes
	return state
}

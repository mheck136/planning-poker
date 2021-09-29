package game

import (
	"github.com/google/uuid"
	"log"
)

type Aggregate struct {
	id         uuid.UUID
	version    int
	game       Game
	cmdIn      chan Command
	repository EventRepository
}

func (a *Aggregate) run() {
	a.rehydrate()
	for nextCmd := range a.cmdIn {
		log.Printf("processing command")
		events, err := nextCmd.execute(a.game)
		if err != nil {
			log.Printf("error while executing command: %v", err)
			continue
		}
		err = a.storeEvents(events, a.version)
		if err != nil {
			log.Printf("error while storing events: %v", err)
		}
		for _, event := range events {
			a.game = event.apply(a.game)
		}
		a.version += len(events)
		if replier, ok := nextCmd.(Replier); ok {
			replier.reply(a.game, events)
		}
	}
}

func (a *Aggregate) rehydrate() {
	events, err := a.repository.RetrieveEvents(a.id)
	if err != nil {
		log.Printf("error retrieving events: %v", err)
	} else {
		for _, event := range events {
			a.game = event.apply(a.game)
		}
	}
}

func (a *Aggregate) storeEvents(events []Event, version int) error {
	_ = a.repository.StoreEvents(a.id, version, events)
	return nil
}

func (a *Aggregate) SubmitCommand(c Command) {
	a.cmdIn <- c
}

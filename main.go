package main

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/mheck136/planning-poker/game"
	"log"
)

type InMemoryRepository struct {

}

func (r *InMemoryRepository) StoreEvents(id uuid.UUID, version int, events []game.Event) error {
	for i, event := range events {
		fmt.Printf("storing event for id %v with version %v: %v\n", id, version + i, event)
	}
	return nil
}

func (r *InMemoryRepository) RetrieveEvents(id uuid.UUID) ([]game.Event, error) {
	fmt.Printf("retrieving events for id %v\n", id)
	return nil, nil
}

func main() {
	rep := InMemoryRepository{}
	reg := game.NewRegistry(&rep)
	id := uuid.New()
	agg := reg.GetAggregate(id)
	cmd := game.NewCreateCommand("My Test Game")
	agg.SubmitCommand(cmd)
	log.Printf("new game: %v\n", cmd.Result())

}

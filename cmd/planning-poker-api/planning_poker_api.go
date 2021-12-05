package main

import (
	"github.com/google/uuid"
	"github.com/mheck136/planning-poker/aggregate"
	"github.com/mheck136/planning-poker/api"
	"github.com/mheck136/planning-poker/notifier"
	"github.com/mheck136/planning-poker/registry"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
)

type mockEventStore int

func (m mockEventStore) StoreEvents(events []registry.PersistedEvent) error {
	for _, event := range events {
		log.Info().Interface("events", events).RawJSON("event", event.Data).Msg("storing event")
	}
	return nil
}

func (m mockEventStore) RetrieveAllEvents(uuid.UUID) ([]registry.PersistedEvent, error) {
	panic("implement me")
}

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	not := notifier.New()

	reg := registry.New(mockEventStore(1), func(id uuid.UUID) *aggregate.Aggregate {
		return aggregate.New(id, not)
	})

	a := api.New(
		reg,
		not,
	)
	log.Fatal().Err(a.ListenAndServe("127.0.0.1:8080")).Msg("")

}

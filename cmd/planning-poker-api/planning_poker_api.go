package main

import (
	"github.com/google/uuid"
	"github.com/mheck136/planning-poker/game"
	"github.com/mheck136/planning-poker/gameapi"
	"github.com/mheck136/planning-poker/notifier"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
)

type mockEventStore int

func (m mockEventStore) StoreEvents(events []game.PersistedEvent) error {
	for _, event := range events {
		log.Info().Interface("events", events).RawJSON("event", event.Data).Msg("storing event")
	}
	return nil
}

func (m mockEventStore) RetrieveAllEvents(uuid.UUID) ([]game.PersistedEvent, error) {
	panic("implement me")
}

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	not := notifier.New()

	reg := game.NewRegistry(mockEventStore(1), not)

	a := gameapi.New(
		reg,
		not,
	)
	log.Fatal().Err(a.ListenAndServe("127.0.0.1:8080")).Msg("")

}

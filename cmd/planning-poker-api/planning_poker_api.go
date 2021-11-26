package main

import (
	"github.com/google/uuid"
	"github.com/mheck136/planning-poker/gameapi"
	"github.com/mheck136/planning-poker/gamedomain"
	"github.com/mheck136/planning-poker/gameregistry"
	"github.com/mheck136/planning-poker/notifications"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
)

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	notifier := notifications.New()
	a := gameapi.New(
		gameregistry.NewGameRegistry(gameregistry.NewAggregateRunnerFactory(gameregistry.AggregateRunnerCfg{
			AggregateFactory: func(uuid uuid.UUID) gameregistry.Aggregate {
				return gamedomain.NewGameAggregate(uuid)
			},
			Publisher: notifier,
		})),
		notifier,
	)
	log.Fatal().Err(a.ListenAndServe("127.0.0.1:8080")).Msg("")

}

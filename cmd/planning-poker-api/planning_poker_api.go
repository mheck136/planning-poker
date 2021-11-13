package main

import (
	"github.com/mheck136/planning-poker/gameapi"
	"github.com/mheck136/planning-poker/gameregistry"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
)

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	a := gameapi.New(gameregistry.NewGameRegistry(gameregistry.NewAggregateRunner))
	log.Fatal().Err(a.ListenAndServe("127.0.0.1:8080")).Msg("")

}

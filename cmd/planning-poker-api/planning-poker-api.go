package main

import (
	"github.com/mheck136/planning-poker/gameapi"
	"github.com/mheck136/planning-poker/gameregistry"
	"github.com/rs/zerolog/log"
)

func main() {
	a := gameapi.New(gameregistry.NewGameRegistry(gameregistry.NewGameAggregateProxy))
	log.Fatal().Err(a.ListenAndServe("127.0.0.1:8080")).Msg("")

}

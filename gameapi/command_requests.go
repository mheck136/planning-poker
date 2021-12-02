package gameapi

import (
	"github.com/google/uuid"
	"github.com/mheck136/planning-poker/game"
)

type CommandRequest interface {
	toCommand(playerId uuid.UUID) game.Command
}

type JoinCommandRequest struct {
	Name string `json:"name"`
}

func (r JoinCommandRequest) toCommand(playerId uuid.UUID) game.Command {
	return game.JoinCommand{
		PlayerId:   playerId,
		PlayerName: r.Name,
	}
}

type StartRoundCommandRequest struct {
	RoundName string `json:"roundName"`
}

func (s StartRoundCommandRequest) toCommand(uuid.UUID) game.Command {
	return game.StartRoundCommand{RoundName: s.RoundName}
}

type CastVoteCommandRequest struct {
	Vote string `json:"vote"`
}

func (r CastVoteCommandRequest) toCommand(playerId uuid.UUID) game.Command {
	return game.CastVoteCommand{
		PlayerId: playerId,
		Vote:     r.Vote,
	}
}

type FinishRoundCommandRequest struct {
	Result string `json:"result"`
}

func (r FinishRoundCommandRequest) toCommand(uuid.UUID) game.Command {
	return game.FinishRoundCommand{
		Result: r.Result,
	}
}

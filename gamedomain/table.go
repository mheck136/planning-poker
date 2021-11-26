package gamedomain

import (
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

type tableState string

const (
	idle         tableState = "IDLE"
	openForVotes tableState = "OPEN_FOR_VOTES"
	deciding     tableState = "DECIDING"
)

type tableError string

func (e tableError) Error() string {
	return string(e)
}

const (
	gameNotIdleError         tableError = "game not idle"
	gameNotOpenForVotesError tableError = "game not open for votes"
)

func NewTable(id uuid.UUID) *Table {
	return &Table{
		id:        id,
		players:   make(map[uuid.UUID]string),
		state:     idle,
		roundName: "",
		votes:     nil,
	}
}

type Table struct {
	id        uuid.UUID
	players   map[uuid.UUID]string
	state     tableState
	roundName string
	votes     map[uuid.UUID]string
}

func (t *Table) Join(playerId uuid.UUID, name string) {
	if currentName, ok := t.players[playerId]; !ok || currentName != name {
		t.players[playerId] = name
		log.Info().Str("playerId", playerId.String()).Str("playerName", name).Msg("player joined")
	}
}

func (t *Table) Leave(playerId uuid.UUID) bool {
	_, ok := t.players[playerId]
	delete(t.players, playerId)
	return ok
}

func (t *Table) StartRound(roundName string) error {
	if t.state != idle {
		return gameNotIdleError
	}
	t.clearVotes()
	t.state = openForVotes
	t.roundName = roundName
	return nil
}

func (t *Table) clearVotes() {
	t.votes = make(map[uuid.UUID]string)
}

func (t *Table) Vote(playerId uuid.UUID, vote string) (bool, error) {
	if t.state != openForVotes {
		return false, gameNotOpenForVotesError
	}
	t.votes[playerId] = vote
	// check whether all players have voted
	for userId := range t.players {
		_, ok := t.votes[userId]
		if !ok {
			return false, nil
		}
	}
	t.state = deciding
	return true, nil
}

func (t *Table) ResetRound() error {
	if t.state != openForVotes {
		return gameNotOpenForVotesError
	}
	t.clearVotes()
	return nil
}

func (t *Table) CloseRound() error {
	if t.state != openForVotes {
		return gameNotOpenForVotesError
	}
	t.state = deciding
	t.roundName = ""
	return nil
}

func (t *Table) HasPlayerJoined(playerId uuid.UUID) bool {
	for id := range t.players {
		if id == playerId {
			return true
		}
	}
	return false
}

func (t *Table) Copy() Table {
	players := map[uuid.UUID]string{}
	for playerId, playerName := range t.players {
		players[playerId] = playerName
	}
	votes := map[uuid.UUID]string{}
	for playerId, vote := range t.votes {
		votes[playerId] = vote
	}
	return Table{
		id:        t.id,
		players:   players,
		state:     t.state,
		roundName: t.roundName,
		votes:     votes,
	}
}

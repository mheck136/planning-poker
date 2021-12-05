package game

import (
	"github.com/google/uuid"
)

type player struct {
	id   uuid.UUID
	name string
}

type boardState int

const (
	_ boardState = iota
	idle
	openForVotes
	deciding
)

func (s boardState) String() string {
	switch s {
	case idle:
		return "IDLE"
	case openForVotes:
		return "OPEN_FOR_VOTES"
	case deciding:
		return "DECIDING"
	default:
		return "UNDEFINED"
	}
}

func NewBoard(id uuid.UUID) *Board {
	return &Board{
		id:              id,
		players:         nil,
		state:           idle,
		activeRoundName: "",
	}
}

type Board struct {
	id              uuid.UUID
	title           string
	initialized     bool
	players         []player
	state           boardState
	activeRoundName string
	votes           map[uuid.UUID]string
}

func (b *Board) IsInitialized() bool {
	return b.initialized
}

func (b *Board) Initialize(title string) {
	b.title = title
	b.initialized = true
}

func (b *Board) KnowsPlayer(id uuid.UUID) bool {
	for _, p := range b.players {
		if p.id == id {
			return true
		}
	}
	return false
}

func (b *Board) AddPlayer(playerId uuid.UUID, name string) {
	if !b.KnowsPlayer(playerId) {
		b.players = append(b.players, player{
			id:   playerId,
			name: name,
		})
	}
}

func (b *Board) IsIdle() bool {
	return b.state == idle
}

func (b *Board) IsOpenForVotes() bool {
	return b.state == openForVotes
}

func (b *Board) IsDeciding() bool {
	return b.state == deciding
}

func (b *Board) StartRound(roundName string) {
	if b.state == idle {
		b.activeRoundName = roundName
		b.state = openForVotes
		b.votes = make(map[uuid.UUID]string)
	}
}

func (b *Board) CastVote(playerId uuid.UUID, vote string) {
	if b.state == openForVotes && b.KnowsPlayer(playerId) {
		b.votes[playerId] = vote
	}
}

func (b *Board) allVoted() (allVoted bool) {
	allVoted = true
	for _, p := range b.players {
		_, ok := b.votes[p.id]
		allVoted = allVoted && ok
		if !allVoted {
			return
		}
	}
	return
}

func (b *Board) RevealCards() {
	if b.state == openForVotes {
		b.state = deciding
	}
}

func (b *Board) FinishRound() {
	if b.state == deciding {
		b.votes = nil
		b.state = idle
		b.activeRoundName = ""
	}
}

func (b *Board) ToSnapshot() Snapshot {
	snapshot := Snapshot{
		GameId:          b.id,
		Players:         make(map[uuid.UUID]string),
		GameState:       b.state.String(),
		ActiveRoundName: b.activeRoundName,
	}
	for _, p := range b.players {
		snapshot.Players[p.id] = p.name
	}
	votes := make(map[uuid.UUID]string, len(b.players))
	if b.IsOpenForVotes() || b.IsDeciding() {
		for playerId, vote := range b.votes {
			if b.IsOpenForVotes() {
				votes[playerId] = ""
			} else {
				votes[playerId] = vote
			}
		}
	}
	snapshot.Votes = votes
	return snapshot
}

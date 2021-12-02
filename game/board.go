package game

import "github.com/google/uuid"

type player struct {
	id   uuid.UUID
	name string
}

type boardState int

const (
	undefinedState = iota
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

func newBoard() *board {
	return &board{
		players:         nil,
		state:           idle,
		activeRoundName: "",
	}
}

type board struct {
	players         []player
	state           boardState
	activeRoundName string
	votes           map[uuid.UUID]string
}

func (b *board) knowsPlayer(id uuid.UUID) bool {
	for _, p := range b.players {
		if p.id == id {
			return true
		}
	}
	return false
}

func (b *board) addPlayer(p player) {
	if !b.knowsPlayer(p.id) {
		b.players = append(b.players, p)
	}
}

func (b *board) isIdle() bool {
	return b.state == idle
}

func (b *board) isOpenForVotes() bool {
	return b.state == openForVotes
}

func (b *board) isDeciding() bool {
	return b.state == deciding
}

func (b *board) startRound(roundName string) {
	if b.state == idle {
		b.activeRoundName = roundName
		b.state = openForVotes
		b.votes = make(map[uuid.UUID]string)
	}
}

func (b *board) castVote(playerId uuid.UUID, vote string) {
	if b.state == openForVotes && b.knowsPlayer(playerId) {
		b.votes[playerId] = vote

		allVoted := true
		for _, p := range b.players {
			_, ok := b.votes[p.id]
			allVoted = allVoted && ok
		}
		if allVoted {
			b.state = deciding
		}

	}
}

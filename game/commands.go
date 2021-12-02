package game

import "github.com/google/uuid"

type JoinCommand struct {
	PlayerId   uuid.UUID
	PlayerName string
}

func (c JoinCommand) CanExecute(r *AggregateRoot) bool {
	return !r.board.knowsPlayer(c.PlayerId)
}

func (c JoinCommand) ExecuteCommand(*AggregateRoot) []Event {
	return []Event{
		PlayerJoinedEvent{
			PlayerId:   c.PlayerId,
			PlayerName: c.PlayerName,
		},
	}
}

type StartRoundCommand struct {
	RoundName string
}

func (c StartRoundCommand) CanExecute(r *AggregateRoot) bool {
	return r.board.isIdle()
}

func (c StartRoundCommand) ExecuteCommand(*AggregateRoot) []Event {
	return []Event{
		RoundStartedEvent{RoundName: c.RoundName},
	}
}

type CastVoteCommand struct {
	PlayerId uuid.UUID
	Vote     string
}

func (c CastVoteCommand) CanExecute(r *AggregateRoot) bool {
	return r.board.isOpenForVotes()
}

func (c CastVoteCommand) ExecuteCommand(*AggregateRoot) []Event {
	return []Event{
		VoteCastEvent{
			PlayerId: c.PlayerId,
			Vote:     c.Vote,
		},
	}
}

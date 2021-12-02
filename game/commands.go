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

func (StartRoundCommand) CanExecute(r *AggregateRoot) bool {
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

func (CastVoteCommand) CanExecute(r *AggregateRoot) bool {
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

type FinishRoundCommand struct {
	Result string
}

func (FinishRoundCommand) CanExecute(r *AggregateRoot) bool {
	return r.board.isDeciding()
}

func (c FinishRoundCommand) ExecuteCommand(*AggregateRoot) []Event {
	return []Event{
		RoundFinishedEvent{Result: c.Result},
	}
}

type RevealCardsCommand struct {
}

func (RevealCardsCommand) CanExecute(r *AggregateRoot) bool {
	return r.board.isOpenForVotes()
}

func (c RevealCardsCommand) ExecuteCommand(*AggregateRoot) []Event {
	return []Event{
		CardsRevealedEvent{},
	}
}

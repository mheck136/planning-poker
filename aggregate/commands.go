package aggregate

import "github.com/google/uuid"

func NewCreateGameCommand(gameTitle string) CreateGameCommand {
	return CreateGameCommand{
		gameTitle: gameTitle,
	}
}

type CreateGameCommand struct {
	gameTitle string
}

func (c CreateGameCommand) CanExecute(a *Aggregate) bool {
	return !a.board.IsInitialized()
}

func (c CreateGameCommand) ExecuteCommand(*Aggregate) []Event {
	return []Event{
		GameCreatedEvent{Title: c.gameTitle},
	}
}

func NewJoinCommand(playerId uuid.UUID, playerName string) JoinCommand {
	return JoinCommand{
		playerId:   playerId,
		playerName: playerName,
	}
}

type JoinCommand struct {
	playerId   uuid.UUID
	playerName string
}

func (c JoinCommand) CanExecute(r *Aggregate) bool {
	return !r.board.KnowsPlayer(c.playerId)
}

func (c JoinCommand) ExecuteCommand(*Aggregate) []Event {
	return []Event{
		PlayerJoinedEvent{
			PlayerId:   c.playerId,
			PlayerName: c.playerName,
		},
	}
}

func (JoinCommand) Finish(aggregate *Aggregate, err error) {
	if err == nil {
		aggregate.SendSnapshot()
	}
}

func NewStartRoundCommand(roundName string) StartRoundCommand {
	return StartRoundCommand{
		roundName: roundName,
	}
}

type StartRoundCommand struct {
	roundName string
}

func (StartRoundCommand) CanExecute(r *Aggregate) bool {
	return r.board.IsIdle()
}

func (c StartRoundCommand) ExecuteCommand(*Aggregate) []Event {
	return []Event{
		RoundStartedEvent{RoundName: c.roundName},
	}
}

func (StartRoundCommand) Finish(aggregate *Aggregate, err error) {
	if err == nil {
		aggregate.SendSnapshot()
	}
}

func NewCastVoteCommand(playerId uuid.UUID, vote string) CastVoteCommand {
	return CastVoteCommand{
		playerId: playerId,
		vote:     vote,
	}
}

type CastVoteCommand struct {
	playerId uuid.UUID
	vote     string
}

func (CastVoteCommand) CanExecute(r *Aggregate) bool {
	return r.board.IsOpenForVotes()
}

func (c CastVoteCommand) ExecuteCommand(*Aggregate) []Event {
	return []Event{
		VoteCastEvent{
			PlayerId: c.playerId,
			Vote:     c.vote,
		},
	}
}

func (CastVoteCommand) Finish(aggregate *Aggregate, err error) {
	if err == nil {
		aggregate.SendSnapshot()
	}
}

func NewFinishRoundCommand(result string) FinishRoundCommand {
	return FinishRoundCommand{
		result: result,
	}
}

type FinishRoundCommand struct {
	result string
}

func (FinishRoundCommand) CanExecute(r *Aggregate) bool {
	return r.board.IsDeciding()
}

func (c FinishRoundCommand) ExecuteCommand(*Aggregate) []Event {
	return []Event{
		RoundFinishedEvent{Result: c.result},
	}
}

func (FinishRoundCommand) Finish(aggregate *Aggregate, err error) {
	if err == nil {
		aggregate.SendSnapshot()
	}
}

func NewRevealCardsCommand() RevealCardsCommand {
	return RevealCardsCommand{}
}

type RevealCardsCommand struct {
}

func (RevealCardsCommand) CanExecute(r *Aggregate) bool {
	return r.board.IsOpenForVotes()
}

func (c RevealCardsCommand) ExecuteCommand(*Aggregate) []Event {
	return []Event{
		CardsRevealedEvent{},
	}
}

func (RevealCardsCommand) Finish(aggregate *Aggregate, err error) {
	if err == nil {
		aggregate.SendSnapshot()
	}
}

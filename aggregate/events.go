package aggregate

import "github.com/google/uuid"

type EventName string

const (
	_                      EventName = ""
	gameCreatedEventName   EventName = "GAME_CREATED"
	playerJoinedEventName  EventName = "PLAYER_JOINED"
	roundStartedEventName  EventName = "ROUND_STARTED"
	voteCastEventName      EventName = "VOTE_CAST"
	finishedRoundEventName EventName = "FINISHED_ROUND"
	cardsRevealedEventName EventName = "CARDS_REVEALED"
)

type GameCreatedEvent struct {
	Title string
}

func (e GameCreatedEvent) ExecuteEvent(a *Aggregate) {
	a.board.Initialize(e.Title)
}

func (e GameCreatedEvent) EventName() EventName {
	return gameCreatedEventName
}

type PlayerJoinedEvent struct {
	PlayerId   uuid.UUID
	PlayerName string
}

func (e PlayerJoinedEvent) EventName() EventName {
	return playerJoinedEventName
}

func (e PlayerJoinedEvent) ExecuteEvent(a *Aggregate) {
	a.board.AddPlayer(e.PlayerId, e.PlayerName)
}

type RoundStartedEvent struct {
	RoundName string
}

func (e RoundStartedEvent) EventName() EventName {
	return roundStartedEventName
}

func (e RoundStartedEvent) ExecuteEvent(a *Aggregate) {
	a.board.StartRound(e.RoundName)
}

type VoteCastEvent struct {
	PlayerId uuid.UUID
	Vote     string
}

func (e VoteCastEvent) ExecuteEvent(a *Aggregate) {
	a.board.CastVote(e.PlayerId, e.Vote)
}

func (e VoteCastEvent) EventName() EventName {
	return voteCastEventName
}

type RoundFinishedEvent struct {
	Result string
}

func (e RoundFinishedEvent) ExecuteEvent(a *Aggregate) {
	a.board.FinishRound()
}

func (e RoundFinishedEvent) EventName() EventName {
	return finishedRoundEventName
}

type CardsRevealedEvent struct{}

func (e CardsRevealedEvent) ExecuteEvent(a *Aggregate) {
	a.board.RevealCards()
}

func (e CardsRevealedEvent) EventName() EventName {
	return cardsRevealedEventName
}

package game

import "github.com/google/uuid"

type EventName string

const (
	undefinedEventName    EventName = ""
	playerJoinedEventName EventName = "PLAYER_JOINED"
	roundStartedEventName EventName = "ROUND_STARTED"
	voteCastEventName     EventName = "VOTE_CAST"
)

type PlayerJoinedEvent struct {
	PlayerId   uuid.UUID
	PlayerName string
}

func (e PlayerJoinedEvent) EventName() EventName {
	return playerJoinedEventName
}

func (e PlayerJoinedEvent) ExecuteEvent(a *AggregateRoot) {
	a.board.addPlayer(player{
		id:   e.PlayerId,
		name: e.PlayerName,
	})
}

type RoundStartedEvent struct {
	RoundName string
}

func (e RoundStartedEvent) EventName() EventName {
	return roundStartedEventName
}

func (e RoundStartedEvent) ExecuteEvent(a *AggregateRoot) {
	a.board.startRound(e.RoundName)
}

type VoteCastEvent struct {
	PlayerId uuid.UUID
	Vote     string
}

func (e VoteCastEvent) ExecuteEvent(a *AggregateRoot) {
	a.board.castVote(e.PlayerId, e.Vote)
}

func (e VoteCastEvent) EventName() EventName {
	return voteCastEventName
}

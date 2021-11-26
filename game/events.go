package game

const (
	playerJoinedEventName = "PLAYER_JOINED"
)

type PlayerJoinedEvent struct {
}

func (e PlayerJoinedEvent) ExecuteEvent(a *AggregateRoot) {

}

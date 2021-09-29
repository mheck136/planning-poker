package game

type CreatedEvent struct {
	Name string
}

func (e CreatedEvent) apply(game Game) Game {
	game.Name = e.Name
	return game
}

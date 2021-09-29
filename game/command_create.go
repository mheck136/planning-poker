package game

import "fmt"

func NewCreateCommand(name string) CreateCommand {
	return CreateCommand{
		name: name,
		result: make(chan Game, 1),
	}
}

type CreateCommand struct {
	name   string
	result chan Game
}

func (c CreateCommand) reply(game Game, _ []Event) {
	c.result <- game
	close(c.result)
}

func (c CreateCommand) execute(game Game) ([]Event, error) {
	if len(game.Name) > 0 {
		return nil, fmt.Errorf("game has already been created")
	}
	return []Event{
		CreatedEvent{
			Name: c.name,
		},
	}, nil
}

func (c CreateCommand) Result() Game {
	return <- c.result
}



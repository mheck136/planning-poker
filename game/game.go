package game

import "encoding/json"

type Game struct {
	Name string
}

type Command interface {
	execute(game Game) ([]Event, error)
}

type Replier interface {
	reply(game Game, events []Event)
}

type Event interface {
	apply(game Game) Game
}

func Serialize(event Event) (string, []byte) {
	marshal, _ := json.Marshal(event)
	switch event.(type) {
	case CreatedEvent:
		return "CREATED", marshal
	}
	return "", nil
}

func Deserialize(eventName string, raw []byte) Event {
	switch eventName {
	case "CREATED":
		e := CreatedEvent{}
		_ = json.Unmarshal(raw, &e)
		return e
	}
	return nil
}

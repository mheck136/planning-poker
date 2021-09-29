package game

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCreateCommand_execute(t *testing.T) {
	c := NewCreateCommand("name")
	_, err := c.execute(Game{Name: "test"})
	assert.Error(t, err)

	c = NewCreateCommand("name")
	events, err := c.execute(Game{})
	assert.Nil(t, err)
	if assert.Equal(t, 1, len(events)) {
		createdEvent, ok := events[0].(CreatedEvent)
		if ok {
			assert.Equal(t, "name", createdEvent.Name)
		} else {
			t.Errorf("unexpected event type")
		}
	}
}

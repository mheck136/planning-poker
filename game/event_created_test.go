package game

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCreatedEvent_Apply(t *testing.T) {
	name := "TEST_NAME"
	e := CreatedEvent{Name: name}
	g := e.apply(Game{})
	assert.Equal(t, name, g.Name)
}

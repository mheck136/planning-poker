package gamedomain

import (
	"github.com/google/uuid"
)

type Game struct {
	Name            string
	JoinInformation JoinInformation
	Users           []User
	Admins          []User
}

func (g *Game) IsUserAllowedToJoin(token uuid.UUID, asAdmin bool) bool {
	if asAdmin {
		return g.JoinInformation.AdminToken == token
	} else {
		return g.JoinInformation.PlayerToken == token
	}
}

type JoinInformation struct {
	PlayerToken uuid.UUID
	AdminToken  uuid.UUID
}

type User struct {
	Id   uuid.UUID
	Name string
}

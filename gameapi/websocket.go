package gameapi

import (
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/rs/zerolog/log"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type liveUpdateController struct {
}

type GameUpdate struct {
	Users map[uuid.UUID]string `json:"users"`
}

func startNewConnectionHandler(conn *websocket.Conn) *connectionHandler {
	handler := connectionHandler{
		id:   uuid.New(),
		conn: conn,
		in:   make(chan GameUpdate, 1),
	}
	go handler.readLoop()
	go handler.writeLoop()
	return &handler
}

type connectionHandler struct {
	id   uuid.UUID
	conn *websocket.Conn
	in   chan GameUpdate
}

func (c *connectionHandler) readLoop() {
	for {
		if _, _, err := c.conn.NextReader(); err != nil {
			log.Err(err).Msg("error while reading from websocket connection")
			_ = c.conn.Close()
			break
		}
	}
}

func (c *connectionHandler) writeLoop() {
	for update := range c.in {
		if err := c.conn.WriteJSON(update); err != nil {
			log.Err(err).Msg("error while writing update to websocket connection")
			_ = c.conn.Close()
			break
		}
	}
}

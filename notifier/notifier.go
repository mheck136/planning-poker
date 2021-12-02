package notifier

import (
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/mheck136/planning-poker/game"
	"github.com/rs/zerolog/log"
)

type connection struct {
	conn     *websocket.Conn
	playerId uuid.UUID
	gameId   uuid.UUID
}

func New() *Notifier {
	n := &Notifier{
		newConnections:   make(chan connection),
		connectionClosed: make(chan *websocket.Conn),
		updates:          make(chan game.State),
		connections:      make(map[*websocket.Conn]connection),
	}
	go n.run()
	return n
}

type Notifier struct {
	newConnections   chan connection
	connectionClosed chan *websocket.Conn
	updates          chan game.State
	connections      map[*websocket.Conn]connection
}

func (n *Notifier) Notify(state game.State) {
	n.updates <- state
}

func (n *Notifier) HandleNewConnection(conn *websocket.Conn, playerId, gameId uuid.UUID) {
	n.newConnections <- connection{
		conn:     conn,
		playerId: playerId,
		gameId:   gameId,
	}
}

func (n *Notifier) handleClosedConnection(conn *websocket.Conn) {
	n.connectionClosed <- conn
}

func (n *Notifier) run() {
	for {
		select {
		case newConnection := <-n.newConnections:
			_, ok := n.connections[newConnection.conn]
			if !ok {
				log.Info().Str("gameId", newConnection.gameId.String()).Msg("adding connection for game")
				n.connections[newConnection.conn] = newConnection
				go func(conn *websocket.Conn) {
					for {
						if _, _, err := conn.NextReader(); err != nil {
							log.Err(err).Msg("error while reading from websocket connection")
							_ = conn.Close()
							delete(n.connections, conn)
							break
						}
					}
				}(newConnection.conn)
			}
		case connectionClosed := <-n.connectionClosed:
			delete(n.connections, connectionClosed)
		case update := <-n.updates:
			for _, conn := range n.connections {
				if conn.gameId == update.GameId {
					err := conn.conn.WriteJSON(update)
					if err != nil {
						log.Error().Err(err).Msg("error while sending update")
						_ = conn.conn.Close()
						delete(n.connections, conn.conn)
					}
				}
			}
		}
	}
}

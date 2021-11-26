package notifications

import (
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/mheck136/planning-poker/gamedomain"
	"github.com/rs/zerolog/log"
)

type registration struct {
	gameId   uuid.UUID
	playerId uuid.UUID
	conn     *websocket.Conn
}

type jsonNotification struct {
	gameId  uuid.UUID
	message interface{}
}

func New() *Notifier {
	n := &Notifier{
		registrations:             make(chan registration),
		jsonNotifications:         make(chan jsonNotification),
		connectionsByGameByPlayer: map[uuid.UUID]map[uuid.UUID]map[uuid.UUID]connection{},
	}
	go n.run()
	return n
}

type Notifier struct {
	registrations             chan registration
	jsonNotifications         chan jsonNotification
	connectionsByGameByPlayer map[uuid.UUID]map[uuid.UUID]map[uuid.UUID]connection // TODO
}

type connection struct {
	id       uuid.UUID
	playerId uuid.UUID
	gameId   uuid.UUID
	conn     *websocket.Conn
}

func (n *Notifier) run() {
	for {
		select {
		case registration := <-n.registrations:
			n.handleRegistration(registration)
		case jsonNotification := <-n.jsonNotifications:
			n.handleJsonNotification(jsonNotification)
		}
	}
}

func (n *Notifier) handleRegistration(r registration) {
	gameConnections, gameFound := n.connectionsByGameByPlayer[r.gameId]
	if !gameFound {
		gameConnections = map[uuid.UUID]map[uuid.UUID]connection{}
		n.connectionsByGameByPlayer[r.gameId] = gameConnections
	}
	playerConnections, playerFound := gameConnections[r.playerId]
	if !playerFound {
		playerConnections = map[uuid.UUID]connection{}
		gameConnections[r.playerId] = playerConnections
	}
	connectionId := uuid.New()
	playerConnections[connectionId] = connection{
		id:       connectionId,
		playerId: r.playerId,
		gameId:   r.gameId,
		conn:     r.conn,
	}
}

func (n *Notifier) handleJsonNotification(j jsonNotification) {
	gameConnections, gameFound := n.connectionsByGameByPlayer[j.gameId]
	if gameFound {
		for _, playerConnections := range gameConnections {
			for _, connection := range playerConnections {
				err := connection.conn.WriteJSON(j.message)
				if err != nil {
					log.Error().Err(err).Msg("error when sending message to client")
					// TODO: unsubscribe consumer
				}
			}
		}
	}
}

func (n *Notifier) PublishState(table gamedomain.Table) {
	n.getConnectionsForGameId()
	panic("implement me") // TODO
}

func (n *Notifier) Register(gameId, playerId uuid.UUID, conn *websocket.Conn) {
	n.registrations <- registration{
		gameId:   gameId,
		playerId: playerId,
		conn:     conn,
	}
}

func (n *Notifier) SendJsonNotification(gameId uuid.UUID, message interface{}) {
	n.jsonNotifications <- jsonNotification{
		gameId:  gameId,
		message: message,
	}
}

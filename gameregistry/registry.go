package gameregistry

import (
	"github.com/google/uuid"
	"github.com/mheck136/planning-poker/gamecommands"
	"sync"
)

func NewGameRegistry(proxyFactory ProxyFactory) *GameRegistry {
	return &GameRegistry{
		lock:         &sync.RWMutex{},
		proxies:      make(map[uuid.UUID]GameAggregateProxy),
		proxyFactory: proxyFactory,
	}
}

type GameAggregateProxy interface {
	SendJoinCommand(c gamecommands.JoinCommand)
}

type ProxyFactory func(id uuid.UUID) GameAggregateProxy

type GameRegistry struct {
	lock         *sync.RWMutex
	proxies      map[uuid.UUID]GameAggregateProxy
	proxyFactory ProxyFactory
}

func (g *GameRegistry) GetGameAggregateProxy(id uuid.UUID) GameAggregateProxy {
	g.lock.RLock()
	proxy, ok := g.proxies[id]
	g.lock.RUnlock()
	if !ok {
		g.lock.Lock()
		defer g.lock.Unlock()
		proxy = g.proxyFactory(id)
		g.proxies[id] = proxy
	}
	return proxy
}

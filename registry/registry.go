package registry

import (
	"github.com/google/uuid"
	"github.com/mheck136/planning-poker/aggregate"
	"sync"
)

func New(store EventStore, aggregateFactory func(id uuid.UUID) *aggregate.Aggregate) *Registry {
	return &Registry{
		lock:             &sync.RWMutex{},
		eventStore:       store,
		aggregateFactory: aggregateFactory,
		aggregateRoots:   make(map[uuid.UUID]*Root),
	}
}

type Registry struct {
	lock             *sync.RWMutex
	eventStore       EventStore
	aggregateFactory func(id uuid.UUID) *aggregate.Aggregate
	aggregateRoots   map[uuid.UUID]*Root
}

func (r *Registry) GetAggregateRoot(id uuid.UUID) *Root {
	r.lock.RLock()
	ar, ok := r.aggregateRoots[id]
	r.lock.RUnlock()
	if ok {
		return ar
	}
	r.lock.Lock()
	defer r.lock.Unlock()
	ar, ok = r.aggregateRoots[id]
	if ok {
		return ar
	}
	newAr := r.createAggregateRoot(id)
	r.aggregateRoots[id] = newAr
	return newAr
}

func (r *Registry) createAggregateRoot(id uuid.UUID) *Root {
	return newRoot(id, r.eventStore, r.aggregateFactory)
}

package game

import (
	"github.com/google/uuid"
	"sync"
)

func NewRegistry(store EventStore) *Registry {
	return &Registry{
		lock:           &sync.RWMutex{},
		eventStore:     store,
		aggregateRoots: make(map[uuid.UUID]*AggregateRoot),
	}
}

type Registry struct {
	lock           *sync.RWMutex
	eventStore     EventStore
	aggregateRoots map[uuid.UUID]*AggregateRoot
}

func (r *Registry) GetAggregateRoot(id uuid.UUID) *AggregateRoot {
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

func (r *Registry) createAggregateRoot(id uuid.UUID) *AggregateRoot {
	return newAggregateRoot(id, r.eventStore)
}

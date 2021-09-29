package game

import (
	"github.com/google/uuid"
	"sync"
)

type Registry struct {
	mu              sync.RWMutex
	aggregates      map[uuid.UUID]*Aggregate
	eventRepository EventRepository
}

func NewRegistry(repository EventRepository) *Registry {
	return &Registry{
		mu:              sync.RWMutex{},
		aggregates:      make(map[uuid.UUID]*Aggregate),
		eventRepository: repository,
	}
}

func (r *Registry) GetAggregate(id uuid.UUID) *Aggregate {
	r.mu.RLock()
	aggregate, ok := r.aggregates[id]
	r.mu.RUnlock()
	if !ok {
		r.mu.Lock()
		defer r.mu.Unlock()
		aggregate, ok = r.aggregates[id]
	}
	if !ok {
		aggregate = &Aggregate{
			id: id,
			cmdIn: make(chan Command, 10),
			repository: r.eventRepository,
		}
		r.aggregates[id] = aggregate
		go aggregate.run()
	}
	return aggregate
}

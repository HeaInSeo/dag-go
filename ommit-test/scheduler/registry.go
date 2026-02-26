package scheduler

import "sync"

// TODO 이거 수정해야 하지 않을까?

// Registry is the storage contract for Actor lookup by spawnID.
type Registry interface {
	Get(spawnID string) (*Actor, bool)
	Put(spawnID string, a *Actor)
	Delete(spawnID string)
}

// InMemoryRegistry is a thread-safe in-process implementation of Registry.
type InMemoryRegistry struct {
	mu   sync.RWMutex
	data map[string]*Actor
}

// NewInMemoryRegistry returns an initialized InMemoryRegistry.
func NewInMemoryRegistry() *InMemoryRegistry {
	return &InMemoryRegistry{data: make(map[string]*Actor)}
}

// Get returns the Actor registered under spawnID, or (nil, false) if absent.
func (r *InMemoryRegistry) Get(spawnID string) (*Actor, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	a, ok := r.data[spawnID]
	return a, ok
}

// Put registers the Actor under spawnID, overwriting any existing entry.
func (r *InMemoryRegistry) Put(spawnID string, a *Actor) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.data[spawnID] = a
}

// Delete removes the Actor registered under spawnID.
func (r *InMemoryRegistry) Delete(spawnID string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.data, spawnID)
}

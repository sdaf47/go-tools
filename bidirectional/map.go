package bidirectional

import "sync"

type BiMap[T comparable, E comparable] struct {
	s       sync.RWMutex
	forward map[T]E
	inverse map[E]T
}

// NewBiMap returns an empty mutable biMap
func NewBiMap[T comparable, E comparable]() *BiMap[T, E] {
	return &BiMap[T, E]{forward: make(map[T]E), inverse: make(map[E]T)}
}

// Insert puts a key value pair into BiMap
func (b *BiMap[T, E]) Insert(k T, v E) {
	b.s.Lock()
	defer b.s.Unlock()

	b.forward[k] = v
	b.inverse[v] = k
}

// Exists checks a key exists in the BiMap
func (b *BiMap[T, E]) Exists(k T) bool {
	b.s.RLock()
	defer b.s.RUnlock()

	_, ok := b.forward[k]
	return ok
}

// ExistsInverse checks a value exists in the BiMap
func (b *BiMap[T, E]) ExistsInverse(v E) bool {
	b.s.RLock()
	defer b.s.RUnlock()

	_, ok := b.inverse[v]
	return ok
}

// Get returns the value for a given key in the BiMap and the element was present
func (b *BiMap[T, E]) Get(k T) (E, bool) {
	b.s.RLock()
	defer b.s.RUnlock()

	v, ok := b.forward[k]

	return v, ok
}

// Inverse returns the key by value
func (b *BiMap[T, E]) Inverse(v E) (T, bool) {
	b.s.RLock()
	defer b.s.RUnlock()
	k, ok := b.inverse[v]

	return k, ok
}

// Delete deletes item by key
func (b *BiMap[T, E]) Delete(k T) {
	b.s.Lock()
	defer b.s.Unlock()

	v, ok := b.forward[k]
	if !ok {
		return
	}

	delete(b.inverse, v)
	delete(b.forward, k)
}

// DeleteInverse deletes item by its value
func (b *BiMap[T, E]) DeleteInverse(v E) {
	b.s.Lock()
	defer b.s.Unlock()

	k, ok := b.inverse[v]
	if !ok {
		return
	}

	delete(b.inverse, v)
	delete(b.forward, k)
}

// Len returns the number of elements in the BiMap
func (b *BiMap[T, E]) Len() int {
	b.s.RLock()
	defer b.s.RUnlock()
	return len(b.forward)
}

// ForwardMap returns a regular go map from the BiMap
func (b *BiMap[T, E]) ForwardMap() map[T]E {
	return b.forward
}

// InverseMap returns a reverted go map from the BiMap
func (b *BiMap[T, E]) InverseMap() map[E]T {
	return b.inverse
}

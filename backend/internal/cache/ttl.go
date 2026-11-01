// Package cache implementa um cache em memória com expiração por TTL.
package cache

import (
	"sync"
	"time"
)

type entry[V any] struct {
	value     V
	expiresAt time.Time
}

// TTL é um mapa protegido por mutex onde cada valor expira depois de ttl.
// Itens expirados são removidos na leitura e por Sweep.
type TTL[K comparable, V any] struct {
	mu    sync.Mutex
	items map[K]entry[V]
	ttl   time.Duration
	now   func() time.Time
}

// New cria um cache com o TTL dado.
func New[K comparable, V any](ttl time.Duration) *TTL[K, V] {
	return &TTL[K, V]{items: make(map[K]entry[V]), ttl: ttl, now: time.Now}
}

// Get devolve o valor e true se a chave existir e não tiver expirado.
func (c *TTL[K, V]) Get(key K) (V, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	e, ok := c.items[key]
	if !ok {
		var zero V
		return zero, false
	}
	if c.now().After(e.expiresAt) {
		delete(c.items, key)
		var zero V
		return zero, false
	}
	return e.value, true
}

// Set guarda o valor com o TTL do cache.
func (c *TTL[K, V]) Set(key K, value V) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.items[key] = entry[V]{value: value, expiresAt: c.now().Add(c.ttl)}
}

// Sweep remove todos os itens expirados e devolve quantos saíram.
func (c *TTL[K, V]) Sweep() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	now := c.now()
	removed := 0
	for k, e := range c.items {
		if now.After(e.expiresAt) {
			delete(c.items, k)
			removed++
		}
	}
	return removed
}

// Len devolve o número de itens guardados (incluindo os ainda não varridos).
func (c *TTL[K, V]) Len() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return len(c.items)
}

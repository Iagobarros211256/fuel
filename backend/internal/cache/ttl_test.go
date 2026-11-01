package cache

import (
	"testing"
	"time"
)

func TestTTLExpira(t *testing.T) {
	clock := time.Date(2026, 11, 1, 10, 0, 0, 0, time.UTC)
	c := New[string, int](time.Minute)
	c.now = func() time.Time { return clock }

	c.Set("a", 1)
	if v, ok := c.Get("a"); !ok || v != 1 {
		t.Fatalf("esperava 1, true; veio %d, %v", v, ok)
	}

	clock = clock.Add(61 * time.Second)
	if _, ok := c.Get("a"); ok {
		t.Fatal("o item deveria ter expirado")
	}
	if c.Len() != 0 {
		t.Fatal("Get deveria remover o item expirado")
	}
}

func TestSweep(t *testing.T) {
	clock := time.Date(2026, 11, 1, 10, 0, 0, 0, time.UTC)
	c := New[string, int](time.Minute)
	c.now = func() time.Time { return clock }

	c.Set("velho", 1)
	clock = clock.Add(30 * time.Second)
	c.Set("novo", 2)
	clock = clock.Add(45 * time.Second) // "velho" passou de 1 min, "novo" não

	if n := c.Sweep(); n != 1 {
		t.Fatalf("esperava remover 1, removeu %d", n)
	}
	if _, ok := c.Get("novo"); !ok {
		t.Fatal("'novo' ainda deveria existir")
	}
}

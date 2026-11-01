package station

import (
	"context"
	"errors"
	"testing"
	"time"
)

type stubProvider struct {
	calls    int
	stations []Station
	err      error
}

func (s *stubProvider) Nearby(context.Context, Query) ([]Station, error) {
	s.calls++
	return s.stations, s.err
}

func TestNearbyOrdenaPorDistancia(t *testing.T) {
	p := &stubProvider{stations: []Station{
		{ID: "longe", Lat: -3.7400, Lng: -38.5267},
		{ID: "perto", Lat: -3.7320, Lng: -38.5267},
		{ID: "meio", Lat: -3.7350, Lng: -38.5267},
	}}
	svc := NewService(p, time.Minute)

	got, err := svc.Nearby(context.Background(), Query{Lat: -3.7319, Lng: -38.5267, RadiusMeters: 3000})
	if err != nil {
		t.Fatal(err)
	}
	ids := []string{got[0].ID, got[1].ID, got[2].ID}
	want := []string{"perto", "meio", "longe"}
	for i := range want {
		if ids[i] != want[i] {
			t.Fatalf("ordem errada: %v, esperado %v", ids, want)
		}
	}
	if got[0].DistanceMeters <= 0 || got[0].Prices == nil {
		t.Fatal("distância deveria ser calculada e Prices nunca nil")
	}
}

func TestNearbyUsaCache(t *testing.T) {
	p := &stubProvider{stations: []Station{{ID: "a", Lat: -3.73, Lng: -38.52}}}
	svc := NewService(p, time.Minute)
	q := Query{Lat: -3.73191, Lng: -38.52671, RadiusMeters: 3000}

	for range 3 {
		if _, err := svc.Nearby(context.Background(), q); err != nil {
			t.Fatal(err)
		}
	}
	// ponto a poucos metros: mesma chave de cache
	q.Lat, q.Lng = -3.73194, -38.52668
	if _, err := svc.Nearby(context.Background(), q); err != nil {
		t.Fatal(err)
	}
	if p.calls != 1 {
		t.Fatalf("provider chamado %d vezes, esperado 1", p.calls)
	}
}

func TestNearbyErroDoProvider(t *testing.T) {
	p := &stubProvider{err: errors.New("timeout")}
	svc := NewService(p, time.Minute)

	_, err := svc.Nearby(context.Background(), Query{Lat: 0, Lng: 0, RadiusMeters: 1000})
	if !errors.Is(err, ErrUpstream) {
		t.Fatalf("esperava ErrUpstream, veio %v", err)
	}
}

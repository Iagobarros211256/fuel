package httpapi

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Iagobarros211256/abastece/backend/internal/station"
)

type stubFinder struct {
	got station.Query
	err error
}

func (s *stubFinder) Nearby(_ context.Context, q station.Query) ([]station.Station, error) {
	s.got = q
	if s.err != nil {
		return nil, s.err
	}
	return []station.Station{{ID: "1", Name: "Posto", Prices: []station.FuelPrice{}}}, nil
}

func newTestHandler(f NearbyFinder) http.Handler {
	return New(f, slog.New(slog.NewTextHandler(io.Discard, nil)))
}

func TestNearbyOK(t *testing.T) {
	f := &stubFinder{}
	rec := httptest.NewRecorder()
	newTestHandler(f).ServeHTTP(rec, httptest.NewRequest("GET", "/v1/stations/nearby?lat=-3.73&lng=-38.52", nil))

	if rec.Code != http.StatusOK {
		t.Fatalf("status %d: %s", rec.Code, rec.Body)
	}
	if f.got.RadiusMeters != defaultRadius {
		t.Fatalf("raio padrão deveria ser %v, veio %v", defaultRadius, f.got.RadiusMeters)
	}
	var body nearbyResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatal(err)
	}
	if body.Count != 1 || body.Stations[0].ID != "1" {
		t.Fatalf("corpo inesperado: %+v", body)
	}
	if rec.Header().Get("Access-Control-Allow-Origin") != "*" {
		t.Fatal("faltou header de CORS")
	}
}

func TestNearbyValidacao(t *testing.T) {
	casos := map[string]string{
		"sem lat":         "/v1/stations/nearby?lng=-38.52",
		"lat texto":       "/v1/stations/nearby?lat=abc&lng=-38.52",
		"lat fora":        "/v1/stations/nearby?lat=91&lng=-38.52",
		"lng fora":        "/v1/stations/nearby?lat=-3&lng=181",
		"raio pequeno":    "/v1/stations/nearby?lat=-3&lng=-38&radius=50",
		"raio muito alto": "/v1/stations/nearby?lat=-3&lng=-38&radius=60000",
	}
	for nome, url := range casos {
		t.Run(nome, func(t *testing.T) {
			rec := httptest.NewRecorder()
			newTestHandler(&stubFinder{}).ServeHTTP(rec, httptest.NewRequest("GET", url, nil))
			if rec.Code != http.StatusBadRequest {
				t.Fatalf("esperava 400, veio %d", rec.Code)
			}
		})
	}
}

func TestNearbyUpstreamFalhou(t *testing.T) {
	f := &stubFinder{err: errors.Join(station.ErrUpstream, errors.New("timeout"))}
	rec := httptest.NewRecorder()
	newTestHandler(f).ServeHTTP(rec, httptest.NewRequest("GET", "/v1/stations/nearby?lat=-3&lng=-38", nil))
	if rec.Code != http.StatusBadGateway {
		t.Fatalf("esperava 502, veio %d", rec.Code)
	}
}

func TestHealthz(t *testing.T) {
	rec := httptest.NewRecorder()
	newTestHandler(&stubFinder{}).ServeHTTP(rec, httptest.NewRequest("GET", "/healthz", nil))
	if rec.Code != http.StatusOK {
		t.Fatalf("status %d", rec.Code)
	}
}

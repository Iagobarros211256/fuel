// Package httpapi expõe a API HTTP consumida pelo app.
package httpapi

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/Iagobarros211256/abastece/backend/internal/station"
)

const (
	defaultRadius = 3000.0
	minRadius     = 100.0
	maxRadius     = 50000.0 // limite do Places API
)

// NearbyFinder é o que o handler precisa do serviço (facilita testar).
type NearbyFinder interface {
	Nearby(ctx context.Context, q station.Query) ([]station.Station, error)
}

// New monta o roteador com middlewares.
func New(svc NearbyFinder, log *slog.Logger) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, _ *http.Request) {
		writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
	})
	mux.HandleFunc("GET /v1/stations/nearby", nearbyHandler(svc, log))
	return logging(log, cors(mux))
}

type nearbyResponse struct {
	Count    int               `json:"count"`
	Stations []station.Station `json:"stations"`
}

func nearbyHandler(svc NearbyFinder, log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		q, err := parseQuery(r)
		if err != nil {
			writeError(w, http.StatusBadRequest, err.Error())
			return
		}

		stations, err := svc.Nearby(r.Context(), q)
		if err != nil {
			log.Error("nearby falhou", "err", err)
			if errors.Is(err, station.ErrUpstream) {
				writeError(w, http.StatusBadGateway, "não foi possível consultar os postos agora")
				return
			}
			writeError(w, http.StatusInternalServerError, "erro interno")
			return
		}
		writeJSON(w, http.StatusOK, nearbyResponse{Count: len(stations), Stations: stations})
	}
}

func parseQuery(r *http.Request) (station.Query, error) {
	v := r.URL.Query()

	lat, err := parseFloat(v.Get("lat"), "lat")
	if err != nil {
		return station.Query{}, err
	}
	lng, err := parseFloat(v.Get("lng"), "lng")
	if err != nil {
		return station.Query{}, err
	}
	if lat < -90 || lat > 90 {
		return station.Query{}, errors.New("lat deve estar entre -90 e 90")
	}
	if lng < -180 || lng > 180 {
		return station.Query{}, errors.New("lng deve estar entre -180 e 180")
	}

	radius := defaultRadius
	if s := v.Get("radius"); s != "" {
		radius, err = parseFloat(s, "radius")
		if err != nil {
			return station.Query{}, err
		}
		if radius < minRadius || radius > maxRadius {
			return station.Query{}, errors.New("radius deve estar entre 100 e 50000 metros")
		}
	}
	return station.Query{Lat: lat, Lng: lng, RadiusMeters: radius}, nil
}

func parseFloat(s, name string) (float64, error) {
	if s == "" {
		return 0, errors.New(name + " é obrigatório")
	}
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0, errors.New(name + " deve ser um número")
	}
	return f, nil
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}

// cors libera chamadas do Expo web em desenvolvimento.
func cors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}

type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (s *statusRecorder) WriteHeader(code int) {
	s.status = code
	s.ResponseWriter.WriteHeader(code)
}

func logging(log *slog.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rec := &statusRecorder{ResponseWriter: w, status: http.StatusOK}
		next.ServeHTTP(rec, r)
		log.Info("http", "method", r.Method, "path", r.URL.Path, "status", rec.status, "dur", time.Since(start))
	})
}

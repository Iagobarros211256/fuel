// Package station contém o modelo de posto e o serviço que busca postos próximos.
package station

import (
	"context"
	"errors"
	"sort"
	"time"

	"github.com/Iagobarros211256/abastece/backend/internal/cache"
	"github.com/Iagobarros211256/abastece/backend/internal/geo"
)

// Combustíveis normalizados usados pelo app.
const (
	FuelGasolinaComum     = "gasolina_comum"
	FuelGasolinaAditivada = "gasolina_aditivada"
	FuelGasolinaPremium   = "gasolina_premium"
	FuelEtanol            = "etanol"
	FuelDiesel            = "diesel"
	FuelGNV               = "gnv"
	FuelGLP               = "glp"
)

// Fontes de preço (ver ADR 0002).
const (
	SourceGoogle  = "google"
	SourceANP     = "anp"
	SourceUsuario = "usuario"
	SourceFake    = "fake"
)

// FuelPrice é o preço de um combustível num posto, sempre com origem e data.
type FuelPrice struct {
	Fuel        string    `json:"fuel"`
	Price       float64   `json:"price"`
	Currency    string    `json:"currency"`
	Source      string    `json:"source"`
	CollectedAt time.Time `json:"collectedAt"`
}

// Station é um posto de combustível.
type Station struct {
	ID             string      `json:"id"`
	Name           string      `json:"name"`
	Address        string      `json:"address"`
	Lat            float64     `json:"lat"`
	Lng            float64     `json:"lng"`
	Rating         float64     `json:"rating,omitempty"`
	DistanceMeters float64     `json:"distanceMeters"`
	Prices         []FuelPrice `json:"prices"`
}

// Query descreve uma busca por proximidade.
type Query struct {
	Lat, Lng     float64
	RadiusMeters float64
}

// Provider é qualquer fonte de postos (Google Places, fake, futuramente ANP).
type Provider interface {
	Nearby(ctx context.Context, q Query) ([]Station, error)
}

// ErrUpstream indica falha no provedor externo.
var ErrUpstream = errors.New("falha ao consultar o provedor de postos")

// Service busca postos usando um Provider, com cache e ordenação por distância.
type Service struct {
	provider Provider
	cache    *cache.TTL[string, []Station]
}

// NewService cria o serviço. ttl controla por quanto tempo uma busca fica em cache.
func NewService(p Provider, ttl time.Duration) *Service {
	return &Service{provider: p, cache: cache.New[string, []Station](ttl)}
}

// Nearby devolve os postos no raio, ordenados do mais perto para o mais longe.
// A distância é sempre recalculada a partir do ponto exato do usuário, mesmo
// quando o resultado vem do cache (que usa a coordenada arredondada).
func (s *Service) Nearby(ctx context.Context, q Query) ([]Station, error) {
	key := geo.CacheKey(q.Lat, q.Lng, q.RadiusMeters)

	stations, ok := s.cache.Get(key)
	if !ok {
		var err error
		stations, err = s.provider.Nearby(ctx, q)
		if err != nil {
			return nil, errors.Join(ErrUpstream, err)
		}
		s.cache.Set(key, stations)
	}
	return withDistance(stations, q), nil
}

// SweepCache remove entradas expiradas; chamado periodicamente pelo main.
func (s *Service) SweepCache() int { return s.cache.Sweep() }

func withDistance(in []Station, q Query) []Station {
	out := make([]Station, len(in)) // cópia: não altera o slice que está no cache
	copy(out, in)
	for i := range out {
		out[i].DistanceMeters = geo.DistanceMeters(q.Lat, q.Lng, out[i].Lat, out[i].Lng)
		if out[i].Prices == nil {
			out[i].Prices = []FuelPrice{}
		}
	}
	sort.SliceStable(out, func(i, j int) bool {
		return out[i].DistanceMeters < out[j].DistanceMeters
	})
	return out
}

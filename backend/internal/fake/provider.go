// Package fake gera postos falsos em volta do usuário, para desenvolver sem chave do Google.
package fake

import (
	"context"
	"fmt"
	"math"
	"time"

	"github.com/Iagobarros211256/fuel/backend/internal/station"
)

// Provider devolve sempre os mesmos 8 postos espalhados no raio pedido.
type Provider struct {
	Now func() time.Time
}

// Nearby implementa station.Provider.
func (p Provider) Nearby(_ context.Context, q station.Query) ([]station.Station, error) {
	now := time.Now
	if p.Now != nil {
		now = p.Now
	}

	const n = 8
	out := make([]station.Station, 0, n)
	for i := range n {
		// distribui em espiral: ângulo crescente, distância de 15% a 85% do raio
		angle := float64(i) * 2 * math.Pi / n
		dist := q.RadiusMeters * (0.15 + 0.7*float64(i)/float64(n-1))
		lat, lng := offset(q.Lat, q.Lng, dist, angle)

		gasolina := 6.09 + 0.07*float64((i*3)%n)
		out = append(out, station.Station{
			ID:      fmt.Sprintf("fake-%d", i+1),
			Name:    fmt.Sprintf("Posto Exemplo %d", i+1),
			Address: "Endereço fictício (provider fake)",
			Lat:     lat,
			Lng:     lng,
			Rating:  3.5 + 0.2*float64(i%5),
			Prices: []station.FuelPrice{
				price(station.FuelGasolinaComum, gasolina, now()),
				price(station.FuelEtanol, gasolina*0.72, now()),
			},
		})
	}
	return out, nil
}

func price(fuel string, v float64, at time.Time) station.FuelPrice {
	return station.FuelPrice{
		Fuel:        fuel,
		Price:       math.Round(v*100) / 100,
		Currency:    "BRL",
		Source:      station.SourceFake,
		CollectedAt: at.Add(-24 * time.Hour).UTC().Truncate(time.Hour),
	}
}

// offset desloca um ponto dist metros na direção angle (radianos, 0 = leste).
func offset(lat, lng, dist, angle float64) (float64, float64) {
	const metersPerDegLat = 111_320.0
	dLat := dist * math.Sin(angle) / metersPerDegLat
	dLng := dist * math.Cos(angle) / (metersPerDegLat * math.Cos(lat*math.Pi/180))
	return lat + dLat, lng + dLng
}

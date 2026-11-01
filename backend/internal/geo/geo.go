// Package geo reúne cálculos geográficos simples usados pelo backend.
package geo

import (
	"fmt"
	"math"
)

// earthRadiusMeters é o raio médio da Terra (esfera), suficiente para distâncias urbanas.
const earthRadiusMeters = 6_371_000.0

// DistanceMeters calcula a distância em metros entre dois pontos pela fórmula de haversine.
func DistanceMeters(lat1, lng1, lat2, lng2 float64) float64 {
	φ1 := toRad(lat1)
	φ2 := toRad(lat2)
	Δφ := toRad(lat2 - lat1)
	Δλ := toRad(lng2 - lng1)

	a := math.Sin(Δφ/2)*math.Sin(Δφ/2) +
		math.Cos(φ1)*math.Cos(φ2)*math.Sin(Δλ/2)*math.Sin(Δλ/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	return earthRadiusMeters * c
}

// CacheKey arredonda a coordenada para 3 casas decimais (~110 m) e junta com o raio.
// Usuários a poucos metros um do outro caem na mesma chave e reaproveitam o cache.
func CacheKey(lat, lng, radiusMeters float64) string {
	return fmt.Sprintf("%.3f,%.3f,%d", round3(lat), round3(lng), int(radiusMeters))
}

func round3(v float64) float64 {
	r := math.Round(v*1000) / 1000
	if r == 0 {
		return 0 // evita "-0.000"
	}
	return r
}

func toRad(deg float64) float64 { return deg * math.Pi / 180 }

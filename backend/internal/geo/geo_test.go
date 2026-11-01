package geo

import (
	"math"
	"testing"
)

func TestDistanceMeters(t *testing.T) {
	tests := []struct {
		name                   string
		lat1, lng1, lat2, lng2 float64
		want, tolerance        float64
	}{
		{"mesmo ponto", -3.7319, -38.5267, -3.7319, -38.5267, 0, 0.001},
		// dois pontos no centro de Fortaleza, ~1 km em linha reta
		{"centro de Fortaleza", -3.7272, -38.5274, -3.7225, -38.5195, 1020, 60},
		// 1 grau de latitude ~ 111,2 km
		{"um grau de latitude", 0, 0, 1, 0, 111_195, 50},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := DistanceMeters(tt.lat1, tt.lng1, tt.lat2, tt.lng2)
			if math.Abs(got-tt.want) > tt.tolerance {
				t.Fatalf("got %.1f m, want %.1f ± %.1f", got, tt.want, tt.tolerance)
			}
		})
	}
}

func TestCacheKey(t *testing.T) {
	a := CacheKey(-3.73191, -38.52671, 3000)
	b := CacheKey(-3.73194, -38.52668, 3000)
	if a != b {
		t.Fatalf("pontos vizinhos deveriam ter a mesma chave: %q vs %q", a, b)
	}
	if CacheKey(-3.7319, -38.5267, 3000) == CacheKey(-3.7319, -38.5267, 5000) {
		t.Fatal("raios diferentes deveriam gerar chaves diferentes")
	}
	if got := CacheKey(-0.0001, 0.0001, 100); got != "0.000,0.000,100" {
		t.Fatalf("zero negativo não normalizado: %q", got)
	}
}

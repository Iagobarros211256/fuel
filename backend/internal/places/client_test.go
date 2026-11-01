package places

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Iagobarros211256/abastece/backend/internal/station"
)

const respostaExemplo = `{
  "places": [
    {
      "id": "abc123",
      "displayName": {"text": "Posto Teste", "languageCode": "pt-BR"},
      "formattedAddress": "Av. Beira Mar, 100 - Fortaleza, CE",
      "location": {"latitude": -3.7250, "longitude": -38.4950},
      "rating": 4.3,
      "fuelOptions": {
        "fuelPrices": [
          {"type": "REGULAR_UNLEADED", "price": {"currencyCode": "BRL", "units": "6", "nanos": 290000000}, "updateTime": "2026-10-30T12:00:00Z"},
          {"type": "E100", "price": {"currencyCode": "BRL", "units": "4", "nanos": 590000000}, "updateTime": "2026-10-30T12:00:00Z"}
        ]
      }
    },
    {
      "id": "sem-preco",
      "displayName": {"text": "Posto Sem Preço"},
      "location": {"latitude": -3.73, "longitude": -38.50}
    }
  ]
}`

func TestNearbyMontaRequisicaoEParseia(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || r.URL.Path != "/places:searchNearby" {
			t.Errorf("rota inesperada: %s %s", r.Method, r.URL.Path)
		}
		if r.Header.Get("X-Goog-Api-Key") != "chave-teste" {
			t.Error("faltou o header X-Goog-Api-Key")
		}
		if !strings.Contains(r.Header.Get("X-Goog-FieldMask"), "places.fuelOptions") {
			t.Error("field mask deveria pedir places.fuelOptions")
		}

		var body searchNearbyRequest
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatal(err)
		}
		if body.IncludedTypes[0] != "gas_station" || body.LocationRestriction.Circle.Radius != 3000 {
			t.Errorf("corpo inesperado: %+v", body)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(respostaExemplo))
	}))
	defer srv.Close()

	c := New("chave-teste", WithBaseURL(srv.URL))
	got, err := c.Nearby(context.Background(), station.Query{Lat: -3.7319, Lng: -38.5267, RadiusMeters: 3000})
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 2 {
		t.Fatalf("esperava 2 postos, veio %d", len(got))
	}

	p := got[0]
	if p.Name != "Posto Teste" || p.Rating != 4.3 || len(p.Prices) != 2 {
		t.Fatalf("posto mal convertido: %+v", p)
	}
	if p.Prices[0].Fuel != station.FuelGasolinaComum || p.Prices[0].Price != 6.29 || p.Prices[0].Source != station.SourceGoogle {
		t.Fatalf("preço mal convertido: %+v", p.Prices[0])
	}
	if p.Prices[1].Fuel != station.FuelEtanol {
		t.Fatalf("E100 deveria virar etanol, veio %q", p.Prices[1].Fuel)
	}
	if got[1].Prices == nil || len(got[1].Prices) != 0 {
		t.Fatal("posto sem fuelOptions deveria ter lista de preços vazia")
	}
}

func TestNearbyErroHTTP(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, `{"error":{"message":"API key not valid"}}`, http.StatusForbidden)
	}))
	defer srv.Close()

	c := New("ruim", WithBaseURL(srv.URL))
	_, err := c.Nearby(context.Background(), station.Query{RadiusMeters: 1000})
	if err == nil || !strings.Contains(err.Error(), "403") {
		t.Fatalf("esperava erro com status 403, veio %v", err)
	}
}

func TestMoneyValue(t *testing.T) {
	m := money{Units: "5", Nanos: 990000000}
	if v, _ := m.value(); v != 5.99 {
		t.Fatalf("esperava 5.99, veio %v", v)
	}
	if _, err := (money{Units: "x"}).value(); err == nil {
		t.Fatal("units inválido deveria dar erro")
	}
}

// Package places é o cliente HTTP do Google Places API (New) — Nearby Search.
//
// Docs: https://developers.google.com/maps/documentation/places/web-service/nearby-search
package places

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/Iagobarros211256/fuel/backend/internal/station"
)

const defaultBaseURL = "https://places.googleapis.com/v1"

// fieldMask define o que o Google devolve (e quanto cobra). Só pedimos o necessário.
var fieldMask = strings.Join([]string{
	"places.id",
	"places.displayName",
	"places.formattedAddress",
	"places.location",
	"places.rating",
	"places.fuelOptions",
}, ",")

// Client chama o endpoint places:searchNearby.
type Client struct {
	apiKey  string
	baseURL string
	http    *http.Client
}

// Option configura o Client.
type Option func(*Client)

// WithBaseURL troca a URL base (usado nos testes com httptest).
func WithBaseURL(u string) Option { return func(c *Client) { c.baseURL = strings.TrimRight(u, "/") } }

// WithHTTPClient troca o http.Client.
func WithHTTPClient(h *http.Client) Option { return func(c *Client) { c.http = h } }

// New cria o cliente com a chave da API.
func New(apiKey string, opts ...Option) *Client {
	c := &Client{
		apiKey:  apiKey,
		baseURL: defaultBaseURL,
		http:    &http.Client{Timeout: 10 * time.Second},
	}
	for _, o := range opts {
		o(c)
	}
	return c
}

// --- formato da requisição/resposta (só os campos que usamos) ---

type latLng struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type searchNearbyRequest struct {
	IncludedTypes       []string `json:"includedTypes"`
	MaxResultCount      int      `json:"maxResultCount"`
	RankPreference      string   `json:"rankPreference"`
	LanguageCode        string   `json:"languageCode"`
	RegionCode          string   `json:"regionCode"`
	LocationRestriction struct {
		Circle struct {
			Center latLng  `json:"center"`
			Radius float64 `json:"radius"`
		} `json:"circle"`
	} `json:"locationRestriction"`
}

// money segue google.type.Money: units vem como string (int64 em JSON) e nanos como int.
type money struct {
	CurrencyCode string `json:"currencyCode"`
	Units        string `json:"units"`
	Nanos        int64  `json:"nanos"`
}

func (m money) value() (float64, error) {
	var units int64
	if m.Units != "" {
		var err error
		units, err = strconv.ParseInt(m.Units, 10, 64)
		if err != nil {
			return 0, fmt.Errorf("units inválido %q: %w", m.Units, err)
		}
	}
	return float64(units) + float64(m.Nanos)/1e9, nil
}

type place struct {
	ID          string `json:"id"`
	DisplayName struct {
		Text string `json:"text"`
	} `json:"displayName"`
	FormattedAddress string  `json:"formattedAddress"`
	Location         latLng  `json:"location"`
	Rating           float64 `json:"rating"`
	FuelOptions      *struct {
		FuelPrices []struct {
			Type       string    `json:"type"`
			Price      money     `json:"price"`
			UpdateTime time.Time `json:"updateTime"`
		} `json:"fuelPrices"`
	} `json:"fuelOptions"`
}

type searchNearbyResponse struct {
	Places []place `json:"places"`
}

// Nearby implementa station.Provider.
func (c *Client) Nearby(ctx context.Context, q station.Query) ([]station.Station, error) {
	body := searchNearbyRequest{
		IncludedTypes:  []string{"gas_station"},
		MaxResultCount: 20, // máximo permitido pela API
		RankPreference: "DISTANCE",
		LanguageCode:   "pt-BR",
		RegionCode:     "BR",
	}
	body.LocationRestriction.Circle.Center = latLng{q.Lat, q.Lng}
	body.LocationRestriction.Circle.Radius = q.RadiusMeters

	buf, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/places:searchNearby", bytes.NewReader(buf))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Goog-Api-Key", c.apiKey)
	req.Header.Set("X-Goog-FieldMask", fieldMask)

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		msg, _ := io.ReadAll(io.LimitReader(resp.Body, 2048))
		return nil, fmt.Errorf("places API respondeu %d: %s", resp.StatusCode, strings.TrimSpace(string(msg)))
	}

	var out searchNearbyResponse
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return nil, fmt.Errorf("decodificando resposta do places: %w", err)
	}
	return toStations(out.Places), nil
}

func toStations(ps []place) []station.Station {
	stations := make([]station.Station, 0, len(ps))
	for _, p := range ps {
		s := station.Station{
			ID:      p.ID,
			Name:    p.DisplayName.Text,
			Address: p.FormattedAddress,
			Lat:     p.Location.Latitude,
			Lng:     p.Location.Longitude,
			Rating:  p.Rating,
			Prices:  []station.FuelPrice{},
		}
		if p.FuelOptions != nil {
			for _, fp := range p.FuelOptions.FuelPrices {
				v, err := fp.Price.value()
				if err != nil || v <= 0 {
					continue // preço ilegível não entra
				}
				s.Prices = append(s.Prices, station.FuelPrice{
					Fuel:        normalizeFuel(fp.Type),
					Price:       v,
					Currency:    fp.Price.CurrencyCode,
					Source:      station.SourceGoogle,
					CollectedAt: fp.UpdateTime,
				})
			}
		}
		stations = append(stations, s)
	}
	return stations
}

// normalizeFuel traduz o enum do Google para os nomes usados no app.
// Mapeamento provisório: revisar quando houver dados reais do Brasil (S1-10).
func normalizeFuel(googleType string) string {
	switch googleType {
	case "REGULAR_UNLEADED":
		return station.FuelGasolinaComum
	case "MIDGRADE":
		return station.FuelGasolinaAditivada
	case "PREMIUM":
		return station.FuelGasolinaPremium
	case "E100":
		return station.FuelEtanol
	case "DIESEL", "DIESEL_PLUS", "BIO_DIESEL":
		return station.FuelDiesel
	case "METHANE":
		return station.FuelGNV
	case "LPG":
		return station.FuelGLP
	default:
		return strings.ToLower(googleType)
	}
}

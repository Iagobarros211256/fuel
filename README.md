# Fuel

App que encontra postos de combustível perto de você e (em breve) mostra os 5 mais
baratos, com camadas de radares, semáforos e hospitais.

Projeto de estudo — trilha "me desafiando".

```
fuel/
├── backend/   API em Go (stdlib) — integra Google Places, ANP e OpenStreetMap
├── app/       App mobile em Expo (React Native + TypeScript)
└── docs/      PRD, ADRs e backlog
```

## Status

**Sprint 1 — postos próximos.** Ver [docs/BACKLOG.md](docs/BACKLOG.md).

## Rodando o backend

Requer Go 1.24+.

```bash
cd backend
go test ./...
go run ./cmd/api            # sem chave: usa postos fictícios (provider fake)
```

Com dados reais do Google:

```bash
export GOOGLE_MAPS_API_KEY=sua_chave   # com "Places API (New)" habilitada
go run ./cmd/api
```

| Variável | Padrão | O que faz |
|---|---|---|
| `PORT` | `8080` | porta HTTP |
| `GOOGLE_MAPS_API_KEY` | — | chave do Google; se vazia, usa o fake |
| `PROVIDER` | `google` se houver chave, senão `fake` | força o provider |
| `CACHE_TTL` | `10m` | tempo de cache das buscas |

### API

```
GET /healthz
GET /v1/stations/nearby?lat=-3.7319&lng=-38.5267&radius=3000
```

`radius` em metros (100 a 50000, padrão 3000). Resposta:

```json
{
  "count": 8,
  "stations": [
    {
      "id": "…", "name": "…", "address": "…",
      "lat": -3.73, "lng": -38.52, "rating": 4.3,
      "distanceMeters": 300,
      "prices": [
        { "fuel": "gasolina_comum", "price": 6.29, "currency": "BRL",
          "source": "google", "collectedAt": "2026-10-30T12:00:00Z" }
      ]
    }
  ]
}
```

## Rodando o app

Requer Node 20+ e o app **Expo Go** no celular.

```bash
cd app
npm install
# troque pelo IP do seu PC na rede (o celular não enxerga "localhost" do PC)
EXPO_PUBLIC_API_URL=http://192.168.0.10:8080 npx expo start
```

Escaneie o QR code com o Expo Go. PC e celular precisam estar na mesma rede, e a
porta 8080 liberada no firewall (no Nobara: `sudo firewall-cmd --add-port=8080/tcp`).

> No Expo Go o mapa funciona sem chave. Para gerar o APK/IPA com Google Maps, siga a
> seção "Deploy app with Google Maps" da doc do `react-native-maps` no Expo.

## Fontes de dados

- Postos e (quando houver) preços: Google Places API (New)
- Preços oficiais: ANP — Levantamento de Preços de Combustíveis (fase 2)
- Radares e semáforos: © OpenStreetMap contributors, via Overpass API (fase 4)

# Backlog

## Sprint 1 — Postos próximos (fase 1) ✅ em andamento

> 2026-11-01: testado num Android real via Expo Go com o provider fake — mapa, lista,
> chips de raio e origem/idade dos preços funcionando.

- [x] S1-01 Backend: `GET /healthz`
- [x] S1-02 Backend: `GET /v1/stations/nearby?lat&lng&radius` com distância (haversine) e ordenação
- [x] S1-03 Backend: cliente do Places API (New) `places:searchNearby` com field mask
- [x] S1-04 Backend: provider *fake* (roda sem chave do Google)
- [x] S1-05 Backend: cache em memória com TTL (chave = coordenada arredondada + raio)
- [x] S1-06 Backend: testes (haversine, cache, handler, cliente Places com `httptest`)
- [x] S1-07 App: pedir permissão e ler a localização
- [x] S1-08 App: mapa com marcadores + lista por distância
- [x] S1-09 App: abrir rota no app de mapas ao tocar no posto
- [ ] S1-10 Testar com chave real em Fortaleza e anotar se `fuelOptions` vem preenchido
- [ ] S1-11 Tema: mapa segue o modo escuro do sistema, lista fica clara — unificar

## Sprint 2 — Preços da ANP (fase 2)

- [ ] S2-01 Baixar e parsear o CSV "preços por posto revendedor" da ANP
- [ ] S2-02 Persistência (SQLite) de postos ANP e preços
- [ ] S2-03 Job semanal de importação
- [ ] S2-04 Cruzamento ANP ↔ Google (normalização de endereço + distância < 80 m)
- [ ] S2-05 `GET /v1/stations/cheapest?fuel=gasolina&lat&lng&radius` → top 5
- [ ] S2-06 App: aba "Mais baratos" com origem e data do preço

## Sprint 3 — Crowdsourcing

- [ ] S3-01 `POST /v1/stations/{id}/prices`
- [ ] S3-02 Anti-abuso básico (rate limit por dispositivo, faixa de preço plausível)
- [ ] S3-03 App: formulário "informar preço"

## Sprint 4 — Camadas

- [ ] S4-01 Hospitais via Places
- [ ] S4-02 Radares e semáforos via Overpass (cache 24 h)
- [ ] S4-03 App: seletor de camadas + atribuição OSM

## Sprint 5 — Alertas

- [ ] S5-01 Geofencing de radares com `expo-location`
- [ ] S5-02 Notificação local ao se aproximar

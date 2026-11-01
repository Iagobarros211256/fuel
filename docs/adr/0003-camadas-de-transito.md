# ADR 0003 — Radares, semáforos e hospitais

**Status:** aceito · 2026-11-01

## Contexto

O Google Maps não expõe radares nem semáforos pela API. O OpenStreetMap mapeia
radares fixos com `highway=speed_camera` e semáforos com `highway=traffic_signals`,
consultáveis de graça pela Overpass API. Hospitais existem no Places (`hospital`).

## Decisão

- Radares e semáforos: **Overpass API**, com cache agressivo no backend (esses pontos
  mudam pouco). Se a AMC/Detran-CE publicar dados abertos, entram como segunda fonte.
- Hospitais: **Google Places**, reaproveitando o cliente da fase 1.

## Consequências

- (+) Custo zero para radares/semáforos.
- (−) A qualidade depende do mapeamento do OSM em cada região; o app deve deixar claro
  que a camada pode estar incompleta.
- Atribuição obrigatória: "© OpenStreetMap contributors" (licença ODbL).

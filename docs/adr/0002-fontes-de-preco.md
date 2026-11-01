# ADR 0002 — Fontes de preço de combustível

**Status:** aceito · 2026-11-01

## Contexto

- O Google Places API (New) tem o campo `fuelOptions.fuelPrices`, mas a cobertura no
  Brasil é incerta.
- A ANP publica semanalmente, em dados abertos, preços **por posto revendedor**. É a
  fonte oficial, mas é amostral (nem todo posto aparece) e pode ter até uma semana de
  atraso.

## Decisão

Cada preço é um registro com `fonte` (`google` | `anp` | `usuario`), `valor` e
`coletado_em`. O ranking usa o preço **mais recente** disponível para cada posto e
combustível, independentemente da fonte.

Ordem de implementação:
1. Fase 1: o backend já pede `places.fuelOptions` ao Google e repassa o que vier.
2. Fase 2: job semanal da ANP + cruzamento com os postos do Google.
3. Fase 3: preço informado pelo usuário.

## Consequências

- (+) O app nunca mostra um preço sem dizer de onde veio e de quando é.
- (+) Se o Google passar a cobrir o Brasil, o dado entra sem mudar o modelo.
- (−) O cruzamento ANP ↔ Google é heurístico e vai errar às vezes; precisa de testes
  com dados reais de Fortaleza.

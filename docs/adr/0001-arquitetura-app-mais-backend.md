# ADR 0001 — App mobile + backend próprio

**Status:** aceito · 2026-11-01

## Contexto

O app precisa de dados de três fontes (Google Places, ANP, OpenStreetMap). A chave do
Google é paga por requisição e não pode ficar exposta no app. O dado da ANP é um CSV
semanal que precisa ser baixado, tratado e cruzado com os postos do Google.

## Decisão

- **Backend em Go** (stdlib `net/http`, sem framework) que concentra as integrações,
  guarda cache e expõe uma API simples (`/v1/...`) para o app.
- **App em Expo (React Native + TypeScript)** que só fala com o backend.

## Consequências

- (+) A chave do Google fica só no servidor; dá pra trocar de provedor sem mexer no app.
- (+) Cache e job da ANP vivem num lugar só.
- (+) Go é a linguagem que estou aprofundando; Expo deixa testar no celular sem build nativo.
- (−) Mais uma peça para hospedar. Na fase 1 roda local; deploy fica para depois.

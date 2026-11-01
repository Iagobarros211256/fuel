# PRD — Abastece

> Nome provisório. Projeto da trilha "me desafiando".

## 1. Problema

Quem dirige não sabe, sem rodar a cidade, qual posto perto de si tem o combustível
mais barato. Os dados existem (ANP, Google), mas espalhados e sem cruzamento.
Além disso, informações úteis de trânsito — radares/fotossensores, semáforos e
hospitais — ficam em apps diferentes.

## 2. Objetivo

Um app mobile que, a partir da localização do usuário:

1. mostra os postos de combustível próximos num mapa;
2. ranqueia os **5 mais baratos** por tipo de combustível, deixando claro de onde
   veio o preço e **quando** ele foi coletado;
3. mostra camadas opcionais: radares (fotossensores), semáforos e hospitais.

## 3. Usuário

Motorista/motociclista urbano (primeiro alvo: Fortaleza-CE).

## 4. Requisitos funcionais

| ID   | Requisito | Fase |
|------|-----------|------|
| RF01 | Obter a localização atual do usuário (com permissão) | 1 |
| RF02 | Listar postos num raio configurável (padrão 3 km), ordenados por distância | 1 |
| RF03 | Exibir os postos como marcadores no mapa e em lista | 1 |
| RF04 | Importar semanalmente os preços por posto da ANP (CSV de dados abertos) | 2 |
| RF05 | Cruzar posto da ANP com posto do Google (CNPJ/endereço/proximidade) | 2 |
| RF06 | Top 5 mais baratos por combustível (gasolina comum, aditivada, etanol, diesel, GNV) | 2 |
| RF07 | Mostrar origem (`ANP`, `Google`, `usuário`) e data de cada preço | 2 |
| RF08 | Usuário informa o preço que viu no posto (crowdsourcing) | 3 |
| RF09 | Camada de hospitais (Places, tipo `hospital`) | 4 |
| RF10 | Camadas de radares e semáforos (OpenStreetMap / Overpass) | 4 |
| RF11 | Alerta de proximidade de radar (geofencing) | 5 |

## 5. Requisitos não funcionais

- **Custo:** chamadas ao Google Places passam pelo backend, com cache; a chave nunca
  fica no app.
- **Transparência:** nenhum preço sem data. Preço com mais de 7 dias é sinalizado.
- **Rodar sem chave:** o backend tem um provider *fake* para desenvolvimento.
- **Privacidade:** a localização do usuário não é armazenada.

## 6. Fora de escopo (por enquanto)

- Navegação turn-by-turn (o app abre a rota no Google Maps/Waze).
- Contas de usuário (só entram na fase 3, para o crowdsourcing).

## 7. Métricas de sucesso (pessoais)

- Fase 1 rodando no celular via Expo Go, mostrando postos reais de Fortaleza.
- Fase 2 com top 5 batendo com o CSV da ANP da semana.

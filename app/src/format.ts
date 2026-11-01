import type { FuelPrice } from './types';

export function formatDistance(meters: number): string {
  if (meters < 1000) return `${Math.round(meters)} m`;
  return `${(meters / 1000).toFixed(1).replace('.', ',')} km`;
}

export function formatPrice(p: FuelPrice): string {
  return `R$ ${p.price.toFixed(2).replace('.', ',')}`;
}

const FUEL_LABELS: Record<string, string> = {
  gasolina_comum: 'Gasolina',
  gasolina_aditivada: 'Aditivada',
  gasolina_premium: 'Premium',
  etanol: 'Etanol',
  diesel: 'Diesel',
  gnv: 'GNV',
  glp: 'GLP',
};

export function fuelLabel(fuel: string): string {
  return FUEL_LABELS[fuel] ?? fuel;
}

const SOURCE_LABELS: Record<string, string> = {
  google: 'Google',
  anp: 'ANP',
  usuario: 'usuário',
  fake: 'dado fictício',
};

/** "ANP · há 3 dias" — todo preço aparece com origem e idade (ADR 0002). */
export function priceOrigin(p: FuelPrice, now = new Date()): string {
  const days = Math.floor((now.getTime() - new Date(p.collectedAt).getTime()) / 86_400_000);
  const age = days <= 0 ? 'hoje' : days === 1 ? 'ontem' : `há ${days} dias`;
  return `${SOURCE_LABELS[p.source] ?? p.source} · ${age}`;
}

export function isStale(p: FuelPrice, now = new Date()): boolean {
  return now.getTime() - new Date(p.collectedAt).getTime() > 7 * 86_400_000;
}

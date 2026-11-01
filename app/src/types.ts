// Espelha o JSON de GET /v1/stations/nearby do backend.

export type FuelSource = 'google' | 'anp' | 'usuario' | 'fake';

export type FuelPrice = {
  fuel: string;
  price: number;
  currency: string;
  source: FuelSource;
  collectedAt: string; // ISO 8601
};

export type Station = {
  id: string;
  name: string;
  address: string;
  lat: number;
  lng: number;
  rating?: number;
  distanceMeters: number;
  prices: FuelPrice[];
};

export type NearbyResponse = {
  count: number;
  stations: Station[];
};

export type Coords = { latitude: number; longitude: number };

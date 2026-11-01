import type { Coords, NearbyResponse } from './types';

// No celular com Expo Go, "localhost" é o próprio celular. Use o IP do PC na rede:
//   EXPO_PUBLIC_API_URL=http://192.168.0.10:8080 npx expo start
export const API_URL = process.env.EXPO_PUBLIC_API_URL ?? 'http://localhost:8080';

export class ApiError extends Error {}

export async function fetchNearbyStations(
  coords: Coords,
  radiusMeters: number,
  signal?: AbortSignal,
): Promise<NearbyResponse> {
  const params = new URLSearchParams({
    lat: String(coords.latitude),
    lng: String(coords.longitude),
    radius: String(radiusMeters),
  });

  let res: Response;
  try {
    res = await fetch(`${API_URL}/v1/stations/nearby?${params}`, { signal });
  } catch (err) {
    if ((err as Error).name === 'AbortError') throw err;
    throw new ApiError(`Não consegui falar com o servidor (${API_URL}).`);
  }

  const body = await res.json().catch(() => null);
  if (!res.ok) {
    throw new ApiError(body?.error ?? `Erro ${res.status} ao buscar postos.`);
  }
  return body as NearbyResponse;
}

import * as Location from 'expo-location';
import { useCallback, useEffect, useState } from 'react';
import type { Coords } from './types';

type State =
  | { status: 'loading' }
  | { status: 'denied' }
  | { status: 'error'; message: string }
  | { status: 'ready'; coords: Coords };

/** Pede permissão e lê a posição atual. `refresh` lê de novo. */
export function useUserLocation() {
  const [state, setState] = useState<State>({ status: 'loading' });

  const refresh = useCallback(async () => {
    setState({ status: 'loading' });
    try {
      const { status } = await Location.requestForegroundPermissionsAsync();
      if (status !== 'granted') {
        setState({ status: 'denied' });
        return;
      }
      // posição recente em cache aparece rápido; depois refinamos com o GPS
      const last = await Location.getLastKnownPositionAsync({ maxAge: 60_000 });
      if (last) setState({ status: 'ready', coords: pick(last) });

      const current = await Location.getCurrentPositionAsync({
        accuracy: Location.Accuracy.Balanced,
      });
      setState({ status: 'ready', coords: pick(current) });
    } catch (err) {
      setState({ status: 'error', message: (err as Error).message });
    }
  }, []);

  useEffect(() => {
    refresh();
  }, [refresh]);

  return { ...state, refresh };
}

function pick(loc: Location.LocationObject): Coords {
  return { latitude: loc.coords.latitude, longitude: loc.coords.longitude };
}

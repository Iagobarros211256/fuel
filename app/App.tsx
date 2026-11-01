import { StatusBar } from 'expo-status-bar';
import { type ReactNode, useEffect, useRef, useState } from 'react';
import {
  ActivityIndicator,
  FlatList,
  Linking,
  Pressable,
  StyleSheet,
  Text,
  View,
} from 'react-native';
import MapView, { Circle, Marker } from 'react-native-maps';
import { SafeAreaProvider, SafeAreaView } from 'react-native-safe-area-context';

import { ApiError, fetchNearbyStations } from './src/api';
import { StationRow } from './src/components/StationRow';
import type { Station } from './src/types';
import { useUserLocation } from './src/useUserLocation';

const RADII = [1000, 3000, 5000];

export default function App() {
  return (
    <SafeAreaProvider>
      <Home />
    </SafeAreaProvider>
  );
}

function Home() {
  const location = useUserLocation();
  const [radius, setRadius] = useState(3000);
  const [stations, setStations] = useState<Station[]>([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [selectedId, setSelectedId] = useState<string | null>(null);
  const mapRef = useRef<MapView>(null);
  const listRef = useRef<FlatList<Station>>(null);

  const coords = location.status === 'ready' ? location.coords : null;

  // busca os postos sempre que a posição ou o raio mudam
  useEffect(() => {
    if (!coords) return;
    const ctrl = new AbortController();
    setLoading(true);
    setError(null);
    fetchNearbyStations(coords, radius, ctrl.signal)
      .then((res) => setStations(res.stations))
      .catch((err) => {
        if (err.name === 'AbortError') return;
        setError(err instanceof ApiError ? err.message : 'Erro inesperado.');
      })
      .finally(() => setLoading(false));
    return () => ctrl.abort();
  }, [coords?.latitude, coords?.longitude, radius]);

  // enquadra o círculo de busca quando o raio muda
  useEffect(() => {
    if (!coords) return;
    const delta = (radius * 2.4) / 111_320;
    mapRef.current?.animateToRegion(
      { ...coords, latitudeDelta: delta, longitudeDelta: delta },
      400,
    );
  }, [coords?.latitude, coords?.longitude, radius]);

  function select(station: Station, from: 'map' | 'list') {
    setSelectedId(station.id);
    if (from === 'map') {
      const index = stations.findIndex((s) => s.id === station.id);
      if (index >= 0) listRef.current?.scrollToIndex({ index, viewPosition: 0 });
    } else {
      mapRef.current?.animateCamera({ center: { latitude: station.lat, longitude: station.lng } });
    }
  }

  function openRoute(station: Station) {
    const dest = `${station.lat},${station.lng}`;
    Linking.openURL(`https://www.google.com/maps/dir/?api=1&destination=${dest}`);
  }

  if (location.status === 'loading' && !coords) {
    return (
      <Centered>
        <ActivityIndicator size="large" />
        <Text style={styles.muted}>Pegando sua localização…</Text>
      </Centered>
    );
  }
  if (location.status === 'denied') {
    return (
      <Centered>
        <Text style={styles.title}>Preciso da sua localização</Text>
        <Text style={styles.muted}>Sem ela não dá pra achar os postos perto de você.</Text>
        <Button label="Tentar de novo" onPress={location.refresh} />
        <Button label="Abrir configurações" onPress={() => Linking.openSettings()} />
      </Centered>
    );
  }
  if (location.status === 'error') {
    return (
      <Centered>
        <Text style={styles.title}>Não consegui sua localização</Text>
        <Text style={styles.muted}>{location.message}</Text>
        <Button label="Tentar de novo" onPress={location.refresh} />
      </Centered>
    );
  }

  return (
    <SafeAreaView style={styles.container}>
      <StatusBar style="dark" />

      <MapView
        ref={mapRef}
        style={styles.map}
        showsUserLocation
        initialRegion={coords ? { ...coords, latitudeDelta: 0.06, longitudeDelta: 0.06 } : undefined}
      >
        {coords && (
          <Circle
            center={coords}
            radius={radius}
            strokeColor="rgba(29,78,216,0.5)"
            fillColor="rgba(29,78,216,0.06)"
          />
        )}
        {stations.map((s) => (
          <Marker
            key={s.id}
            coordinate={{ latitude: s.lat, longitude: s.lng }}
            title={s.name}
            pinColor={s.id === selectedId ? '#1d4ed8' : '#dc2626'}
            onPress={() => select(s, 'map')}
          />
        ))}
      </MapView>

      <View style={styles.toolbar}>
        <View style={styles.chips}>
          {RADII.map((r) => (
            <Pressable
              key={r}
              onPress={() => setRadius(r)}
              style={[styles.chip, r === radius && styles.chipActive]}
            >
              <Text style={[styles.chipText, r === radius && styles.chipTextActive]}>{r / 1000} km</Text>
            </Pressable>
          ))}
        </View>
        <Text style={styles.muted}>
          {loading ? 'Buscando…' : `${stations.length} postos`}
        </Text>
      </View>

      {error ? (
        <View style={styles.errorBox}>
          <Text style={styles.errorText}>{error}</Text>
          <Button label="Tentar de novo" onPress={location.refresh} />
        </View>
      ) : (
        <FlatList
          ref={listRef}
          data={stations}
          keyExtractor={(s) => s.id}
          renderItem={({ item }) => (
            <StationRow
              station={item}
              selected={item.id === selectedId}
              onPress={() => select(item, 'list')}
              onRoute={() => openRoute(item)}
            />
          )}
          onScrollToIndexFailed={() => {}}
          ListEmptyComponent={
            loading ? undefined : <Text style={[styles.muted, styles.empty]}>Nenhum posto nesse raio.</Text>
          }
        />
      )}
    </SafeAreaView>
  );
}

function Centered({ children }: { children: ReactNode }) {
  return <View style={styles.centered}>{children}</View>;
}

function Button({ label, onPress }: { label: string; onPress: () => void }) {
  return (
    <Pressable onPress={onPress} style={styles.button} accessibilityRole="button">
      <Text style={styles.buttonText}>{label}</Text>
    </Pressable>
  );
}

const styles = StyleSheet.create({
  container: { flex: 1, backgroundColor: '#fff' },
  map: { height: '45%' },
  toolbar: {
    flexDirection: 'row',
    alignItems: 'center',
    justifyContent: 'space-between',
    paddingHorizontal: 16,
    paddingVertical: 10,
    borderBottomWidth: StyleSheet.hairlineWidth,
    borderBottomColor: '#d9dde3',
  },
  chips: { flexDirection: 'row', gap: 8 },
  chip: { paddingHorizontal: 12, paddingVertical: 6, borderRadius: 999, backgroundColor: '#eef1f5' },
  chipActive: { backgroundColor: '#1d4ed8' },
  chipText: { fontSize: 13, fontWeight: '600', color: '#374151' },
  chipTextActive: { color: '#fff' },
  centered: { flex: 1, alignItems: 'center', justifyContent: 'center', padding: 24, gap: 12 },
  title: { fontSize: 18, fontWeight: '700', textAlign: 'center' },
  muted: { fontSize: 13, color: '#6b7280', textAlign: 'center' },
  empty: { marginTop: 24 },
  errorBox: { padding: 24, gap: 12, alignItems: 'center' },
  errorText: { color: '#b91c1c', textAlign: 'center' },
  button: { backgroundColor: '#1d4ed8', paddingHorizontal: 16, paddingVertical: 10, borderRadius: 8 },
  buttonText: { color: '#fff', fontWeight: '600' },
});

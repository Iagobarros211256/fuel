import { Pressable, StyleSheet, Text, View } from 'react-native';
import { formatDistance, formatPrice, fuelLabel, isStale, priceOrigin } from '../format';
import type { Station } from '../types';

type Props = {
  station: Station;
  selected: boolean;
  onPress: () => void;
  onRoute: () => void;
};

export function StationRow({ station, selected, onPress, onRoute }: Props) {
  return (
    <Pressable onPress={onPress} style={[styles.row, selected && styles.selected]}>
      <View style={styles.header}>
        <Text style={styles.name} numberOfLines={1}>
          {station.name}
        </Text>
        <Text style={styles.distance}>{formatDistance(station.distanceMeters)}</Text>
      </View>

      {!!station.address && (
        <Text style={styles.address} numberOfLines={1}>
          {station.address}
        </Text>
      )}

      {station.prices.length === 0 ? (
        <Text style={styles.noPrice}>Sem preço informado</Text>
      ) : (
        <View style={styles.prices}>
          {station.prices.map((p) => (
            <View key={p.fuel} style={styles.price}>
              <Text style={styles.priceValue}>
                {fuelLabel(p.fuel)} {formatPrice(p)}
              </Text>
              <Text style={[styles.priceOrigin, isStale(p) && styles.stale]}>{priceOrigin(p)}</Text>
            </View>
          ))}
        </View>
      )}

      {selected && (
        <Pressable onPress={onRoute} style={styles.routeButton} accessibilityRole="button">
          <Text style={styles.routeText}>Como chegar</Text>
        </Pressable>
      )}
    </Pressable>
  );
}

const styles = StyleSheet.create({
  row: {
    paddingHorizontal: 16,
    paddingVertical: 12,
    borderBottomWidth: StyleSheet.hairlineWidth,
    borderBottomColor: '#d9dde3',
    backgroundColor: '#fff',
  },
  selected: { backgroundColor: '#f1f6ff' },
  header: { flexDirection: 'row', justifyContent: 'space-between', gap: 12 },
  name: { flex: 1, fontSize: 16, fontWeight: '600', color: '#111827' },
  distance: { fontSize: 14, fontWeight: '600', color: '#1d4ed8' },
  address: { marginTop: 2, fontSize: 13, color: '#6b7280' },
  noPrice: { marginTop: 6, fontSize: 13, color: '#9ca3af', fontStyle: 'italic' },
  prices: { marginTop: 6, flexDirection: 'row', flexWrap: 'wrap', gap: 12 },
  price: {},
  priceValue: { fontSize: 14, fontWeight: '600', color: '#111827' },
  priceOrigin: { fontSize: 11, color: '#6b7280' },
  stale: { color: '#b45309' },
  routeButton: {
    marginTop: 10,
    alignSelf: 'flex-start',
    backgroundColor: '#1d4ed8',
    paddingHorizontal: 14,
    paddingVertical: 8,
    borderRadius: 8,
  },
  routeText: { color: '#fff', fontWeight: '600' },
});

import { Pressable, StyleSheet, Text, View } from "react-native";

interface Props {
  page: number;
  totalPages: number;
  onChange: (page: number) => void;
}

export function Pagination({ page, totalPages, onChange }: Props) {
  const canPrev = page > 1;
  const canNext = page < totalPages;

  return (
    <View style={styles.row} testID="pagination">
      <Pressable
        testID="pagination-prev"
        disabled={!canPrev}
        onPress={() => canPrev && onChange(page - 1)}
        style={[styles.btn, !canPrev && styles.btnDisabled]}
      >
        <Text style={styles.btnText}>‹ Prev</Text>
      </Pressable>

      <Text style={styles.info} testID="pagination-info">
        {totalPages === 0 ? "0 / 0" : `${page} / ${totalPages}`}
      </Text>

      <Pressable
        testID="pagination-next"
        disabled={!canNext}
        onPress={() => canNext && onChange(page + 1)}
        style={[styles.btn, !canNext && styles.btnDisabled]}
      >
        <Text style={styles.btnText}>Next ›</Text>
      </Pressable>
    </View>
  );
}

const styles = StyleSheet.create({
  row: {
    flexDirection: "row",
    justifyContent: "space-between",
    alignItems: "center",
    marginTop: 12,
  },
  btn: {
    paddingHorizontal: 14,
    paddingVertical: 8,
    borderRadius: 8,
    borderWidth: 1,
    borderColor: "#d1d5db",
    backgroundColor: "#fff",
  },
  btnDisabled: { opacity: 0.4 },
  btnText: { color: "#374151" },
  info: { color: "#6b7280" },
});

import { Pressable, StyleSheet, Text, View } from "react-native";
import { TaskStatus } from "../types/task";

interface Props {
  value: TaskStatus | "";
  onChange: (s: TaskStatus | "") => void;
}

const OPTIONS: { value: TaskStatus | ""; label: string }[] = [
  { value: "", label: "Semua" },
  { value: "todo", label: "Todo" },
  { value: "in_progress", label: "In Progress" },
  { value: "done", label: "Done" },
];

export function StatusFilter({ value, onChange }: Props) {
  return (
    <View style={styles.row} testID="status-filter">
      {OPTIONS.map((opt) => {
        const active = value === opt.value;
        return (
          <Pressable
            key={opt.value || "all"}
            testID={`status-chip-${opt.value || "all"}`}
            onPress={() => onChange(opt.value)}
            style={[styles.chip, active && styles.chipActive]}
          >
            <Text style={[styles.chipText, active && styles.chipTextActive]}>
              {opt.label}
            </Text>
          </Pressable>
        );
      })}
    </View>
  );
}

const styles = StyleSheet.create({
  row: { flexDirection: "row", flexWrap: "wrap", gap: 8, marginBottom: 12 },
  chip: {
    paddingHorizontal: 12,
    paddingVertical: 6,
    borderRadius: 999,
    borderWidth: 1,
    borderColor: "#d1d5db",
    backgroundColor: "#fff",
  },
  chipActive: { backgroundColor: "#2563eb", borderColor: "#2563eb" },
  chipText: { color: "#374151", fontSize: 13 },
  chipTextActive: { color: "#fff", fontWeight: "600" },
});

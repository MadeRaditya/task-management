import { FlatList, Pressable, StyleSheet, Text, View } from "react-native";
import { Task } from "../types/task";

interface Props {
  tasks: Task[];
  onEdit: (task: Task) => void;
  onDelete: (task: Task) => void;
  onRefresh?: () => void;
  refreshing?: boolean;
  emptyMessage?: string;
}

const STATUS_LABEL: Record<string, string> = {
  todo: "Todo",
  in_progress: "In Progress",
  done: "Done",
};

export function TaskList({
  tasks,
  onEdit,
  onDelete,
  onRefresh,
  refreshing,
  emptyMessage = "Tidak ada task.",
}: Props) {
  return (
    <FlatList
      testID="task-list"
      data={tasks}
      keyExtractor={(item) => String(item.id)}
      onRefresh={onRefresh}
      refreshing={!!refreshing}
      ListEmptyComponent={
        <View style={styles.empty}>
          <Text style={styles.emptyText}>{emptyMessage}</Text>
        </View>
      }
      renderItem={({ item }) => (
        <View style={styles.item} testID={`task-item-${item.id}`}>
          <View style={styles.itemHeader}>
            <Text style={styles.itemTitle} numberOfLines={1}>
              {item.title}
            </Text>
            <View
              style={[styles.badge, styles[`badge_${item.status}` as const]]}
            >
              <Text style={styles.badgeText}>{STATUS_LABEL[item.status]}</Text>
            </View>
          </View>

          {!!item.description && (
            <Text style={styles.itemDesc} numberOfLines={2}>
              {item.description}
            </Text>
          )}

          <View style={styles.itemFooter}>
            <Text style={styles.itemMeta}>
              {item.assignee ? `👤 ${item.assignee}` : "👤 -"}
            </Text>
            <View style={styles.actions}>
              <Pressable
                testID={`task-edit-${item.id}`}
                onPress={() => onEdit(item)}
                style={[styles.btn, styles.btnGhost]}
              >
                <Text style={styles.btnGhostText}>Edit</Text>
              </Pressable>
              <Pressable
                testID={`task-delete-${item.id}`}
                onPress={() => onDelete(item)}
                style={[styles.btn, styles.btnDanger]}
              >
                <Text style={styles.btnDangerText}>Hapus</Text>
              </Pressable>
            </View>
          </View>
        </View>
      )}
    />
  );
}

const styles = StyleSheet.create({
  item: {
    backgroundColor: "#fff",
    borderRadius: 10,
    padding: 12,
    marginBottom: 8,
    borderWidth: 1,
    borderColor: "#e5e7eb",
  },
  itemHeader: {
    flexDirection: "row",
    alignItems: "center",
    justifyContent: "space-between",
    gap: 8,
  },
  itemTitle: { flex: 1, fontSize: 16, fontWeight: "600", color: "#111827" },
  itemDesc: { color: "#4b5563", marginTop: 4 },
  itemFooter: {
    flexDirection: "row",
    alignItems: "center",
    justifyContent: "space-between",
    marginTop: 8,
  },
  itemMeta: { color: "#6b7280", fontSize: 13 },
  actions: { flexDirection: "row", gap: 8 },
  badge: {
    paddingHorizontal: 8,
    paddingVertical: 2,
    borderRadius: 999,
  },
  badge_todo: { backgroundColor: "#e5e7eb" },
  badge_in_progress: { backgroundColor: "#fef3c7" },
  badge_done: { backgroundColor: "#d1fae5" },
  badgeText: { fontSize: 11, color: "#111827", fontWeight: "600" },
  btn: { paddingHorizontal: 10, paddingVertical: 6, borderRadius: 6 },
  btnGhost: { backgroundColor: "#f3f4f6" },
  btnGhostText: { color: "#374151", fontSize: 13 },
  btnDanger: { backgroundColor: "#fee2e2" },
  btnDangerText: { color: "#b91c1c", fontSize: 13 },
  empty: { padding: 24, alignItems: "center" },
  emptyText: { color: "#6b7280" },
});

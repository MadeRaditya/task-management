import { useEffect, useMemo, useState } from "react";
import {
  Alert,
  Platform,
  SafeAreaView,
  StyleSheet,
  Text,
  View,
} from "react-native";
import { StatusBar } from "expo-status-bar";
import { deleteTask, updateTask } from "./src/api/tasks";
import { EditTaskModal } from "./src/components/EditTaskModal";
import { LoadingState } from "./src/components/LoadingState";
import { Pagination } from "./src/components/Pagination";
import { SearchBar } from "./src/components/SearchBar";
import { StatusFilter } from "./src/components/StatusFilter";
import { TaskList } from "./src/components/TaskList";
import { useTasks } from "./src/hooks/useTasks";
import { useDebounce } from "./src/utils/useDebounce";
import { Task, TaskStatus } from "./src/types/task";

export default function App() {
  const [searchInput, setSearchInput] = useState("");
  const debouncedKeyword = useDebounce(searchInput, 400);

  const {
    tasks,
    loading,
    error,
    filters,
    total,
    totalPages,
    setKeyword,
    setStatus,
    setPage,
    refresh,
  } = useTasks({ keyword: "", status: "", page: 1, limit: 10 });

  useEffect(() => {
    setKeyword(debouncedKeyword);
  }, [debouncedKeyword]);

  const [editing, setEditing] = useState<Task | null>(null);

  const handleUpdate = async (
    id: number,
    payload: {
      title: string;
      description: string;
      status: TaskStatus;
      assignee: string;
    },
  ) => {
    await updateTask(id, payload);
    refresh();
  };

  const handleDelete = (task: Task) => {
    const doDelete = async () => {
      try {
        await deleteTask(task.id);
        refresh();
      } catch (e: any) {
        if (Platform.OS === "web") {
          window.alert(e?.message ?? "Gagal menghapus");
        } else {
          Alert.alert("Gagal", e?.message ?? "Gagal menghapus");
        }
      }
    };

    if (Platform.OS === "web") {
      if (window.confirm(`Hapus task "${task.title}"?`)) doDelete();
    } else {
      Alert.alert("Konfirmasi", `Hapus task "${task.title}"?`, [
        { text: "Batal", style: "cancel" },
        { text: "Hapus", style: "destructive", onPress: doDelete },
      ]);
    }
  };

  return (
    <SafeAreaView style={styles.safe}>
      <StatusBar style="auto" />
      <View style={styles.container}>
        <Text style={styles.heading}>Task Management</Text>
        <Text style={styles.subheading}>
          {total} task{filters.status ? ` · ${filters.status}` : ""}
        </Text>

        <SearchBar value={searchInput} onChangeText={setSearchInput} />
        <StatusFilter value={filters.status} onChange={(s) => setStatus(s)} />

        {error && <Text style={styles.error}>{error}</Text>}

        {loading && tasks.length === 0 ? (
          <LoadingState />
        ) : (
          <TaskList
            tasks={tasks}
            onEdit={(t) => setEditing(t)}
            onDelete={handleDelete}
            onRefresh={refresh}
            refreshing={loading}
          />
        )}

        <Pagination
          page={filters.page}
          totalPages={totalPages}
          onChange={setPage}
        />
      </View>

      <EditTaskModal
        visible={!!editing}
        task={editing}
        onClose={() => setEditing(null)}
        onSubmit={handleUpdate}
      />
    </SafeAreaView>
  );
}

const styles = StyleSheet.create({
  safe: { flex: 1, backgroundColor: "#f9fafb" },
  container: {
    flex: 1,
    padding: 16,
    maxWidth: 720,
    width: "100%",
    alignSelf: "center",
  },
  heading: { fontSize: 24, fontWeight: "800", color: "#111827" },
  subheading: { color: "#6b7280", marginBottom: 12 },
  error: { color: "#dc2626", marginBottom: 8 },
});
